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
	"github.com/palantir/bouncer/bouncer"
	"github.com/palantir/bouncer/slowcanary"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var slowCanaryCmd = &cobra.Command{
	Use:   "slow-canary",
	Short: "Run bouncer in slow-canary",
	Long:  `Run bouncer in slow-canary mode, where we add a new node to an ASG, then remove an old, and repeat until we've cycled all the nodes.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(logLevelFromViper())

		log.Debug("slow-canary called")
		if log.GetLevel() == log.DebugLevel {
			cmd.DebugFlags()
			viper.Debug()
		}

		asgString := viper.GetString("slow-canary.asg")
		if asgString == "" {
			log.Fatal("You must specify ASG to cycle nodes from")
		}

		commandString := viper.GetString("slow-canary.command")
		noop := viper.GetBool("slow-canary.noop")
		force := viper.GetBool("slow-canary.force")
		termHook := viper.GetString("terminate-hook")
		pendHook := viper.GetString("pending-hook")
		timeout := timeoutFromViper()

		log.Debugf("Binding vars, got %+v %+v %+v %+v", asgString, noop, version, commandString)

		log.Info("Beginning bouncer slow-canary run")

		opts := bouncer.RunnerOpts{
			Noop:            noop,
			Force:           force,
			AsgString:       asgString,
			CommandString:   commandString,
			DefaultCapacity: nil,
			TerminateHook:   termHook,
			PendingHook:     pendHook,
			ItemTimeout:     timeout,
		}

		r, err := slowcanary.NewRunner(&opts)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error initializing runner"))
		}

		r.MustValidatePrereqs()

		err = r.Run()
		if err != nil {
			log.Fatal(errors.Wrap(err, "error in run"))
		}
	},
}

func init() {
	RootCmd.AddCommand(slowCanaryCmd)

	slowCanaryCmd.Flags().BoolP("noop", "n", false, "Run this in noop mode, and only print what you would do")
	err := viper.BindPFlag("slow-canary.noop", slowCanaryCmd.Flags().Lookup("noop"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'noop' to viper var 'slow-canary.noop' failed: %s"))
	}

	slowCanaryCmd.Flags().StringP("asg", "a", "", "ASG to refresh")
	err = viper.BindPFlag("slow-canary.asg", slowCanaryCmd.Flags().Lookup("asg"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'asg' to viper var 'slow-canary.asg' failed: %s"))
	}

	slowCanaryCmd.Flags().StringP("preterminatecall", "p", "", "External command to run before host is removed from its ELB & terminate process begins")
	err = viper.BindPFlag("slow-canary.command", slowCanaryCmd.Flags().Lookup("preterminatecall"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'command' to viper var 'slow-canary.command' failed: %s"))
	}

	slowCanaryCmd.Flags().BoolP("force", "f", false, "Force all nodes to be recycled, even if they're running the latest launch config")
	err = viper.BindPFlag("slow-canary.force", slowCanaryCmd.Flags().Lookup("force"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'force' to viper var 'slow-canary.force' failed: %s"))
	}
}
