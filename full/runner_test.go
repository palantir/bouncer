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

package full

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/palantir/bouncer/bouncer"
	"github.com/stretchr/testify/assert"
)

func asgSliceTestConstructor(len int) []*bouncer.ASG {
	var asgs []*bouncer.ASG

	for i := 0; i < len; i++ {
		name := fmt.Sprintf("asg-%v", i)
		asgs = append(asgs, &bouncer.ASG{
			ASG: &autoscaling.Group{
				AutoScalingGroupName: &name,
			},
		})
	}

	return asgs
}

func TestReverseASGSlice(t *testing.T) {
	asgs := asgSliceTestConstructor(5)
	rev := reverseASGSetOrder(asgs)

	assert.NotEqual(t, asgs, rev, "The original asg slice should NOT equal the reversed")
	assert.Equal(t, asgs, asgSliceTestConstructor(5), "The original asg slice should still equal a fresh slice with the same number of instances, even after being reversed")
	assert.Equal(t, asgs, reverseASGSetOrder(rev), "Reversing the already reversed asg set, should equal the original")

	// tests on varying slice length
	assert.Equal(t, asgSliceTestConstructor(1), reverseASGSetOrder(asgSliceTestConstructor(1)), "An asg slice of 1 does not change when reversed")
	assert.Equal(t, asgSliceTestConstructor(3)[0], reverseASGSetOrder(asgSliceTestConstructor(3))[2], "The first asg in a set of 3 should equal the last asg in the reverse")
	assert.Equal(t, asgSliceTestConstructor(4)[0], reverseASGSetOrder(asgSliceTestConstructor(4))[3], "The first asg in a set of 4 should equal the last asg in the reverse")
	assert.Equal(t, asgSliceTestConstructor(7)[6], reverseASGSetOrder(asgSliceTestConstructor(7))[0], "The last asg in a set of 7 should equal the first entry in the reverse")
	assert.Equal(t, asgSliceTestConstructor(8)[7], reverseASGSetOrder(asgSliceTestConstructor(8))[0], "The last asg in a set of 8 should equal the first entry in the reverse")
	assert.Equal(t, asgSliceTestConstructor(15)[7], reverseASGSetOrder(asgSliceTestConstructor(15))[7], "The middle asg in a set of 15 should equal the middle asg in the reverse as it is an odd numbered slice")
}
