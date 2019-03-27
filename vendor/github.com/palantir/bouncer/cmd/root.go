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
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	version     = "unspecified"
	versionFlag bool
	validate    *validator.Validate
)

const killswitchVar = "BOUNCER_KILLSWITCH"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "bouncer",
	Short: "An app that bounces AWS instances in the given ASGs.",
	Long:  `Bounces AWS instances that are due to be cycled in the ASGs passed-in.`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Println(version)
			os.Exit(0)
		}
		switch cmdName := cmd.Name(); cmdName {
		case "serial":
			log.SetFormatter(&log.TextFormatter{})
			log.SetOutput(os.Stdout)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("You must provide a command")
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	validate = validator.New()

	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable Verbose debugging output")
	err := viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error binding verbose flag"))
	}

	RootCmd.PersistentFlags().Int32P("timeout", "t", 20, "Timeout for each AWS mutation action (in minutes)")
	err = viper.BindPFlag("timeout", RootCmd.PersistentFlags().Lookup("timeout"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error binding timeout flag"))
	}

	RootCmd.PersistentFlags().BoolVar(&versionFlag, "version", false, "Print Version and exit")

	RootCmd.PersistentFlags().String("terminate-hook", "terminate-hook", "Name of the hook on the autoscaling:EC2_INSTANCE_TERMINATING transition")
	err = viper.BindPFlag("terminate-hook", RootCmd.PersistentFlags().Lookup("terminate-hook"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error binding terminate-hook flag"))
	}

	RootCmd.PersistentFlags().String("pending-hook", "pending-hook", "Name of the hook on the autoscaling:EC2_INSTANCE_LAUNCHING transition")
	err = viper.BindPFlag("pending-hook", RootCmd.PersistentFlags().Lookup("pending-hook"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error binding pending-hook flag"))
	}

	// Check for special killswitch
	val := os.Getenv(killswitchVar)
	if val != "" {
		log.Warn("Killswitch variable found, skipping all actions and exiting with success.")
		os.Exit(0)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

func timeoutFromViper() time.Duration {
	return time.Duration(viper.GetInt("timeout")) * time.Minute
}

func logLevelFromViper() log.Level {
	if viper.GetBool("verbose") {
		return log.DebugLevel
	}
	logLevel := viper.GetString("bouncer.log")
	level := strings.ToLower(logLevel)
	switch level {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}
