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

package cmd

import (
	"context"

	"github.com/palantir/bouncer/batchcanary"
	"github.com/palantir/bouncer/bouncer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var batchCanaryCmd = &cobra.Command{
	Use:   "batch-canary",
	Short: "Run bouncer in batch canary",
	Long:  `Run bouncer in batch canary mode, where we add a new node to an ASG, then if it's successful, cycle the rest of the nodes in batches.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(logLevelFromViper())

		log.Debug("canary called")
		if log.GetLevel() == log.DebugLevel {
			cmd.DebugFlags()
			viper.Debug()
		}

		asgString := viper.GetString("batchcanary.asg")
		if asgString == "" {
			log.Fatal("You must specify ASG to cycle nodes from")
		}

		commandString := viper.GetString("batchcanary.command")
		noop := viper.GetBool("batchcanary.noop")
		force := viper.GetBool("batchcanary.force")
		batchSize := viper.GetInt32("batchcanary.batchsize")
		termHook := viper.GetString("terminate-hook")
		pendHook := viper.GetString("pending-hook")
		timeout := timeoutFromViper()

		if batchSize < 0 {
			log.Fatalf("Batch size must be >= 0, got %d", batchSize)
		}

		log.Debugf("Binding vars, got %+v %+v %+v %+v", asgString, noop, version, commandString)

		log.Info("Beginning bouncer batch canary run")

		opts := bouncer.RunnerOpts{
			Noop:          noop,
			BatchSize:     &batchSize,
			Force:         force,
			AsgString:     asgString,
			CommandString: commandString,
			TerminateHook: termHook,
			PendingHook:   pendHook,
			ItemTimeout:   timeout,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		log.RegisterExitHandler(cancel)

		r, err := batchcanary.NewRunner(ctx, &opts)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error initializing runner"))
		}

		err = r.ValidatePrereqs(ctx)
		if err != nil {
			log.Fatal(err)
		}

		err = r.Run()
		if err != nil {
			log.Fatal(errors.Wrap(err, "error in run"))
		}
	},
}

func init() {
	RootCmd.AddCommand(batchCanaryCmd)

	batchCanaryCmd.Flags().BoolP("noop", "n", false, "Run this in noop mode, and only print what you would do")
	err := viper.BindPFlag("batchcanary.noop", batchCanaryCmd.Flags().Lookup("noop"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'noop' to viper var 'batchcanary.noop' failed: %s"))
	}

	batchCanaryCmd.Flags().StringP("asg", "a", "", "ASG to refresh")
	err = viper.BindPFlag("batchcanary.asg", batchCanaryCmd.Flags().Lookup("asg"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'asg' to viper var 'batchcanary.asg' failed: %s"))
	}

	batchCanaryCmd.Flags().Int32P("batchsize", "b", 0, "Max number of nodes to refresh at a time after the single canary. Defaults to all remaining nodes.")
	err = viper.BindPFlag("batchcanary.batchsize", batchCanaryCmd.Flags().Lookup("batchsize"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'batchsize' to viper var 'batchcanary.batchsize' failed: %s"))
	}

	batchCanaryCmd.Flags().StringP("preterminatecall", "p", "", "External command to run before host is removed from its ELB & terminate process begins")
	err = viper.BindPFlag("batchcanary.command", batchCanaryCmd.Flags().Lookup("preterminatecall"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'command' to viper var 'batchcanary.command' failed: %s"))
	}

	batchCanaryCmd.Flags().BoolP("force", "f", false, "Force all nodes to be recycled, even if they're running the latest launch config")
	err = viper.BindPFlag("batchcanary.force", batchCanaryCmd.Flags().Lookup("force"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'force' to viper var 'batchcanary.force' failed: %s"))
	}
}
