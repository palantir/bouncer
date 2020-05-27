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
	"github.com/palantir/bouncer/aws"
	"github.com/pkg/errors"
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
func NewInstance(ac *aws.Clients, asg *autoscaling.Group, asgInst *autoscaling.Instance, force bool, startTime time.Time, preTerminateCmd *string) (*Instance, error) {
	var ec2Inst *ec2.Instance
	err := retry(apiRetryCount, apiRetrySleep, func() (err error) {
		ec2Inst, err = ac.ASGInstToEC2Inst(asgInst)
		return
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error converting ASG Inst to EC2 inst for %s", *asgInst.InstanceId)
	}

	lts := ac.GetLaunchTemplateSpec(asg)

	var ec2LTplVersion *string
	err = retry(apiRetryCount, apiRetrySleep, func() (err error) {
		ec2LTplVersion, err = ac.ASGLTplVersionToEC2LTplVersion(lts)
		return
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error resolving LaunchTemplate %s Version to actual version number", *lts.LaunchTemplateId)
	}

	inst := Instance{
		EC2Instance:      ec2Inst,
		ASGInstance:      asgInst,
		AutoscalingGroup: asg,
		IsOld:            isInstanceOld(asgInst, ec2Inst, asg.LaunchConfigurationName, lts, ec2LTplVersion, force, startTime),
		IsHealthy:        isInstanceHealthy(asgInst, ec2Inst),
		PreTerminateCmd:  preTerminateCmd,
	}

	return &inst, nil
}

func isInstanceOld(asgInst *autoscaling.Instance, ec2Inst *ec2.Instance, asgLCName *string, asgLT *autoscaling.LaunchTemplateSpecification, asgLTVer *string, force bool, startTime time.Time) bool {
	if asgLCName != nil {
		// This machine is using LaunchConfigs

		if asgInst.LaunchConfigurationName == nil {
			log.WithFields(log.Fields{
				"InstanceID": *asgInst.InstanceId,
			}).Debug("Instance marked as old because launch config has been deleted")

			return true
		} else if *asgInst.LaunchConfigurationName != *asgLCName {
			log.WithFields(log.Fields{
				"InstanceID":           *asgInst.InstanceId,
				"InstanceLaunchConfig": *asgInst.LaunchConfigurationName,
				"GroupLaunchConfig":    *asgLCName,
			}).Debug("Instance marked as old because launch config differs from that of its ASG")

			return true
		}
	} else if asgInst.LaunchTemplate != nil && asgLT != nil {
		// This machine is using LaunchTemplates

		if *asgInst.LaunchTemplate.Version != *asgLTVer {
			log.WithFields(log.Fields{
				"InstanceID":                         *asgInst.InstanceId,
				"InstanceLaunchTemplateId":           *asgInst.LaunchTemplate.LaunchTemplateId,
				"InstanceLaunchTemplateName":         *asgInst.LaunchTemplate.LaunchTemplateName,
				"InstanceLaunchTemplateVersion":      *asgInst.LaunchTemplate.Version,
				"GroupLaunchTemplateId":              *asgLT.LaunchTemplateId,
				"GroupLaunchTemplateName":            *asgLT.LaunchTemplateName,
				"GroupLaunchTemplateVersion":         *asgLT.Version,
				"ResovledGroupLaunchTemplateVersion": *asgLTVer,
			}).Debug("Instance marked as old because launchTemplate version is old")

			return true
		} else if *asgInst.LaunchTemplate.LaunchTemplateId != *asgLT.LaunchTemplateId {
			log.WithFields(log.Fields{
				"InstanceID":                    *asgInst.InstanceId,
				"InstanceLaunchTemplateId":      *asgInst.LaunchTemplate.LaunchTemplateId,
				"InstanceLaunchTemplateName":    *asgInst.LaunchTemplate.LaunchTemplateName,
				"InstanceLaunchTemplateVersion": *asgInst.LaunchTemplate.Version,
				"GroupLaunchTemplateId":         *asgLT.LaunchTemplateId,
				"GroupLaunchTemplateName":       *asgLT.LaunchTemplateName,
				"GroupLaunchTemplateVersion":    *asgLT.Version,
			}).Debug("Instance marked as old because launchTemplate differs from that of its ASG")

			return true
		}
	} else {
		// Using neither - seems to only happen as part of a race condition during migrating from LC to LT

		log.WithFields(log.Fields{
			"InstanceID": *asgInst.InstanceId,
		}).Debug("Instance marked as old because the ASG has neither LC or LT, it must be being transitioned")

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
