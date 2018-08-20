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
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

const apiSleepTime = 200 * time.Millisecond

// Clients holds the clients for this account's invocation of the APIs we'll need
type Clients struct {
	ASGClient *autoscaling.AutoScaling
	EC2Client *ec2.EC2
}

// GetAWSClients returns the AWS client objects we'll need
func GetAWSClients() (*Clients, error) {
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	awsConf := aws.Config{
		Region: &region,
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: awsConf,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error opening AWS session")
	}

	asg := autoscaling.New(sess)
	ec2 := ec2.New(sess)

	ac := Clients{
		ASGClient: asg,
		EC2Client: ec2,
	}

	return &ac, nil
}

// ASGInstToEC2Inst converts a *autoscaling.Instance to its corresponding *ec2.Instance
func (c *Clients) ASGInstToEC2Inst(inst *autoscaling.Instance) (*ec2.Instance, error) {
	var instIDs []*string
	instIDs = append(instIDs, inst.InstanceId)
	input := ec2.DescribeInstancesInput{
		InstanceIds: instIDs,
	}
	output, err := c.EC2Client.DescribeInstances(&input)
	if err != nil {
		return nil, errors.Wrapf(err, "Error describing instance %s", *inst.InstanceId)
	}

	for _, res := range output.Reservations {
		if len(res.Instances) > 1 {
			return nil, errors.New("More than 1 instance found somehow")
		}

		for _, ec2Inst := range res.Instances {
			return ec2Inst, nil
		}
	}

	return nil, errors.Errorf("No instances found for %s", *inst.InstanceId)
}
