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
	"github.com/palantir/bouncer/serial"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serialCmd = &cobra.Command{
	Use:   "serial",
	Short: "Run bouncer in serial",
	Long:  `Run bouncer in serial mode, where we bounce one node at a time from the list of ASGs.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(logLevelFromViper())

		log.Debug("serial called")
		if log.GetLevel() == log.DebugLevel {
			cmd.DebugFlags()
			viper.Debug()
		}

		asgsString := viper.GetString("serial.asgs")
		if asgsString == "" {
			log.Fatal("You must specify ASGs to cycle nodes from (in a comma-delimited list)")
		}

		commandString := viper.GetString("serial.command")
		noop := viper.GetBool("serial.noop")
		force := viper.GetBool("serial.force")
		termHook := viper.GetString("terminate-hook")
		pendHook := viper.GetString("pending-hook")
		timeout := timeoutFromViper()

		log.Debugf("Binding vars, got %+v %+v %+v %+v", asgsString, noop, version, commandString)

		log.Info("Beginning bouncer serial run")

		var defCap int64
		defCap = 1
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

		r, err := serial.NewRunner(&opts)
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
	RootCmd.AddCommand(serialCmd)

	serialCmd.Flags().BoolP("noop", "n", false, "Run this in noop mode, and only print what you would do")
	err := viper.BindPFlag("serial.noop", serialCmd.Flags().Lookup("noop"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'noop' to viper var 'serial.noop' failed: %s"))
	}

	serialCmd.Flags().StringP("asgs", "a", "", "ASGs to check for nodes to cycle in")
	err = viper.BindPFlag("serial.asgs", serialCmd.Flags().Lookup("asgs"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'asgs' to viper var 'serial.asgs' failed: %s"))
	}

	serialCmd.Flags().StringP("preterminatecall", "p", "", "External command to run before host is removed from its ELB & terminate process begins")
	err = viper.BindPFlag("serial.command", serialCmd.Flags().Lookup("preterminatecall"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'command' to viper var 'serial.command' failed: %s"))
	}

	serialCmd.Flags().BoolP("force", "f", false, "Force all nodes to be recycled, even if they're running the latest launch config")
	err = viper.BindPFlag("serial.force", serialCmd.Flags().Lookup("force"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'force' to viper var 'serial.force' failed: %s"))
	}
}
