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
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	at "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	et "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"
)

const apiSleepTime = 200 * time.Millisecond

// Clients holds the clients for this account's invocation of the APIs we'll need
type Clients struct {
	ctx       context.Context
	ASGClient *autoscaling.Client
	EC2Client *ec2.Client
}

// GetAWSClients returns the AWS client objects we'll need
func GetAWSClients() (*Clients, error) {
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, errors.Wrap(err, "Error opening default AWS config")
	}

	asg := autoscaling.NewFromConfig(cfg)
	ec2 := ec2.NewFromConfig(cfg)

	ac := Clients{
		ctx:       ctx,
		ASGClient: asg,
		EC2Client: ec2,
	}

	return &ac, nil
}

// ASGInstToEC2Inst converts a *autoscaling.Instance to its corresponding *ec2.Instance
func (c *Clients) ASGInstToEC2Inst(inst at.Instance) (*et.Instance, error) {
	input := ec2.DescribeInstancesInput{
		InstanceIds: []string{*inst.InstanceId},
	}
	output, err := c.EC2Client.DescribeInstances(c.ctx, &input)
	if err != nil {
		return nil, errors.Wrapf(err, "Error describing instance %s", *inst.InstanceId)
	}

	for _, res := range output.Reservations {
		if len(res.Instances) > 1 {
			return nil, errors.New("More than 1 instance found somehow")
		}

		for _, ec2Inst := range res.Instances {
			return &ec2Inst, nil
		}
	}

	return nil, errors.Errorf("No instances found for %s", *inst.InstanceId)
}

// ASGLTplVersionToEC2LTplVersion resolves ASG Template Versions to its actual *int32 ec2LaunchTemplate Version
func (c Clients) ASGLTplVersionToEC2LTplVersion(asgLaunchTemplate *at.LaunchTemplateSpecification) (*string, error) {
	// No launch template, nothing to do here
	if asgLaunchTemplate == nil {
		return nil, nil
	}

	input := &ec2.DescribeLaunchTemplatesInput{
		LaunchTemplateIds: []string{
			*asgLaunchTemplate.LaunchTemplateId,
		},
	}

	res, err := c.EC2Client.DescribeLaunchTemplates(c.ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "Error describing LaunchTemplate %s", *asgLaunchTemplate.LaunchTemplateId)
	}

	if len(res.LaunchTemplates) != 1 {
		return nil, errors.Wrapf(err,
			"Expected exactly one LaunchTemplate returned for launch template id %s, got %d: %v",
			*asgLaunchTemplate.LaunchTemplateId,
			len(res.LaunchTemplates),
			res.LaunchTemplates,
		)
	}

	ec2LaunchTemplate := res.LaunchTemplates[0]

	targetVersion := asgLaunchTemplate.Version

	// Per https://docs.aws.amazon.com/autoscaling/ec2/APIReference/API_LaunchTemplateSpecification.html
	// version is optional and if unspecified should resolve to default.
	if targetVersion == nil || *targetVersion == "$Default" {
		s := strconv.FormatInt(*ec2LaunchTemplate.DefaultVersionNumber, 10)
		return &s, nil
	} else if *targetVersion == "$Latest" {
		s := strconv.FormatInt(*ec2LaunchTemplate.LatestVersionNumber, 10)
		return &s, nil
	} else {
		_, err := strconv.ParseInt(*targetVersion, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Unexpected TemplateVersion %q conversion to Int64 failed", *targetVersion)
		}
		return targetVersion, nil
	}

}
