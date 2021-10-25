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
	"context"

	"github.com/palantir/bouncer/bouncer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Runner holds data for a particular full run
// Note that in the full case, asgs will always be of length 1
type Runner struct {
	bouncer.BaseRunner
}

// NewRunner instantiates a new full runner
func NewRunner(ctx context.Context, opts *bouncer.RunnerOpts) (*Runner, error) {
	br, err := bouncer.NewBaseRunner(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "error getting base runner")
	}

	r := Runner{
		*br,
	}
	return &r, nil
}

// ValidatePrereqs checks that the batch runner is safe to proceed
func (r *Runner) ValidatePrereqs(ctx context.Context) error {
	asgSet, err := r.NewASGSet(ctx)
	if err != nil {
		return errors.Wrap(err, "error building ASGSet")
	}

	divergedASGs := asgSet.GetDivergedASGs()
	if len(divergedASGs) != 0 {
		for _, badASG := range divergedASGs {
			log.WithFields(log.Fields{
				"ASG":                     *badASG.ASG.AutoScalingGroupName,
				"desired_capacity actual": *badASG.ASG.DesiredCapacity,
				"desired_capacity given":  badASG.DesiredASG.DesiredCapacity,
			}).Error("ASG desired capacity doesn't match expected starting value")
		}
		return errors.New("error validating initial ASG state")
	}

	for _, asg := range asgSet.ASGs {
		if *asg.ASG.DesiredCapacity == 0 {
			log.WithFields(log.Fields{
				"ASG": *asg.ASG.AutoScalingGroupName,
			}).Warn("ASG desired capacity is 0 - nothing to do here")
			return errors.New("error validating initial ASG state")
		}

		if *asg.ASG.MinSize != 0 {
			log.WithFields(log.Fields{
				"ASG":      *asg.ASG.AutoScalingGroupName,
				"min_size": *asg.ASG.MinSize,
			}).Error("ASG min size must equal 0")
			return errors.New("error validating initial ASG state")
		}
	}

	return nil
}

func reverseASGSetOrder(asg []*bouncer.ASG) []*bouncer.ASG {
	// copy to new slice
	new := append(asg[:0:0], asg...)

	// reverse order of new slice
	for i := len(new)/2 - 1; i >= 0; i-- {
		rev := len(new) - 1 - i
		new[i], new[rev] = new[rev], new[i]
	}

	return new
}

func asgSetWrapper(asg *bouncer.ASG) *bouncer.ASGSet {
	return &bouncer.ASGSet{
		ASGs: []*bouncer.ASG{asg},
	}
}

// Run has the meat of the batch job
func (r *Runner) Run() error {
	var newDesiredCapacity int32

	ctx, cancel := r.NewContext()
	defer cancel()

start:
	for {
		// Rebuild the state of the world every iteration of the loop because instance and ASG statuses are changing
		log.Debug("Beginning new full run check")
		asgSet, err := r.NewASGSet(ctx)
		if err != nil {
			return errors.Wrap(err, "error building ASGSet")
		}

		// See if we're still waiting on a change we made previously to finish or settle
		if asgSet.IsTransient() {
			r.Sleep(ctx)
			continue
		}

		// drain one ASG at a time one instance at a time until no ASGs have any old instances
		for _, asg := range asgSet.ASGs {
			set := asgSetWrapper(asg)

			if set.IsOldInstance() {
				decrement := true
				err := r.KillInstance(ctx, set.GetBestOldInstance(), &decrement)
				if err != nil {
					return errors.Wrap(err, "failed to kill instance")
				}

				ctx, cancel = r.NewContext()
				defer cancel()
				r.Sleep(ctx)

				continue start
			}
		}

		// restore ASG's in reversed order
		for _, asg := range reverseASGSetOrder(asgSet.ASGs) {
			// restore one instance at a time until back to desired cap
			if *asg.ASG.DesiredCapacity < asg.DesiredASG.DesiredCapacity {
				newDesiredCapacity = *asg.ASG.DesiredCapacity + 1

				err = r.SetDesiredCapacity(ctx, asg, &newDesiredCapacity)
				if err != nil {
					return errors.Wrap(err, "error setting desired capacity")
				}

				ctx, cancel = r.NewContext()
				defer cancel()
				r.Sleep(ctx)

				continue start
			}
		}

		return nil
	}
}
