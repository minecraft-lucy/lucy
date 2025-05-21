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
		flagNoStyle,
	},
	Action: tools.Decorate(
		actionAdd,
		decoratorGlobalFlags,
		decoratorHelpAndExitOnNoInput,
	),
}

var actionAdd cli.ActionFunc = func(
	ctx context.Context,
	cmd *cli.Command,
) error {
	id := syntax.Parse(cmd.Args().First())
	serverInfo := local.GetServerInfo()
	if !serverInfo.HasLucy {
		return errors.New("lucy is not installed, run `lucy init` before downloading mods")
	}

	if serverInfo.Executable == local.UnknownExecutable {
		return errors.New("no executable found, `lucy add` requires a server in current directory")
	}
	if id.Platform != lucytypes.AllPlatform && id.Platform != serverInfo.Executable.Platform {
		logger.Error(errors.New("platform mismatch"))
		return nil
	}

	// var handler remote.SourceHandler
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

	raw, err := sources.Mcdr.Fetch(id)
	if err != nil {
		return err
	}
	remote := raw.ToPackageRemote()
	_, _, err = util.DownloadFile(remote.FileUrl, dir)
	if err != nil {
		logger.ErrorNow(fmt.Errorf("download failed: %w", err))
	}
	return nil
}
