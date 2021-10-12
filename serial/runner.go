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

package serial

import (
	"context"

	"github.com/palantir/bouncer/bouncer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Runner holds data for a particular serial run
type Runner struct {
	bouncer.BaseRunner
}

// NewRunner instantiates a new serial runner
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

func (r *Runner) killBestOldInstance(ctx context.Context, asgSet *bouncer.ASGSet) error {
	bestOld := asgSet.GetBestOldInstance()
	decrement := true
	err := r.KillInstance(ctx, bestOld, &decrement)
	return errors.Wrap(err, "error killing instance")
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

		if *asg.ASG.DesiredCapacity == *asg.ASG.MinSize {
			log.WithFields(log.Fields{
				"ASG":              *asg.ASG.AutoScalingGroupName,
				"desired_capacity": *asg.ASG.DesiredCapacity,
				"min_size":         *asg.ASG.MinSize,
			}).Error("ASG desired capacity must be at least 1 higher than the min size, but they're equal")
			return errors.New("error validating initial ASG state")
		}
	}

	return nil
}

// Run has the meat of the batch job
func (r *Runner) Run(ctx context.Context) error {
	ctx, cancel := r.NewContext(ctx)
	defer cancel()

	for {
		// Rebuild the state of the world every iteration of the loop because instance and ASG statuses are changing
		log.Debug("Beginning new serial run check")
		asgSet, err := r.NewASGSet(ctx)
		if err != nil {
			return errors.Wrap(err, "error building ASGSet")
		}

		// See if we're still waiting on a change we made previously to finish or settle
		if asgSet.IsTransient() {
			r.Sleep(ctx)
			continue
		}

		// See if anyone's desired capacity needs to be reset, and fix it if so (then sleep so it propagates)
		divergedASGs := asgSet.GetDivergedASGs()
		for _, asg := range divergedASGs {
			err := r.SetDesiredCapacity(ctx, asg, &asg.DesiredASG.DesiredCapacity)
			if err != nil {
				return errors.Wrap(err, "error setting desired capacity of ASG")
			}
		}

		if len(divergedASGs) != 0 {
			ctx, cancel = r.ResetAndSleep(ctx)
			defer cancel()

			continue
		}

		// If there are any old instances which are now ready to be terminated, let's do it
		if asgSet.IsOldInstance() {
			err = r.killBestOldInstance(ctx, asgSet)
			if err != nil {
				return errors.Wrap(err, "error finding or killing best old instance")
			}

			ctx, cancel = r.ResetAndSleep(ctx)
			defer cancel()

			continue
		}

		log.Info("Didn't find any old instances or ASGs - we're done here!")
		return nil
	}
}
