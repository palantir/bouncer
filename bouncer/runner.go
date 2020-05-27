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
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/palantir/bouncer/aws"
	"github.com/pkg/errors"
)

// RunnerOpts is user-supplied options to any flavor of runner
type RunnerOpts struct {
	Noop            bool
	Force           bool
	AsgString       string
	CommandString   string
	DefaultCapacity *int64
	TerminateHook   string
	PendingHook     string
	ItemTimeout     time.Duration
}

// BaseRunner is the base struct for any runner
type BaseRunner struct {
	opts       *RunnerOpts
	startTime  time.Time
	awsClients *aws.Clients
	asgs       []*DesiredASG
	// curTimerStart is mutable - it's tracking the most recent API call that we need to time
	curTimerStart time.Time
}

const (
	waitBetweenChecks = 15 * time.Second
	// Sleep time and number of times to retry non-destructive AWS API calls
	apiRetryCount = 10
	apiRetrySleep = 10 * time.Second

	asgSeparator        = ","
	desiredCapSeparator = ":"
)

func retry(attempts int, sleep time.Duration, callback func() error) (err error) {
	for i := 0; ; i++ {
		err = callback()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep)

		log.Warn(errors.Wrap(err, "found error, retrying"))
	}
	return errors.Wrapf(err, "error persists after %v tries", attempts)
}

// NewBaseRunner instantiates a BaseRunner
func NewBaseRunner(opts *RunnerOpts) (*BaseRunner, error) {
	awsClients, err := aws.GetAWSClients()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting AWS Creds")
	}

	asgs, err := getASGList(opts)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing ASG list")
	}

	r := BaseRunner{
		opts:          opts,
		startTime:     time.Now(),
		awsClients:    awsClients,
		asgs:          asgs,
		curTimerStart: time.Now(),
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

		curAsg, err := extractDesiredASG(asgItem, desiredCapSeparator, opts.DefaultCapacity, command)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing ASG item")
		}
		asgs = append(asgs, curAsg)
	}

	return asgs, nil
}

func (r *BaseRunner) noopCheck() {
	if r.opts.Noop {
		log.Warn("NOOP only - not actually performing previous action, and exiting script with success")
		os.Exit(0)
	}
}

func (r *BaseRunner) abandonLifecycle(inst *Instance, hook *string) error {
	log.WithFields(log.Fields{
		"InstanceID":     *inst.ASGInstance.InstanceId,
		"Hook":           *hook,
		"LifecycleState": *inst.ASGInstance.LifecycleState,
	}).Warn("Issuing ABANDON to hook instead of terminating")
	result := "ABANDON"
	r.resetTimeout()
	r.noopCheck()
	err := r.awsClients.CompleteLifecycleAction(inst.AutoscalingGroup.AutoScalingGroupName, inst.ASGInstance.InstanceId, hook, &result)
	return errors.Wrap(err, "error completing lifecycle action")
}

// KillInstance calls TerminateInstanceInAutoscalingGroup, or, if the instance is stuck
// in a lifecycle hook, issues an ABANDON to it, killing it more forcefully
func (r *BaseRunner) KillInstance(inst *Instance) error {
	log.WithFields(log.Fields{
		"ASG":        *inst.AutoscalingGroup.AutoScalingGroupName,
		"InstanceID": *inst.ASGInstance.InstanceId,
	}).Info("Picked instance to die next")
	var hook string

	if *inst.ASGInstance.LifecycleState == autoscaling.LifecycleStatePendingWait {
		hook = r.opts.PendingHook
	}

	if *inst.ASGInstance.LifecycleState == autoscaling.LifecycleStateTerminatingWait {
		hook = r.opts.TerminateHook
	}

	if hook != "" {
		err := r.abandonLifecycle(inst, &hook)
		return errors.Wrapf(err, "error abandoning hook %s", hook)
	}

	if inst.PreTerminateCmd != nil {
		err := r.executeExternalCommand(*inst.PreTerminateCmd)
		if err != nil {
			return errors.Wrap(err, "error executing pre-terminate command")
		}
	}
	err := r.terminateInstanceInASG(inst)
	return errors.Wrap(err, "error terminating instance")
}

func (r *BaseRunner) terminateInstanceInASG(inst *Instance) error {
	log.WithFields(log.Fields{
		"ASG":        *inst.AutoscalingGroup.AutoScalingGroupName,
		"InstanceID": *inst.ASGInstance.InstanceId,
	}).Info("Terminating instance")
	r.resetTimeout()
	r.noopCheck()
	return r.awsClients.TerminateInstanceInASG(inst.ASGInstance.InstanceId)
}

// SetDesiredCapacity Updates desired capacity of ASG
// This function should only be used to increase desired cap, not decrease, since AWS
// will _always_ remove instances based on AZ before any other criteria
// http://docs.aws.amazon.com/autoscaling/latest/userguide/as-instance-termination.html
func (r *BaseRunner) SetDesiredCapacity(asg *ASG, desiredCapacity *int64) error {

	log.WithFields(log.Fields{
		"ASG":           *asg.ASG.AutoScalingGroupName,
		"CurDesiredCap": *asg.ASG.DesiredCapacity,
		"NewDesiredCap": *desiredCapacity,
	}).Info("Changing desired capacity")
	r.noopCheck()

	r.resetTimeout()
	err := r.awsClients.SetDesiredCapacity(asg.ASG, desiredCapacity)
	return errors.Wrap(err, "error setting desired capacity of ASG")
}

// TimedOut returns whether we've hit our runner's timeout or not
// This is based on curTimerStart, so call SetcurTimerStart whenever a new call is made
// that should be timed
func (r *BaseRunner) TimedOut() bool {
	curTime := time.Now()

	timeout := r.curTimerStart.Add(r.opts.ItemTimeout)

	log.WithFields(log.Fields{
		"Current Time": getHumanShortTime(curTime),
		"Timeout Time": getHumanShortTime(timeout),
	}).Debug("Checking if we're at timeout")

	return curTime.After(timeout)
}

func getHumanShortTime(t time.Time) string {
	zonename, _ := t.In(time.Local).Zone()
	human := fmt.Sprintf("%02v:%02v:%02v %s", t.Hour(), t.Minute(), t.Second(), zonename)
	return human
}

// resetTimeout sets curTimerStart to the current time, so call this
// after making a new AWS API call that should restart the timer with respect to the timeout
func (r *BaseRunner) resetTimeout() {
	now := time.Now()
	log.WithFields(log.Fields{
		"Last timer reset":              getHumanShortTime(r.curTimerStart),
		"Time elapsed since last reset": now.Sub(r.curTimerStart),
	}).Debug("Resetting timer")
	r.curTimerStart = now
}

// Sleep makes us sleep for the constant time - call this when waiting for an AWS change
func (r *BaseRunner) Sleep() {
	log.Debugf("Sleeping for %v", waitBetweenChecks)
	time.Sleep(waitBetweenChecks)
}

// NewASGSet returns an ASGSet pointer
func (r *BaseRunner) NewASGSet() (*ASGSet, error) {
	return newASGSet(r.awsClients, r.asgs, r.opts.Force, r.startTime)
}
