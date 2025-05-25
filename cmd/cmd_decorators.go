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

	"lucy/logger"

	"github.com/urfave/cli/v3"
)

// decoratorBaseCommandFlags provides some the base command `lucy` some necessary
// flag actions.
func decoratorBaseCommandFlags(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Bool("no-style") {
			tools.TurnOffStyles()
		}
		if cmd.Bool("log-file") {
			println("Log file at", logger.LogFile.Name())
		}
		return f(ctx, cmd)
	}
}

// decoratorGlobalFlags is a high-order function that appends global flag actions
// to the action function.
func decoratorGlobalFlags(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {

		if cmd.Bool("verbose") {
			logger.UseConsoleOutput()
		}
		if cmd.Bool("debug") {
			logger.UseDebug()
		}
		if cmd.Bool("no-style") {
			tools.TurnOffStyles()
		}
		return f(ctx, cmd)
	}
}

// decoratorHelpAndExitOnNoArg is a high-order function that takes a cli.ActionFunc and
// returns a cli.ActionFunc that prints help and exit when there's no args specified.
//
// This function is not necessarily applicable to every action function, as some
// sub-commands are expected to have no args, e.g., `lucy status`.
func decoratorHelpAndExitOnNoArg(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() == 0 {
			cli.ShowSubcommandHelpAndExit(cmd, 0)
		}
		return f(ctx, cmd)
	}
}

// decoratorHelpAndExitOnNoFlag is similar to decoratorHelpAndExitOnNoArg, but
// it checks for flags instead of args. This is useful for commands that
// require at least one flag.
func decoratorHelpAndExitOnNoFlag(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if len(cmd.LocalFlagNames()) == 0 {
			cli.ShowSubcommandHelpAndExit(cmd, 0)
		}
		return f(ctx, cmd)
	}
}

func decoratorLogAndExitOnError(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		err := f(ctx, cmd)
		if err != nil {
			logger.ErrorNow(err)
			return err
		}
		return nil
	}
}

// decoratorHelpAndExitOnError exits with an error code and prints the help
//
// This means, with this decorator, you MUST NOT throw unexpected errors
// in your action function, as it will be caught and printed to the
// user.
//
// ONLY errors readable by the user should be thrown.
//
// Comparingly, decoratorLogAndExitOnError is more suitable for
// most of the action functions.
func decoratorHelpAndExitOnError(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		err := f(ctx, cmd)
		if err != nil {
			cli.ShowSubcommandHelpAndExit(cmd, 1)
		}
		return err
	}
}
