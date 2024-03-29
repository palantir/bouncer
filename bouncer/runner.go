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
	"context"
	"os"
	"strings"
	"time"

	at "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/palantir/bouncer/aws"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// RunnerOpts is user-supplied options to any flavor of runner
type RunnerOpts struct {
	Noop            bool
	Force           bool
	Fast            bool
	BatchSize       *int32
	AsgString       string
	CommandString   string
	DefaultCapacity *int32
	TerminateHook   string
	PendingHook     string
	ItemTimeout     time.Duration
}

// BaseRunner is the base struct for any runner
type BaseRunner struct {
	Opts       *RunnerOpts
	startTime  time.Time
	awsClients *aws.Clients
	asgs       []*DesiredASG
}

const (
	waitBetweenChecks = 15 * time.Second

	asgSeparator = ","

	debugTimeFormat = "2006-01-02 15:04:05 MST"
)

// NewBaseRunner instantiates a BaseRunner
func NewBaseRunner(ctx context.Context, opts *RunnerOpts) (*BaseRunner, error) {
	awsClients, err := aws.GetAWSClients(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting AWS Creds")
	}

	asgs, err := getASGList(opts)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing ASG list")
	}

	r := BaseRunner{
		Opts:       opts,
		startTime:  time.Now(),
		awsClients: awsClients,
		asgs:       asgs,
	}

	return &r, nil
}

func getASGList(opts *RunnerOpts) ([]*DesiredASG, error) {
	var asgs []*DesiredASG
	var cmdStringItems []string

	asgStringItems := strings.Split(opts.AsgString, asgSeparator)
	if opts.CommandString != "" {
		cmdStringItems = strings.Split(opts.CommandString, asgSeparator)

		if len(cmdStringItems) != len(asgStringItems) {
			return nil, errors.Errorf("You've provided %v ASGs, but %v external commands, counts must match", len(asgStringItems), len(cmdStringItems))
		}
	}

	for i, asgItem := range asgStringItems {
		var command *string

		if len(cmdStringItems) > 0 {
			command = &cmdStringItems[i]
		} else {
			command = nil
		}

		curAsg, err := ExtractDesiredASG(asgItem, opts.DefaultCapacity, command)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing ASG item")
		}
		asgs = append(asgs, curAsg)
	}

	return asgs, nil
}

func (r *BaseRunner) noopCheck() {
	if r.Opts.Noop {
		log.Warn("NOOP only - not actually performing previous action, and exiting script with success")
		os.Exit(0)
	}
}

func (r *BaseRunner) abandonLifecycle(ctx context.Context, inst *Instance, hook *string) error {
	log.WithFields(log.Fields{
		"InstanceID":     *inst.ASGInstance.InstanceId,
		"Hook":           *hook,
		"LifecycleState": inst.ASGInstance.LifecycleState,
	}).Warn("Issuing ABANDON to hook instead of terminating")
	result := "ABANDON"
	r.noopCheck()
	err := r.awsClients.CompleteLifecycleAction(ctx, inst.AutoscalingGroup.AutoScalingGroupName, inst.ASGInstance.InstanceId, hook, &result)
	return errors.Wrap(err, "error completing lifecycle action")
}

// KillInstance calls TerminateInstanceInAutoscalingGroup, or, if the instance is stuck
// in a lifecycle hook, issues an ABANDON to it, killing it more forcefully
func (r *BaseRunner) KillInstance(ctx context.Context, inst *Instance, decrement *bool) error {
	log.WithFields(log.Fields{
		"ASG":        *inst.AutoscalingGroup.AutoScalingGroupName,
		"InstanceID": *inst.ASGInstance.InstanceId,
	}).Info("Picked instance to die next")
	var hook string

	if inst.ASGInstance.LifecycleState == at.LifecycleStatePendingWait {
		hook = r.Opts.PendingHook
	}

	if inst.ASGInstance.LifecycleState == at.LifecycleStateTerminatingWait {
		hook = r.Opts.TerminateHook
	}

	if hook != "" {
		err := r.abandonLifecycle(ctx, inst, &hook)
		return errors.Wrapf(err, "error abandoning hook %s", hook)
	}

	if inst.PreTerminateCmd != nil {
		err := r.executeExternalCommand(ctx, *inst.PreTerminateCmd)
		if err != nil {
			return errors.Wrap(err, "error executing pre-terminate command")
		}
	}
	err := r.terminateInstanceInASG(ctx, inst, decrement)
	return errors.Wrap(err, "error terminating instance")
}

func (r *BaseRunner) terminateInstanceInASG(ctx context.Context, inst *Instance, decrement *bool) error {
	log.WithFields(log.Fields{
		"ASG":        *inst.AutoscalingGroup.AutoScalingGroupName,
		"InstanceID": *inst.ASGInstance.InstanceId,
	}).Info("Terminating instance")
	r.noopCheck()

	err := r.awsClients.TerminateInstanceInASG(ctx, inst.ASGInstance.InstanceId, decrement)

	return err
}

// SetDesiredCapacity Updates desired capacity of ASG
// This function should only be used to increase desired cap, not decrease, since AWS
// will _always_ remove instances based on AZ before any other criteria
// http://docs.aws.amazon.com/autoscaling/latest/userguide/as-instance-termination.html
func (r *BaseRunner) SetDesiredCapacity(ctx context.Context, asg *ASG, desiredCapacity *int32) error {

	log.WithFields(log.Fields{
		"ASG":           *asg.ASG.AutoScalingGroupName,
		"CurDesiredCap": *asg.ASG.DesiredCapacity,
		"NewDesiredCap": *desiredCapacity,
	}).Info("Changing desired capacity")
	r.noopCheck()

	err := r.awsClients.SetDesiredCapacity(ctx, asg.ASG, desiredCapacity)

	return errors.Wrapf(err, "error setting desired capacity of ASG")
}

// NewContext generates a context with the ItemTimeout from the parent context given
func (r *BaseRunner) NewContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), r.Opts.ItemTimeout)
	dn, _ := ctx.Deadline()

	l := log.WithFields(log.Fields{
		"Context deadline": dn.Format(debugTimeFormat),
		"Current time":     getHumanCurrentTime(),
	})

	l.Debug("Generating fresh context")

	return ctx, cancel
}

func getHumanCurrentTime() string {
	return time.Now().Format(debugTimeFormat)
}

// Sleep makes us sleep for the constant time - call this when waiting for an AWS change
func (r *BaseRunner) Sleep(ctx context.Context) {
	l := log.WithFields(log.Fields{
		"Sleep Duration": waitBetweenChecks,
		"Current time":   getHumanCurrentTime(),
	})

	l.Debug("Sleeping between checks")

	select {
	case <-time.After(waitBetweenChecks):
		return
	case <-ctx.Done():
		log.Fatal("timeout exceeded, something is probably wrong with the rollout")
	}
}

// NewASGSet returns an ASGSet pointer
func (r *BaseRunner) NewASGSet(ctx context.Context) (*ASGSet, error) {
	return newASGSet(ctx, r.awsClients, r.asgs, r.Opts.Force, r.startTime)
}
