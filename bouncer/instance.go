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
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"

	"github.com/palantir/bouncer/aws"
)

// Instance tracks the AWS representations of an EC2 instance as well as the metadata we care about it
type Instance struct {
	EC2Instance      *ec2.Instance
	ASGInstance      *autoscaling.Instance
	AutoscalingGroup *autoscaling.Group
	IsOld            bool
	IsHealthy        bool
	PreTerminateCmd  *string
}

// NewInstance returns a new bouncer.Instance object
func NewInstance(ac *aws.Clients, asg *autoscaling.Group, asgInst *autoscaling.Instance, launchConfigName *string, force bool, startTime time.Time, preTerminateCmd *string) (*Instance, error) {
	var ec2Inst *ec2.Instance
	err := retry(apiRetryCount, apiRetrySleep, func() (err error) {
		ec2Inst, err = ac.ASGInstToEC2Inst(asgInst)
		return
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Error converting ASG Inst to EC2 inst for %s", *asgInst.InstanceId)
	}

	inst := Instance{
		EC2Instance:      ec2Inst,
		ASGInstance:      asgInst,
		AutoscalingGroup: asg,
		IsOld:            isInstanceOld(asgInst, ec2Inst, launchConfigName, force, startTime),
		IsHealthy:        isInstanceHealthy(asgInst, ec2Inst),
		PreTerminateCmd:  preTerminateCmd,
	}

	return &inst, nil
}

func isInstanceOld(asgInst *autoscaling.Instance, ec2Inst *ec2.Instance, launchConfigName *string, force bool, startTime time.Time) bool {
	if asgInst.LaunchConfigurationName == nil {
		log.WithFields(log.Fields{
			"InstanceID": *asgInst.InstanceId,
		}).Debug("Instance marked as old because launch config is nil")
		return true
	}

	if *asgInst.LaunchConfigurationName != *launchConfigName {
		log.WithFields(log.Fields{
			"InstanceID":   *asgInst.InstanceId,
			"LaunchConfig": *asgInst.LaunchConfigurationName,
		}).Debug("Instance marked as old because of launch config")
		return true
	}

	// In force mode, mark any node that was launched before this runner was started as old
	if force {
		if startTime.After(*ec2Inst.LaunchTime) {
			log.WithFields(log.Fields{
				"InstanceID": *asgInst.InstanceId,
				"LaunchTime": *ec2Inst.LaunchTime,
			}).Debug("Instance marked as old because of launch time (force mode)")
			return true
		}
	}

	return false
}

func isInstanceHealthy(asgInst *autoscaling.Instance, ec2Inst *ec2.Instance) bool {
	if *ec2Inst.State.Name != "running" {
		return false
	}

	if *asgInst.LifecycleState != "InService" {
		return false
	}

	return true
}
