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

package slowcanary

import (
	"os"

	"github.com/palantir/bouncer/bouncer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Runner holds data for a particular slow-canary run
// Note that in the slow-canary case, asgs will always be of length 1
type Runner struct {
	bouncer.BaseRunner
}

// NewRunner instantiates a new slow-canary runner
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
		}).Error("Slow-canary mode supports only 1 ASG at a time")
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

		if (actualAsg.DesiredASG.DesiredCapacity + 1) > *actualAsg.ASG.MaxSize {
			log.WithFields(log.Fields{
				"ASG":              *actualAsg.ASG.AutoScalingGroupName,
				"min_size":         *actualAsg.ASG.MinSize,
				"max_size":         *actualAsg.ASG.MaxSize,
				"desired_capacity": actualAsg.DesiredASG.DesiredCapacity,
			}).Error("Max capacity set on ASG must be at least 1 + desired")
			os.Exit(1)
		}
	}
}

// Run has the meat of the batch job
func (r *Runner) Run() error {
	var newDesiredCapacity int32

	for {
		if r.TimedOut() {
			return errors.Errorf("timeout exceeded, something is probably wrong with rollout")
		}

		// Rebuild the state of the world every iteration of the loop because instance and ASG statuses are changing
		log.Debug("Beginning new slow-canary run check")
		asgSet, err := r.NewASGSet()
		if err != nil {
			return errors.Wrap(err, "error building ASGSet")
		}

		if asgSet.IsTransient() {
			r.Sleep()
			continue
		}

		// Since we only support one ASG in slow-canary mode
		asg := asgSet.ASGs[0]
		curDesiredCapacity := asg.ASG.DesiredCapacity
		finDesiredCapacity := &asg.DesiredASG.DesiredCapacity
		newCount := int32(len(asgSet.GetNewInstances()))
		oldCount := int32(len(asgSet.GetOldInstances()))

		if *curDesiredCapacity == *finDesiredCapacity {
			if oldCount == 0 {
				log.Info("Didn't find any old instances or ASGs - we're done here!")
				return nil
			}

			log.WithFields(log.Fields{
				"ASG": *asg.ASG.AutoScalingGroupName,
			}).Info("Adding slow-canary node")
			newDesiredCapacity = *curDesiredCapacity + 1

			err = r.SetDesiredCapacity(asg, &newDesiredCapacity)
			if err != nil {
				return errors.Wrap(err, "error setting desired capacity")
			}

			r.Sleep()
			continue
		} else if *curDesiredCapacity == *finDesiredCapacity+1 {
			if oldCount == 0 {
				log.WithFields(log.Fields{
					"ASG":                    *asg.ASG.AutoScalingGroupName,
					"Old instances":          oldCount,
					"New instances":          newCount,
					"Desired Capacity":       *curDesiredCapacity,
					"Final Desired Capacity": *finDesiredCapacity,
				}).Error("Somehow there are no old nodes but capacities are mismatched?")
				return errors.New("capacity mismatch")
			}

			if oldCount == 1 {
				// Kill our last old instance, decrementing our capacity back to our desired value
				log.WithFields(log.Fields{
					"ASG": *asg.ASG.AutoScalingGroupName,
				}).Info("Killing the last old node, so not letting AWS replace it")
				decrement := true
				oldInstances := asgSet.GetOldInstances()
				err := r.KillInstance(oldInstances[0], &decrement)
				if err != nil {
					return errors.Wrap(err, "error killing instance")
				}
				r.Sleep()

				continue
			}

			// Otherwise, we still have more than 1 old instance, so let's terminate w/ replace
			log.WithFields(log.Fields{
				"ASG": *asg.ASG.AutoScalingGroupName,
			}).Info("Killing an old node, and letting AWS replace it")
			decrement := false
			oldInstances := asgSet.GetOldInstances()
			err := r.KillInstance(oldInstances[0], &decrement)
			if err != nil {
				return errors.Wrap(err, "error killing instance")
			}
			r.Sleep()

			continue
		}

		// Don't think this condition should be reachable, but just in case
		log.WithFields(log.Fields{
			"ASG":                    *asg.ASG.AutoScalingGroupName,
			"Old instances":          oldCount,
			"New instances":          newCount,
			"Desired Capacity":       *curDesiredCapacity,
			"Final Desired Capacity": *finDesiredCapacity,
		}).Error("Found capacity mismatch")
		return errors.New("capacity mismatch")
	}
}
