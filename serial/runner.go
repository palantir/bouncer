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
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/palantir/bouncer/bouncer"
	"github.com/pkg/errors"
)

// Runner holds data for a particular serial run
type Runner struct {
	bouncer.BaseRunner
}

// NewRunner instantiates a new serial runner
func NewRunner(opts *bouncer.RunnerOpts) (*Runner, error) {
	br, err := bouncer.NewBaseRunner(opts)
	if err != nil {
		return nil, errors.Wrap(err, "error getting base runner")
	}

	r := Runner{
		*br,
	}
	return &r, nil
}

func (r *Runner) killBestOldInstance(asgSet *bouncer.ASGSet) error {
	bestOld := asgSet.GetBestOldInstance()
	err := r.KillInstance(bestOld)
	return errors.Wrap(err, "error killing instance")
}

// MustValidatePrereqs checks that the batch runner is safe to proceed
func (r *Runner) MustValidatePrereqs() {
	asgSet, err := r.NewASGSet()
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

		if *asg.ASG.DesiredCapacity == *asg.ASG.MinSize {
			log.WithFields(log.Fields{
				"ASG":              *asg.ASG.AutoScalingGroupName,
				"desired_capacity": *asg.ASG.DesiredCapacity,
				"min_size":         *asg.ASG.MinSize,
			}).Error("ASG desired capacity must be at least 1 higher than the min size, but they're equal")
			os.Exit(1)
		}
	}
}

// Run has the meat of the batch job
func (r *Runner) Run() error {
	for {
		if r.TimedOut() {
			return errors.Errorf("timeout exceeded, something is probably wrong with rollout")
		}

		// Rebuild the state of the world every iteration of the loop because instance and ASG statuses are changing
		log.Debug("Beginning new serial run check")
		asgSet, err := r.NewASGSet()
		if err != nil {
			return errors.Wrap(err, "error building ASGSet")
		}

		// See if we're still waiting on a change we made previously to finish or settle
		if asgSet.IsNewUnhealthy() || asgSet.IsTerminating() || asgSet.IsImmutableAutoscalingEvent() || asgSet.IsCountMismatch() {
			r.Sleep()
			continue
		}

		// See if anyone's desired capacity needs to be reset, and fix it if so (then sleep so it propagates)
		divergedASGs := asgSet.GetDivergedASGs()
		for _, asg := range divergedASGs {
			err := r.SetDesiredCapacity(asg, &asg.DesiredASG.DesiredCapacity)
			if err != nil {
				return errors.Wrap(err, "error setting desired capacity of ASG")
			}
		}

		if len(divergedASGs) != 0 {
			r.Sleep()
			continue
		}

		// If there are any old instances which are now ready to be terminated, let's do it
		if asgSet.IsOldInstance() {
			err = r.killBestOldInstance(asgSet)
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
