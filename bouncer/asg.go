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

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/palantir/bouncer/aws"
	"github.com/pkg/errors"
)

// ASG object holds a pointer to an ASG and its Instances
type ASG struct {
	ASG        *autoscaling.Group
	Instances  []*Instance
	DesiredASG *DesiredASG
}

// NewASG creates a new ASG object
func NewASG(ac *aws.Clients, desASG *DesiredASG, force bool, startTime time.Time) (*ASG, error) {
	var awsAsg *autoscaling.Group

	err := retry(apiRetryCount, apiRetrySleep, func() (err error) {
		awsAsg, err = ac.GetASG(&desASG.AsgName)
		return
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS ASG object")
	}

	var instances []*Instance

	for _, asgInst := range awsAsg.Instances {
		inst, err := NewInstance(ac, awsAsg, asgInst, force, startTime, desASG.PreTerminateCmd)
		if err != nil {
			return nil, errors.Wrapf(err, "error generating bouncer.instance for %s", *asgInst.InstanceId)
		}
		instances = append(instances, inst)
	}

	asg := ASG{
		ASG:        awsAsg,
		Instances:  instances,
		DesiredASG: desASG,
	}

	return &asg, nil
}
