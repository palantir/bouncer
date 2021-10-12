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

	at "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	et "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"
)

// GetEC2TagValue returns a pointer to the value of the tag with the given key
func GetEC2TagValue(ec2 *et.Instance, key string) *string {
	for _, tag := range ec2.Tags {
		if tag.Key != nil {
			if strings.ToLower(*tag.Key) == strings.ToLower(key) {
				return tag.Value
			}
		}
	}
	return nil
}

// GetUserData returns a pointer to the value of the instance's userdata
func (c *Clients) GetUserData(ctx context.Context, inst *at.Instance) (*string, error) {
	input := ec2.DescribeInstanceAttributeInput{
		Attribute:  et.InstanceAttributeNameUserData,
		InstanceId: inst.InstanceId,
	}

	output, err := c.EC2Client.DescribeInstanceAttribute(ctx, &input)
	if err != nil {
		return nil, errors.Wrapf(err, "Error describing userdata for instance %s", *inst.InstanceId)
	}

	return output.UserData.Value, nil
}
