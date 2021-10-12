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

	"github.com/palantir/bouncer/bouncer"
	"github.com/palantir/bouncer/rolling"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rollingCmd = &cobra.Command{
	Use:   "rolling",
	Short: "Run bouncer in rolling",
	Long:  `Run bouncer in rolling mode, where we bounce one node at a time from the list of ASGs.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(logLevelFromViper())

		log.Debug("rolling called")
		if log.GetLevel() == log.DebugLevel {
			cmd.DebugFlags()
			viper.Debug()
		}

		asgsString := viper.GetString("rolling.asgs")
		if asgsString == "" {
			log.Fatal("You must specify ASGs to cycle nodes from (in a comma-delimited list)")
		}

		commandString := viper.GetString("rolling.command")
		noop := viper.GetBool("rolling.noop")
		force := viper.GetBool("rolling.force")
		termHook := viper.GetString("terminate-hook")
		pendHook := viper.GetString("pending-hook")
		timeout := timeoutFromViper()

		log.Debugf("Binding vars, got %+v %+v %+v %+v", asgsString, noop, version, commandString)

		log.Info("Beginning bouncer rolling run")

		var defCap int32 = 1
		opts := bouncer.RunnerOpts{
			Noop:            noop,
			Force:           force,
			AsgString:       asgsString,
			CommandString:   commandString,
			DefaultCapacity: &defCap,
			TerminateHook:   termHook,
			PendingHook:     pendHook,
			ItemTimeout:     timeout,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		r, err := rolling.NewRunner(ctx, &opts)
		if err != nil {
			cancel()
			log.Fatal(errors.Wrap(err, "error initializing runner"))
		}

		err = r.ValidatePrereqs(ctx)
		if err != nil {
			cancel()
			log.Fatal(err)
		}

		err = r.Run(ctx)
		if err != nil {
			cancel()
			log.Fatal(errors.Wrap(err, "error in run"))
		}
	},
}

func init() {
	RootCmd.AddCommand(rollingCmd)

	rollingCmd.Flags().BoolP("noop", "n", false, "Run this in noop mode, and only print what you would do")
	err := viper.BindPFlag("rolling.noop", rollingCmd.Flags().Lookup("noop"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'noop' to viper var 'rolling.noop' failed: %s"))
	}

	rollingCmd.Flags().StringP("asgs", "a", "", "ASGs to check for nodes to cycle in")
	err = viper.BindPFlag("rolling.asgs", rollingCmd.Flags().Lookup("asgs"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'asgs' to viper var 'rolling.asgs' failed: %s"))
	}

	rollingCmd.Flags().StringP("preterminatecall", "p", "", "External command to run before host is removed from its ELB & terminate process begins")
	err = viper.BindPFlag("rolling.command", rollingCmd.Flags().Lookup("preterminatecall"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'command' to viper var 'rolling.command' failed: %s"))
	}

	rollingCmd.Flags().BoolP("force", "f", false, "Force all nodes to be recycled, even if they're running the latest launch config")
	err = viper.BindPFlag("rolling.force", rollingCmd.Flags().Lookup("force"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'force' to viper var 'rolling.force' failed: %s"))
	}
}
