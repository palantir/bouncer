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

package rolling

import (
	"context"
	"os"

	"github.com/palantir/bouncer/bouncer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Runner holds data for a particular rolling run
type Runner struct {
	bouncer.BaseRunner
}

// NewRunner instantiates a new rolling runner
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
	decrement := false
	err := r.KillInstance(ctx, bestOld, &decrement)
	return errors.Wrap(err, "error killing instance")
}

// MustValidatePrereqs checks that the batch runner is safe to proceed
func (r *Runner) MustValidatePrereqs(ctx context.Context) {
	asgSet, err := r.NewASGSet(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error building ASGSet"))
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
		os.Exit(1)
	}

	for _, asg := range asgSet.ASGs {
		if *asg.ASG.DesiredCapacity == 0 {
			log.WithFields(log.Fields{
				"ASG": *asg.ASG.AutoScalingGroupName,
			}).Warn("ASG desired capacity is 0 - nothing to do here")
			os.Exit(0)
		}
	}
}

// Run has the meat of the batch job
func (r *Runner) Run(ctx context.Context) error {
	for {
		if r.TimedOut() {
			return errors.Errorf("timeout exceeded, something is probably wrong with rollout")
		}

		// Rebuild the state of the world every iteration of the loop because instance and ASG statuses are changing
		log.Debug("Beginning new rolling run check")
		asgSet, err := r.NewASGSet(ctx)
		if err != nil {
			return errors.Wrap(err, "error building ASGSet")
		}

		// See if we're still waiting on a change we made previously to finish or settle
		if asgSet.IsTransient() {
			r.Sleep()
			continue
		}

		// If there are any old instances which are now ready to be terminated, let's do it
		if asgSet.IsOldInstance() {
			err = r.killBestOldInstance(ctx, asgSet)
			if err != nil {
				return errors.Wrap(err, "error finding or killing best old instance")
			}

			r.Sleep()
			continue
		}

		log.Info("Didn't find any old instances or ASGs - we're done here!")
		return nil
	}
}
