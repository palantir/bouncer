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

package aws

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
)

// GetAllASGs returns all ASGs visible with your client, with no filters
func (c *Clients) GetAllASGs() ([]*autoscaling.Group, error) {
	var nexttoken *string
	var asgs []*autoscaling.Group
	var input *autoscaling.DescribeAutoScalingGroupsInput
	var err error
	var output *autoscaling.DescribeAutoScalingGroupsOutput

	for {
		input = &autoscaling.DescribeAutoScalingGroupsInput{
			NextToken: nexttoken,
		}

		output, err = c.ASGClient.DescribeAutoScalingGroups(input)
		if err != nil {
			return nil, errors.Wrap(err, "Error describing ASGs")
		}

		for _, asg := range output.AutoScalingGroups {
			asgs = append(asgs, asg)
		}
		nexttoken = output.NextToken

		if nexttoken == nil {
			break
		} else {
			time.Sleep(apiSleepTime)
		}
	}

	return asgs, nil
}

// GetASG gets the *autoscaling.Group that matches for the name given
func (c *Clients) GetASG(asgName *string) (*autoscaling.Group, error) {
	var asgs []*autoscaling.Group

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{asgName},
	}

	output, err := c.ASGClient.DescribeAutoScalingGroups(input)
	if err != nil {
		return nil, errors.Wrap(err, "Error describing ASGs")
	}

	for _, asg := range output.AutoScalingGroups {
		asgs = append(asgs, asg)
	}

	if len(asgs) != 1 {
		return nil, errors.Errorf("ASG Name '%s' matched '%v' ASGs, expecting it to match 1 (looking in region %s, have you set AWS_DEFAULT_REGION?)", *asgName, len(asgs), c.ASGClient.SigningRegion)
	}

	return asgs[0], nil
}

// GetASGTagValue returns a pointer to the value for the given tag key
func GetASGTagValue(asg *autoscaling.Group, key string) *string {
	for _, tag := range asg.Tags {
		if tag != nil {
			if strings.ToLower(*tag.Key) == strings.ToLower(key) {
				return tag.Value
			}
		}
	}
	return nil
}

// GetLaunchConfiguration returns the LC object of the given ASG
func (c *Clients) GetLaunchConfiguration(asg *autoscaling.Group) (*autoscaling.LaunchConfiguration, error) {
	var lcs []*string
	lcs = append(lcs, asg.LaunchConfigurationName)
	input := autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: lcs,
	}
	output, err := c.ASGClient.DescribeLaunchConfigurations(&input)
	if err != nil {
		return nil, errors.Wrapf(err, "Error describing launch configuration %s for ASG %s", *asg.LaunchConfigurationName, *asg.AutoScalingGroupName)
	}
	return output.LaunchConfigurations[0], nil
}

// CompleteLifecycleAction calls https://docs.aws.amazon.com/cli/latest/reference/autoscaling/complete-lifecycle-action.html
func (c *Clients) CompleteLifecycleAction(asgName *string, instID *string, lifecycleHook *string, result *string) error {
	input := autoscaling.CompleteLifecycleActionInput{
		AutoScalingGroupName:  asgName,
		InstanceId:            instID,
		LifecycleActionResult: result,
		LifecycleHookName:     lifecycleHook,
	}

	_, err := c.ASGClient.CompleteLifecycleAction(&input)
	return errors.Wrapf(err, "error completing lifecycle hook %s for instance %s", *lifecycleHook, *instID)
}

// TerminateInstanceInASG calls https://docs.aws.amazon.com/cli/latest/reference/autoscaling/terminate-instance-in-auto-scaling-group.html
func (c *Clients) TerminateInstanceInASG(instID *string) error {
	// This call decrements the desired capacity so that we don't get into a race condition
	// where the replacement starts booting before the node we've told to terminate has terminated
	dumb := true
	input := autoscaling.TerminateInstanceInAutoScalingGroupInput{
		InstanceId:                     instID,
		ShouldDecrementDesiredCapacity: &dumb,
	}
	_, err := c.ASGClient.TerminateInstanceInAutoScalingGroup(&input)
	return errors.Wrapf(err, "error terminating instance %s", *instID)
}

// SetDesiredCapacity sets the desired capacity of given ASG to given value
func (c *Clients) SetDesiredCapacity(asg *autoscaling.Group, desiredCapacity *int64) error {
	input := autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: asg.AutoScalingGroupName,
		DesiredCapacity:      desiredCapacity,
	}
	_, err := c.ASGClient.SetDesiredCapacity(&input)
	return errors.Wrapf(err, "error setting desired capacity for %s", *asg.AutoScalingGroupName)
}
