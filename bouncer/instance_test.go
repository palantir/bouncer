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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
)

// TestIsOldLogic mostly to make sure we don't panic because the logic games and everything in AWS being a pointer
func TestIsOldLogic(t *testing.T) {
	startTime := time.Now()
	var isOld bool
	var zeroTime time.Time
	var asgInst autoscaling.Instance

	ec2Inst := &ec2.Instance{
		LaunchTime: &zeroTime,
	}

	iid := "i-123456789abcdefgh"
	lts := &autoscaling.LaunchTemplateSpecification{
		LaunchTemplateId:   aws.String("lt-123456789abcdefgh"),
		LaunchTemplateName: aws.String("test-launch-template"),
		Version:            aws.String("1"),
	}

	// LT instance
	asgInst = autoscaling.Instance{
		LaunchConfigurationName: nil,
		InstanceId:              &iid,
		LaunchTemplate:          lts,
	}

	// old LT
	isOld = isInstanceOld(&asgInst, ec2Inst, nil, lts, aws.String("2"), false, startTime)
	assert.True(t, isOld)

	// not old LT
	isOld = isInstanceOld(&asgInst, ec2Inst, nil, lts, aws.String("1"), false, startTime)
	assert.False(t, isOld)

	// force it to be old
	isOld = isInstanceOld(&asgInst, ec2Inst, nil, lts, aws.String("1"), true, startTime)
	assert.True(t, isOld)

	// malformed ASG for LT instance that should otherwise not be old
	isOld = isInstanceOld(&asgInst, ec2Inst, nil, nil, aws.String("1"), false, startTime)
	assert.True(t, isOld)

	// LC Instance
	asgInst = autoscaling.Instance{
		LaunchConfigurationName: aws.String("hi-there-1"),
		InstanceId:              &iid,
		LaunchTemplate:          nil,
	}

	// old LC
	isOld = isInstanceOld(&asgInst, ec2Inst, aws.String("hi-there-2"), nil, nil, false, startTime)
	assert.True(t, isOld)

	// not old LC
	isOld = isInstanceOld(&asgInst, ec2Inst, aws.String("hi-there-1"), nil, nil, false, startTime)
	assert.False(t, isOld)

	// force it to be old
	isOld = isInstanceOld(&asgInst, ec2Inst, aws.String("hi-there-1"), nil, nil, true, startTime)
	assert.True(t, isOld)

	// malformed ASG for LC instance that should otherwise not be old
	isOld = isInstanceOld(&asgInst, ec2Inst, nil, nil, nil, false, startTime)
	assert.True(t, isOld)
}
