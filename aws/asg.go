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
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	at "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/pkg/errors"
)

// GetAllASGs returns all ASGs visible with your client, with no filters
func (c *Clients) GetAllASGs(ctx context.Context) ([]*at.AutoScalingGroup, error) {
	var nexttoken *string
	var asgs []*at.AutoScalingGroup
	var input *autoscaling.DescribeAutoScalingGroupsInput
	var err error
	var output *autoscaling.DescribeAutoScalingGroupsOutput

	for {
		input = &autoscaling.DescribeAutoScalingGroupsInput{
			NextToken: nexttoken,
		}

		output, err = c.ASGClient.DescribeAutoScalingGroups(ctx, input)
		if err != nil {
			return nil, errors.Wrap(err, "Error describing ASGs")
		}

		for _, asg := range output.AutoScalingGroups {
			asgs = append(asgs, &asg)
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
func (c *Clients) GetASG(ctx context.Context, asgName string) (*at.AutoScalingGroup, error) {
	var asgs []*at.AutoScalingGroup

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{asgName},
	}

	output, err := c.ASGClient.DescribeAutoScalingGroups(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "Error describing ASGs")
	}

	for _, asg := range output.AutoScalingGroups {
		asgs = append(asgs, &asg)
	}

	if len(asgs) != 1 {
		return nil, errors.Errorf("ASG Name '%s' matched '%v' ASGs, expecting it to match 1 (have you set AWS_DEFAULT_REGION?)", asgName, len(asgs))
	}

	return asgs[0], nil
}

// GetASGTagValue returns a pointer to the value for the given tag key
func GetASGTagValue(asg *at.AutoScalingGroup, key string) *string {
	for _, tag := range asg.Tags {
		if tag.Key != nil {
			if strings.EqualFold(*tag.Key, key) {
				return tag.Value
			}
		}
	}
	return nil
}

// GetLaunchConfiguration returns the LC object of the given ASG
func (c *Clients) GetLaunchConfiguration(ctx context.Context, asg *at.AutoScalingGroup) (*at.LaunchConfiguration, error) {
	var lcs []string
	lcs = append(lcs, *asg.LaunchConfigurationName)
	input := autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: lcs,
	}
	output, err := c.ASGClient.DescribeLaunchConfigurations(ctx, &input)
	if err != nil {
		return nil, errors.Wrapf(err, "Error describing launch configuration %s for ASG %s", *asg.LaunchConfigurationName, *asg.AutoScalingGroupName)
	}
	return &output.LaunchConfigurations[0], nil
}

// GetLaunchTemplateSpec returns the LT spec for a given ASG, if it has one
func (c *Clients) GetLaunchTemplateSpec(asg *at.AutoScalingGroup) *at.LaunchTemplateSpecification {
	// First, let's check the direct launch template spec property of the ASG
	if asg.LaunchTemplate != nil {
		return asg.LaunchTemplate
	}

	// Otherwise, let's check for a MixedPolicy
	if asg.MixedInstancesPolicy != nil && asg.MixedInstancesPolicy.LaunchTemplate != nil {
		return asg.MixedInstancesPolicy.LaunchTemplate.LaunchTemplateSpecification
	}

	// Otherwise, it looks like we're probably not using LaunchTemplates
	return nil
}

// CompleteLifecycleAction calls https://docs.aws.amazon.com/cli/latest/reference/autoscaling/complete-lifecycle-action.html
func (c *Clients) CompleteLifecycleAction(ctx context.Context, asgName *string, instID *string, lifecycleHook *string, result *string) error {
	input := autoscaling.CompleteLifecycleActionInput{
		AutoScalingGroupName:  asgName,
		InstanceId:            instID,
		LifecycleActionResult: result,
		LifecycleHookName:     lifecycleHook,
	}

	_, err := c.ASGClient.CompleteLifecycleAction(ctx, &input)
	return errors.Wrapf(err, "error completing lifecycle hook %s for instance %s", *lifecycleHook, *instID)
}

// TerminateInstanceInASG calls https://docs.aws.amazon.com/cli/latest/reference/autoscaling/terminate-instance-in-auto-scaling-group.html
func (c *Clients) TerminateInstanceInASG(ctx context.Context, instID *string, decrement *bool) error {
	input := autoscaling.TerminateInstanceInAutoScalingGroupInput{
		InstanceId:                     instID,
		ShouldDecrementDesiredCapacity: decrement,
	}
	_, err := c.ASGClient.TerminateInstanceInAutoScalingGroup(ctx, &input)
	return errors.Wrapf(err, "error terminating instance %s", *instID)
}

// SetDesiredCapacity sets the desired capacity of given ASG to given value
func (c *Clients) SetDesiredCapacity(ctx context.Context, asg *at.AutoScalingGroup, desiredCapacity *int32) error {
	input := autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: asg.AutoScalingGroupName,
		DesiredCapacity:      desiredCapacity,
	}
	_, err := c.ASGClient.SetDesiredCapacity(ctx, &input)
	return errors.Wrapf(err, "error setting desired capacity for %s", *asg.AutoScalingGroupName)
}
