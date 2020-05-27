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
	log "github.com/Sirupsen/logrus"
	"github.com/palantir/bouncer/bouncer"
	"github.com/palantir/bouncer/canary"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var canaryCmd = &cobra.Command{
	Use:   "canary",
	Short: "Run bouncer in canary",
	Long:  `Run bouncer in canary mode, where we add a new node to an ASG, then if it's successful, cycle the rest of the nodes.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(logLevelFromViper())

		log.Debug("canary called")
		if log.GetLevel() == log.DebugLevel {
			cmd.DebugFlags()
			viper.Debug()
		}

		asgString := viper.GetString("canary.asg")
		if asgString == "" {
			log.Fatal("You must specify ASG to cycle nodes from")
		}

		commandString := viper.GetString("canary.command")
		noop := viper.GetBool("canary.noop")
		force := viper.GetBool("canary.force")
		termHook := viper.GetString("terminate-hook")
		pendHook := viper.GetString("pending-hook")
		timeout := timeoutFromViper()

		log.Debugf("Binding vars, got %+v %+v %+v %+v", asgString, noop, version, commandString)

		log.Info("Beginning bouncer canary run")

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

		r, err := canary.NewRunner(&opts)
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
	RootCmd.AddCommand(canaryCmd)

	canaryCmd.Flags().BoolP("noop", "n", false, "Run this in noop mode, and only print what you would do")
	err := viper.BindPFlag("canary.noop", canaryCmd.Flags().Lookup("noop"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'noop' to viper var 'canary.noop' failed: %s"))
	}

	canaryCmd.Flags().StringP("asg", "a", "", "ASG to refresh")
	err = viper.BindPFlag("canary.asg", canaryCmd.Flags().Lookup("asg"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'asg' to viper var 'canary.asg' failed: %s"))
	}

	canaryCmd.Flags().StringP("preterminatecall", "p", "", "External command to run before host is removed from its ELB & terminate process begins")
	err = viper.BindPFlag("canary.command", canaryCmd.Flags().Lookup("preterminatecall"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'command' to viper var 'canary.command' failed: %s"))
	}

	canaryCmd.Flags().BoolP("force", "f", false, "Force all nodes to be recycled, even if they're running the latest launch config")
	err = viper.BindPFlag("canary.force", canaryCmd.Flags().Lookup("force"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'force' to viper var 'canary.force' failed: %s"))
	}
}
