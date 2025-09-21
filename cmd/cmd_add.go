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
	"errors"
	"fmt"
	"lucy/remote"
	"lucy/remote/sources"
	"lucy/tools"
	"lucy/util"

	"github.com/urfave/cli/v3"
	"lucy/local"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/syntax"
)

var subcmdAdd = &cli.Command{
	Name:  "add",
	Usage: "Add new mods, plugins, or server modules",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "Ignore version, dependency, and platform warnings",
			Value:   false,
		},
		&cli.StringFlag{
			Name:    "source",
			Aliases: []string{"s"},
			Usage:   "Specify the source to download from (modrinth, mcdr)",
			Value:   "none",
		},
		flagNoStyle,
	},
	Action: tools.Decorate(
		actionAdd,
		decoratorGlobalFlags,
		decoratorHelpAndExitOnNoArg,
		decoratorLogAndExitOnError,
	),
}

var actionAdd cli.ActionFunc = func(
ctx context.Context,
cmd *cli.Command,
) error {

	// get id from args
	id := syntax.Parse(cmd.Args().First())

	// probe server info
	serverInfo := local.GetServerInfo()

	// ensure we are in a lucy-managed server
	if !serverInfo.HasLucy {
		return errors.New("lucy is not installed, run `lucy init` before downloading mods")
	}

	if serverInfo.Executable == local.UnknownExecutable {
		return errors.New("no executable found, `lucy add` requires a server in current directory")
	}

	// check if the specified platform matches the server platform
	if id.Platform != lucytypes.AllPlatform && id.Platform != serverInfo.Executable.Platform {
		logger.Error(errors.New("platform mismatch"))
		return nil
	}

	// Get the appropriate directory to download file to.
	// This is a temporary solution. Installation is not supposed to be this simple.
	// The installer should be designed as an injectable interface to allow non-standard
	// installation methods.
	var dir string
	switch id.Platform {
	case lucytypes.AllPlatform:
		logger.InfoNow("no platform specified, attempting to infer")
	case lucytypes.Mcdr:
		if serverInfo.Mcdr == nil {
			return errors.New("mcdr not found, please install mcdr first")
		}
		dir = serverInfo.Mcdr.PluginPaths[0]
	case lucytypes.Forge, lucytypes.Fabric:
		dir = serverInfo.ModPath
	default:
		return errors.New("unsupported platform")
	}

	p := lucytypes.Package{
		Id:           id,
		Dependencies: nil,
		Local:        nil,
		Remote:       nil,
		Supports:     nil,
		Information:  nil,
	}

	var rawRemote remote.RawPackageRemote
	var err error

	switch cmd.String("source") {
	case "none":
		for _, source := range sources.All {
			rawRemote, err = source.Fetch(id)
			if err != nil {
				logger.WarnNow(err)
				err = nil // prevent error got printed twice in the last iteration
				continue
			}
			if rawRemote != nil {
				break
			}
		}
	case sources.Mcdr.Name().String():
		rawRemote, err = sources.Mcdr.Fetch(id)
	case sources.Modrinth.Name().String():
		rawRemote, err = sources.Modrinth.Fetch(id)
	default:
		return fmt.Errorf("unknown source: %s", cmd.String("source"))
	}
	if err != nil {
		logger.WarnNow(err)
	}
	if rawRemote != nil {
		r := rawRemote.ToPackageRemote()
		p.Remote = &r
	}

	// TODO: util.DownloadFile is a temporary solution
	_, _, err = util.DownloadFile(p.Remote.FileUrl, dir)
	if err != nil {
		logger.ErrorNow(fmt.Errorf("download failed: %w", err))
	}
	return nil
}
