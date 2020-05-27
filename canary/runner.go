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

package canary

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/palantir/bouncer/bouncer"
	"github.com/pkg/errors"
)

// Runner holds data for a particular canary run
// Note that in the canary case, asgs will always be of length 1
type Runner struct {
	bouncer.BaseRunner
}

// NewRunner instantiates a new canary runner
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

// MustValidatePrereqs checks that the batch runner is safe to proceed
func (r *Runner) MustValidatePrereqs() {
	asgSet, err := r.NewASGSet()
	if err != nil {
		log.Fatal(errors.Wrap(err, "error building actualASG"))
	}

	if len(asgSet.ASGs) > 1 {
		log.WithFields(log.Fields{
			"count given": len(asgSet.ASGs),
		}).Error("Canary mode supports only 1 ASG at a time")
		os.Exit(1)
	}

	for _, actualAsg := range asgSet.ASGs {
		if actualAsg.DesiredASG.DesiredCapacity != *actualAsg.ASG.DesiredCapacity {
			log.WithFields(log.Fields{
				"ASG":                     *actualAsg.ASG.AutoScalingGroupName,
				"desired_capacity given":  actualAsg.DesiredASG.DesiredCapacity,
				"desired_capacity actual": *actualAsg.ASG.DesiredCapacity,
			}).Error("Desired capacity given must be equal to starting desired_capacity of ASG")
			os.Exit(1)
		}

		if actualAsg.DesiredASG.DesiredCapacity < *actualAsg.ASG.MinSize {
			log.WithFields(log.Fields{
				"ASG":              *actualAsg.ASG.AutoScalingGroupName,
				"min_size":         *actualAsg.ASG.MinSize,
				"max_size":         *actualAsg.ASG.MaxSize,
				"desired_capacity": actualAsg.DesiredASG.DesiredCapacity,
			}).Error("Desired capacity given must be greater than or equal to min ASG size")
			os.Exit(1)
		}

		if (actualAsg.DesiredASG.DesiredCapacity * 2) > *actualAsg.ASG.MaxSize {
			log.WithFields(log.Fields{
				"ASG":              *actualAsg.ASG.AutoScalingGroupName,
				"min_size":         *actualAsg.ASG.MinSize,
				"max_size":         *actualAsg.ASG.MaxSize,
				"desired_capacity": actualAsg.DesiredASG.DesiredCapacity,
			}).Error("Desired capacity given must be less than or equal to 2x max_size")
			os.Exit(1)
		}
	}
}

// Run has the meat of the batch job
func (r *Runner) Run() error {
	var newDesiredCapacity int64

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
		if asgSet.IsTerminating() || asgSet.IsNewUnhealthy() || asgSet.IsImmutableAutoscalingEvent() || asgSet.IsCountMismatch() {
			r.Sleep()
			continue
		}

		// Since we only support one ASG in canary mode
		asg := asgSet.ASGs[0]
		curDesiredCapacity := asg.ASG.DesiredCapacity
		finDesiredCapacity := &asg.DesiredASG.DesiredCapacity
		newCount := int64(len(asgSet.GetNewInstances()))
		oldCount := int64(len(asgSet.GetOldInstances()))

		if newCount == *finDesiredCapacity {
			if *curDesiredCapacity == *finDesiredCapacity {
				if oldCount == 0 {
					log.Info("Didn't find any old instances or ASGs - we're done here!")
					return nil
				}

				// Only wait for terminating instances to finish terminating once all
				// terminate commands have been issued
				if asgSet.IsTerminating() {
					r.Sleep()
					continue
				} else {
					log.WithFields(log.Fields{
						"ASG":           *asg.ASG.AutoScalingGroupName,
						"Old instances": oldCount,
						"New instances": newCount,
					}).Error("I have old instances which aren't terminating")
					return errors.New("old instance mismatch")
				}
			}

			// Don't think this condition should be reachable, but just in case
			if oldCount == 0 {
				log.WithFields(log.Fields{
					"ASG":                    *asg.ASG.AutoScalingGroupName,
					"Old instances":          oldCount,
					"New instances":          newCount,
					"Desired Capacity":       *curDesiredCapacity,
					"Final Desired Capacity": *finDesiredCapacity,
				}).Error("Somehow there are no old nodes but new count is off?")
				return errors.New("capacity mismatch")
			}

			// We have the correct number of new instances, we just need
			// to get rid of the old ones
			// Let's issue all their terminates right here
			for _, oldInst := range asgSet.GetOldInstances() {
				err := r.KillInstance(oldInst)
				if err != nil {
					return errors.Wrap(err, "error killing instance")
				}
				r.Sleep()
			}

			continue
		}

		// Not sure we'll ever hit the IsTerminating case here, we should only hit that inside above if-block
		// The IsCountMismatch check is here and not where IsNewUnhealthy is, because we don't want it
		// to fire when bad nodes are in the process of terminating, since we issue terminates to them one at a time
		if asgSet.IsTerminating() || asgSet.IsCountMismatch() {
			r.Sleep()
			continue
		}

		if oldCount == 0 {
			badCounts := asgSet.GetDivergedASGs()
			if len(badCounts) != 0 {
				return errors.New("somehow our ASG's desired count isn't the canonical count, but we have all new instances, if this is correct, manually set desired capacity")
			}
		}

		if newCount == 0 {
			// We haven't canaried a new instance yet, so let's do that
			log.WithFields(log.Fields{
				"ASG": *asg.ASG.AutoScalingGroupName,
			}).Info("Adding canary node")
			newDesiredCapacity = *curDesiredCapacity + 1
		} else {
			// Otherwise, we've already canaried successfully, so let's expand out to full size
			log.WithFields(log.Fields{
				"ASG": *asg.ASG.AutoScalingGroupName,
			}).Info("Adding in remainder of new nodes")
			// Just set des cap to be current + the number of new nodes that we're short
			newDesiredCapacity = *curDesiredCapacity + (*finDesiredCapacity - newCount)
		}

		err = r.SetDesiredCapacity(asg, &newDesiredCapacity)
		if err != nil {
			return errors.Wrap(err, "error setting desired capacity")
		}

		r.Sleep()
		continue
	}
}
