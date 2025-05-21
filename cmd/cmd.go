/*
Copyright 2024 4rcadia

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"fmt"

	"lucy/tools"

	"github.com/urfave/cli/v3"
)

// Frontend should change when user do not run the program in CLI
// This is prepared for possible GUI implementation
var Frontend = "cli"

// Each subcommand (and its action function) should be in its own file

// Cli is the main command for lucy
var Cli = &cli.Command{
	Name:  "lucy",
	Usage: "The Minecraft server-side package manager",
	Action: tools.Decorate(
		actionEmpty,
		decoratorGlobalFlags,
		decoratorHelpAndExitOnNoInput,
	),
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "log-file",
			Aliases: []string{"l"},
			Usage:   "Output the path to logfile",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Print logs",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "Print debug logs",
			Value:   false,
		},
		flagNoStyle,
	},
	Commands: []*cli.Command{
		subcmdStatus,
		subcmdInfo,
		subcmdSearch,
		subcmdAdd,
		subcmdInit,
	},
	EnableShellCompletion:  true,
	Suggest:                true,
	UseShortOptionHandling: true,
	DefaultCommand:         "help",
	OnUsageError:           helpOnUsageError,
}

var helpOnUsageError cli.OnUsageErrorFunc = func(
	ctx context.Context,
	cmd *cli.Command,
	err error,
	isSubcommand bool,
) error {
	if isSubcommand {
		fmt.Println(fmt.Errorf("invalid command: %s", err).Error())
		cli.ShowAppHelpAndExit(cmd, 1)
	}
	fmt.Println(err.Error())
	return err
}

var actionEmpty cli.ActionFunc = func(
	ctx context.Context,
	cmd *cli.Command,
) error {
	return nil
}
