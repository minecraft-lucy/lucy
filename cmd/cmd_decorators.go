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
	"lucy/tools"

	"github.com/urfave/cli/v3"
	"lucy/lnout"
)

// globalFlagsDecorator is a high-order function that appends global flag actions
// to the action function.
func globalFlagsDecorator(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Bool("log-file") {
			println("Log file at", lnout.LogFile.Name())
		}
		if cmd.Bool("verbose") {
			lnout.UseConsoleOutput()
		}
		if cmd.Bool("debug") {
			lnout.UseDebug()
		}
		if cmd.Bool("no-style") {
			tools.TurnOffStyles()
		}
		return f(ctx, cmd)
	}
}

// helpOnNoInputDecorator is a high-order function that takes a cli.ActionFunc and
// returns a cli.ActionFunc that prints help and exit when there's no args and flags.
//
// This function is not necessarily used for every action function, as some
// action functions are expected to have no args and flags. E.g., `lucy status`.
func helpOnNoInputDecorator(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() == 0 && len(cmd.LocalFlagNames()) == 0 {
			cli.ShowAppHelpAndExit(cmd, 0)
		}
		return f(ctx, cmd)
	}
}

func helpOnErrorDecorator(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		err := f(ctx, cmd)
		if err != nil {
			cli.ShowAppHelpAndExit(cmd, 1)
		}
		return err
	}
}
