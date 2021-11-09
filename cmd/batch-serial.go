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

	"github.com/palantir/bouncer/batchserial"
	"github.com/palantir/bouncer/bouncer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var batchSerialCmd = &cobra.Command{
	Use:   "batch-serial",
	Short: "Run bouncer in batch serial",
	Long:  `Run bouncer in batch serial mode, where we destroy & recreate <batch size> nodes at a time from the list of ASGs.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(logLevelFromViper())

		log.Debug("batch-serial called")
		if log.GetLevel() == log.DebugLevel {
			cmd.DebugFlags()
			viper.Debug()
		}

		asgsString := viper.GetString("batchserial.asgs")
		if asgsString == "" {
			log.Fatal("You must specify ASGs to cycle nodes from (in a comma-delimited list)")
		}

		commandString := viper.GetString("batchserial.command")
		noop := viper.GetBool("batchserial.noop")
		force := viper.GetBool("batchserial.force")
		batchSize := viper.GetInt32("batchserial.batchsize")
		termHook := viper.GetString("terminate-hook")
		pendHook := viper.GetString("pending-hook")
		timeout := timeoutFromViper()

		if batchSize < 1 {
			log.Fatalf("Batch size must be >= 1, got %d", batchSize)
		}

		log.Debugf("Binding vars, got %+v %+v %+v %+v", asgsString, noop, version, commandString)

		log.Info("Beginning bouncer serial run")

		var defCap int32 = 1
		opts := bouncer.RunnerOpts{
			Noop:            noop,
			BatchSize:       &batchSize,
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
		log.RegisterExitHandler(cancel)

		r, err := batchserial.NewRunner(ctx, &opts)
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
	RootCmd.AddCommand(batchSerialCmd)

	batchSerialCmd.Flags().BoolP("noop", "n", false, "Run this in noop mode, and only print what you would do")
	err := viper.BindPFlag("batchserial.noop", batchSerialCmd.Flags().Lookup("noop"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'noop' to viper var 'batchserial.noop' failed: %s"))
	}

	batchSerialCmd.Flags().StringP("asgs", "a", "", "ASGs to check for nodes to cycle in")
	err = viper.BindPFlag("batchserial.asgs", batchSerialCmd.Flags().Lookup("asgs"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'asgs' to viper var 'batchserial.asgs' failed: %s"))
	}

	batchSerialCmd.Flags().Int32P("batchsize", "b", 1, "Max number of nodes to terminate at a time after the single canary. Defaults to all remaining nodes.")
	err = viper.BindPFlag("batchserial.batchsize", batchSerialCmd.Flags().Lookup("batchsize"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'batchsize' to viper var 'batchserial.batchsize' failed: %s"))
	}

	batchSerialCmd.Flags().StringP("preterminatecall", "p", "", "External command to run before host is removed from its ELB & terminate process begins")
	err = viper.BindPFlag("batchserial.command", batchSerialCmd.Flags().Lookup("preterminatecall"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'command' to viper var 'batchserial.command' failed: %s"))
	}

	batchSerialCmd.Flags().BoolP("force", "f", false, "Force all nodes to be recycled, even if they're running the latest launch config")
	err = viper.BindPFlag("batchserial.force", batchSerialCmd.Flags().Lookup("force"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'force' to viper var 'batchserial.force' failed: %s"))
	}
}
