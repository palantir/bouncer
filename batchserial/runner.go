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

package batchserial

import (
	"context"

	at "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/palantir/bouncer/bouncer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Runner holds data for a particular batch-serial run
// Note that in the batch-serial case, asgs will always be of length 1
type Runner struct {
	bouncer.BaseRunner
	batchSize int32 // This field is set in ValidatePrereqs
}

// NewRunner instantiates a new batch-serial runner
func NewRunner(ctx context.Context, opts *bouncer.RunnerOpts) (*Runner, error) {
	br, err := bouncer.NewBaseRunner(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "error getting base runner")
	}

	r := Runner{
		BaseRunner: *br,
		batchSize:  *opts.BatchSize,
	}
	return &r, nil
}

// ValidatePrereqs checks that the batch runner is safe to proceed
func (r *Runner) ValidatePrereqs(ctx context.Context) error {
	asgSet, err := r.NewASGSet(ctx)
	if err != nil {
		return errors.Wrap(err, "error building actualASG")
	}

	if len(asgSet.ASGs) > 1 {
		log.WithFields(log.Fields{
			"count given": len(asgSet.ASGs),
		}).Error("Batch Serial mode supports only 1 ASG at a time")
		return errors.New("error validating ASG input")
	}

	for _, actualAsg := range asgSet.ASGs {
		if actualAsg.DesiredASG.DesiredCapacity != *actualAsg.ASG.DesiredCapacity {
			log.WithFields(log.Fields{
				"desired capacity given":  actualAsg.DesiredASG.DesiredCapacity,
				"desired capacity actual": *actualAsg.ASG.DesiredCapacity,
			}).Error("Desired capacity given must be equal to starting desired_capacity of ASG")
			return errors.New("error validating ASG state")
		}

		if actualAsg.DesiredASG.DesiredCapacity < *actualAsg.ASG.MinSize {
			log.WithFields(log.Fields{
				"min size":         *actualAsg.ASG.MinSize,
				"max size":         *actualAsg.ASG.MaxSize,
				"desired capacity": actualAsg.DesiredASG.DesiredCapacity,
			}).Error("Desired capacity given must be greater than or equal to min ASG size")
			return errors.New("error validating ASG state")
		}

		if *actualAsg.ASG.MinSize > (actualAsg.DesiredASG.DesiredCapacity - r.batchSize) {
			log.WithFields(log.Fields{
				"min size":         *actualAsg.ASG.MinSize,
				"max size":         *actualAsg.ASG.MaxSize,
				"desired capacity": actualAsg.DesiredASG.DesiredCapacity,
				"batch size":       r.batchSize,
			}).Error("Min capacity of ASG must be <= desired capacity - batch size")
			return errors.New("error validating ASG state")
		}
	}

	return nil
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// Run has the meat of the batch job
func (r *Runner) Run() error {
	decrement := true

	ctx, cancel := r.NewContext()
	defer cancel()

	for {
		// Rebuild the state of the world every iteration of the loop because instance and ASG statuses are changing
		log.Debug("Beginning new batch serial run check")
		asgSet, err := r.NewASGSet(ctx)
		if err != nil {
			return errors.Wrap(err, "error building ASGSet")
		}

		// Since we only support one ASG in batch-serial mode
		asg := asgSet.ASGs[0]
		curDesiredCapacity := *asg.ASG.DesiredCapacity
		finDesiredCapacity := asg.DesiredASG.DesiredCapacity

		oldUnhealthy := asgSet.GetUnHealthyOldInstances()
		newHealthy := asgSet.GetHealthyNewInstances()
		oldHealthy := asgSet.GetHealthyOldInstances()

		newCount := int32(len(asgSet.GetNewInstances()))
		oldCount := int32(len(asgSet.GetOldInstances()))
		totalCount := newCount + oldCount

		healthyCount := int32(len(oldHealthy) + len(newHealthy))

		// Never terminate nodes so that we go below finDesiredCapacity - batchSize number of healthy (InService) machines
		minDesiredCapacity := finDesiredCapacity - r.batchSize
		toKill := min(finDesiredCapacity-minDesiredCapacity, oldCount)

		// Clean-out old unhealthy instances in P:W now, as they're just wasting time
		for _, oi := range oldUnhealthy {
			if oi.ASGInstance.LifecycleState == at.LifecycleStatePendingWait {
				err := r.KillInstance(ctx, oi, &decrement)
				if err != nil {
					return errors.Wrap(err, "error killing instance")
				}

				ctx, cancel = r.NewContext()
				defer cancel()
				r.Sleep(ctx)

				continue
			}
		}

		// This check already prints statuses of individual nodes
		if asgSet.IsTransient() {
			log.Info("Waiting for nodes to settle")
			r.Sleep(ctx)
			continue
		}

		// Our exit case - we have exactly the number of nodes we want, they're all new, and they're all InService
		if oldCount == 0 && totalCount == finDesiredCapacity {
			if curDesiredCapacity == finDesiredCapacity {
				log.Info("Didn't find any old instances or ASGs - we're done here!")
				return nil
			}

			// Not sure how this would happen off-hand?
			log.WithFields(log.Fields{
				"Current desired capacity": curDesiredCapacity,
				"Final desired capacity":   finDesiredCapacity,
			}).Error("Capacity mismatch")
			return errors.New("old instance mismatch")
		}

		// If we haven't done the termination piece of canary, let's do that
		if newCount == 0 && totalCount == finDesiredCapacity {
			log.Info("Terminating a canary node")
			oi := asgSet.GetBestOldInstance()

			err := r.KillInstance(ctx, oi, &decrement)
			if err != nil {
				return errors.Wrap(err, "error killing instance")
			}

			ctx, cancel = r.NewContext()
			defer cancel()
			r.Sleep(ctx)

			continue
		}

		// If we haven't done the new instance piece of canary, let's do that
		if newCount == 0 && totalCount < finDesiredCapacity {
			err = r.SetDesiredCapacity(ctx, asg, &finDesiredCapacity)
			if err != nil {
				return errors.Wrap(err, "error setting desired capacity")
			}

			ctx, cancel = r.NewContext()
			defer cancel()
			r.Sleep(ctx)

			continue
		}

		// Scale-in a batch
		if totalCount == finDesiredCapacity && toKill > 0 {
			killed := int32(0)

			log.WithFields(log.Fields{
				"Old nodes":     oldCount,
				"Healthy nodes": healthyCount,
				"Nodes to kill": toKill,
			}).Info("Killing a batch of nodes")

			for _, oi := range oldHealthy {
				err := r.KillInstance(ctx, oi, &decrement)
				if err != nil {
					return errors.Wrap(err, "error killing instance")
				}
				killed++
				if killed == toKill {
					log.WithFields(log.Fields{
						"Killed Nodes": killed,
					}).Info("Already killed max number of nodes to get to min capacity, pausing here")
					break
				}
			}

			ctx, cancel = r.NewContext()
			defer cancel()
			r.Sleep(ctx)

			continue
		}

		// Scale-out a batch to original size to refresh nodes
		if totalCount < finDesiredCapacity {
			err = r.SetDesiredCapacity(ctx, asg, &finDesiredCapacity)
			if err != nil {
				return errors.Wrap(err, "error setting desired capacity")
			}

			ctx, cancel = r.NewContext()
			defer cancel()
			r.Sleep(ctx)

			continue
		}

		// Not sure how this would happen off-hand?
		log.WithFields(log.Fields{
			"Current desired capacity": curDesiredCapacity,
			"Final desired capacity":   finDesiredCapacity,
			"Old nodes":                oldCount,
			"Healthy nodes":            healthyCount,
			"Nodes To Kill":            toKill,
		}).Error("Unknown condition hit")
		return errors.New("undefined error")
	}
}
