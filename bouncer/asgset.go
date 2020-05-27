// Copyright 2017 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bouncer

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/palantir/bouncer/aws"
	"github.com/pkg/errors"
)

// ASGSet has a slice of ASG objects and some functions against them
// This object is recomputed every run of bouncer because it takes actual instance status into account
type ASGSet struct {
	ASGs []*ASG
}

func newASGSet(ac *aws.Clients, desiredASGs []*DesiredASG, force bool, startTime time.Time) (*ASGSet, error) {
	var asgs []*ASG

	for _, desASG := range desiredASGs {
		asg, err := NewASG(ac, desASG, force, startTime)
		if err != nil {
			return nil, errors.Wrapf(err, "Error getting information for ASG %s", desASG.AsgName)
		}

		asgs = append(asgs, asg)
	}

	asgSet := ASGSet{
		ASGs: asgs,
	}

	return &asgSet, nil
}

// GetImmutableInstances returns instances which are in autoscaling events that we can't manipulate by completing lifecycle actions
func (a *ASGSet) GetImmutableInstances() []*Instance {
	var instances []*Instance
	for _, asg := range a.ASGs {
		for _, inst := range asg.Instances {
			if *inst.ASGInstance.LifecycleState == autoscaling.LifecycleStateTerminating || *inst.ASGInstance.LifecycleState == autoscaling.LifecycleStatePending || *inst.ASGInstance.LifecycleState == autoscaling.LifecycleStateTerminatingProceed {
				instances = append(instances, inst)
			}
		}
	}

	return instances
}

// GetUnhealthyNewInstances returns all instances which are on the latest launch configuration but are unhealthy
func (a *ASGSet) GetUnhealthyNewInstances() []*Instance {
	var instances []*Instance
	for _, asg := range a.ASGs {
		for _, inst := range asg.Instances {
			if !inst.IsOld && !inst.IsHealthy {
				instances = append(instances, inst)
			}
		}
	}

	return instances
}

// GetTerminatingInstances returns all instances which are in the process of terminating
func (a *ASGSet) GetTerminatingInstances() []*Instance {
	var terminatingInstances []*Instance
	for _, asg := range a.ASGs {
		for _, inst := range asg.Instances {
			if strings.HasPrefix(*inst.ASGInstance.LifecycleState, "Terminating") {
				terminatingInstances = append(terminatingInstances, inst)
			}
		}
	}
	return terminatingInstances
}

// GetOldInstances returns all instances which are on an outdated launch configuration
func (a *ASGSet) GetOldInstances() []*Instance {
	var oldInstances []*Instance
	for _, asg := range a.ASGs {
		for _, inst := range asg.Instances {
			if inst.IsOld {
				oldInstances = append(oldInstances, inst)
			}
		}
	}
	return oldInstances
}

// GetNewInstances returns all instances which are on an outdated launch configuration
func (a *ASGSet) GetNewInstances() []*Instance {
	var newInstances []*Instance
	for _, asg := range a.ASGs {
		for _, inst := range asg.Instances {
			if !inst.IsOld {
				newInstances = append(newInstances, inst)
			}
		}
	}
	return newInstances
}

// GetBestOldInstance returns the instance which is the best candidate to be bounced
func (a *ASGSet) GetBestOldInstance() *Instance {
	var bestInstance *Instance
	oldInstances := a.GetOldInstances()
	for _, inst := range oldInstances {
		if bestInstance == nil {
			bestInstance = inst
		} else if !inst.IsHealthy && bestInstance.IsHealthy {
			bestInstance = inst
		} else if (inst.IsHealthy == bestInstance.IsHealthy) && inst.EC2Instance.LaunchTime.Before(*bestInstance.EC2Instance.LaunchTime) {
			bestInstance = inst
		}
	}
	return bestInstance
}

// GetActualBadCounts returns all ASGs whose desired counts don't match their actual counts
func (a *ASGSet) GetActualBadCounts() []*ASG {
	var badCountASGs []*ASG
	for _, asg := range a.ASGs {
		if *asg.ASG.DesiredCapacity != int64(len(asg.Instances)) {
			badCountASGs = append(badCountASGs, asg)
		}
	}
	return badCountASGs
}

// GetDivergedASGs returns all ASGs whose desired counts don't match what their desired counts should be
func (a *ASGSet) GetDivergedASGs() []*ASG {
	var badCountASGs []*ASG
	for _, asg := range a.ASGs {
		if *asg.ASG.DesiredCapacity != asg.DesiredASG.DesiredCapacity {
			badCountASGs = append(badCountASGs, asg)
		}
	}
	return badCountASGs
}

// IsOldInstance prints all old instances and returns true/false whether it found any
func (a *ASGSet) IsOldInstance() bool {
	isOldInstance := false

	allOld := a.GetOldInstances()
	for _, old := range allOld {
		log.WithFields(log.Fields{
			"InstanceID": *old.ASGInstance.InstanceId,
			"ASG":        *old.AutoscalingGroup.AutoScalingGroupName,
		}).Info("Instance is old")
		isOldInstance = true
	}

	return isOldInstance
}

// IsNewInstance prints all new instances and returns true/false whether it found any
func (a *ASGSet) IsNewInstance() bool {
	isNewInstance := false

	allNew := a.GetNewInstances()
	for _, new := range allNew {
		log.WithFields(log.Fields{
			"InstanceID": *new.ASGInstance.InstanceId,
			"ASG":        *new.AutoscalingGroup.AutoScalingGroupName,
		}).Info("Instance is new")
		isNewInstance = true
	}

	return isNewInstance
}

// IsTerminating prints all instances in the process of terminating and returns true/false whether it found any
func (a *ASGSet) IsTerminating() bool {
	isTerminating := false

	allTerminating := a.GetTerminatingInstances()
	for _, inst := range allTerminating {
		log.WithFields(log.Fields{
			"ASG":        *inst.AutoscalingGroup.AutoScalingGroupName,
			"InstanceID": *inst.ASGInstance.InstanceId,
			"State":      *inst.ASGInstance.LifecycleState,
		}).Info("Waiting for instance to die")
		isTerminating = true
	}

	return isTerminating
}

// IsCountMismatch prints all instances whose desired_capacity doesn't match running instances and returns true/false whether it found any
func (a *ASGSet) IsCountMismatch() bool {
	isCountMismatch := false

	badActualCounts := a.GetActualBadCounts()
	for _, asg := range badActualCounts {
		log.WithFields(log.Fields{
			"DesiredCapacity": *asg.ASG.DesiredCapacity,
			"InstanceCount":   len(asg.Instances),
			"ASG":             *asg.ASG.AutoScalingGroupName,
		}).Info("Waiting for instance count to match desired capacity")
		isCountMismatch = true
	}

	return isCountMismatch
}

// IsImmutableAutoscalingEvent prints all instances who are in a state we can't affect and returns true/false whether it found any
func (a *ASGSet) IsImmutableAutoscalingEvent() bool {
	isEvent := false

	immutable := a.GetImmutableInstances()
	for _, inst := range immutable {
		log.WithFields(log.Fields{
			"ASG":        *inst.AutoscalingGroup.AutoScalingGroupName,
			"InstanceID": *inst.ASGInstance.InstanceId,
			"State":      *inst.ASGInstance.LifecycleState,
		}).Info("Instance is in transient state")
		isEvent = true
	}

	return isEvent
}

// IsNewUnhealthy prints all instances who are running latest LC but not yet healthy and returns true/false whether it found any
func (a *ASGSet) IsNewUnhealthy() bool {
	isNewUnhealthy := false

	newUnhealthy := a.GetUnhealthyNewInstances()
	for _, inst := range newUnhealthy {
		state := *inst.ASGInstance.LifecycleState
		var msg string

		switch state {
		case autoscaling.LifecycleStateTerminating, autoscaling.LifecycleStateTerminatingProceed, autoscaling.LifecycleStateTerminatingWait:
			msg = "Waiting for unhealthy new instance to get out of the way"
		default:
			msg = "Waiting for new instance to become healthy"
		}

		log.WithFields(log.Fields{
			"ASG":        *inst.AutoscalingGroup.AutoScalingGroupName,
			"InstanceID": *inst.ASGInstance.InstanceId,
			"State":      state,
		}).Info(msg)
		isNewUnhealthy = true
	}

	return isNewUnhealthy
}
