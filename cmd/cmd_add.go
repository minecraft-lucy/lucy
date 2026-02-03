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
	"lucy/probe"
	"lucy/remote"
	"lucy/remote/source"
	"lucy/tools"
	"lucy/util"

	"lucy/logger"
	"lucy/syntax"
	"lucy/types"

	"github.com/urfave/cli/v3"
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
	),
}

var actionAdd cli.ActionFunc = func(
	ctx context.Context,
	cmd *cli.Command,
) error {

	// get id from args
	id := syntax.Parse(cmd.Args().First())

	// probe server info
	serverInfo := probe.GetServerInfo()

	// ensure we are in a lucy-managed server
	// TODO: Disabled for now, the part for building the program directory is not done
	// if !serverInfo.HasLucy {
	// 	return errors.New("lucy is not installed, run `lucy init` before downloading mods")
	// }

	if serverInfo.Executable == probe.UnknownExecutable {
		return errors.New("no executable found, `lucy add` requires a server in current directory")
	}

	// check if the specified platform matches the server platform
	if id.Platform != types.AllPlatform {
		if id.Platform == types.Mcdr {
			// for mcdr, we only need to check if it's mcdr-managed
			if serverInfo.Environments.Mcdr == nil {
				return errors.New("mcdr not found")
			}
		} else if id.Platform != serverInfo.Executable.LoaderPlatform {
			return errors.New("platform mismatch")
		}
	}

	// Get the appropriate directory to download file to.
	// This is a temporary solution. Installation is not supposed to be this simple.
	// The installer should be designed as an injectable interface to allow non-standard
	// installation methods.
	var dir string
	switch id.Platform {
	case types.AllPlatform:
		logger.InfoNow("no platform specified, attempting to infer")
	case types.Mcdr:
		dir = serverInfo.Environments.Mcdr.PluginPaths[0]
	case types.Forge, types.Fabric:
		dir = serverInfo.ModPath
	default:
		return errors.New("unsupported platform")
	}

	p := types.Package{
		Id:           id,
		Dependencies: nil,
		Local:        nil,
		Remote:       nil,
		Supports:     nil,
		Information:  nil,
	}

	// fetch remote data
	var remoteData remote.RawPackageRemote
	var src remote.SourceHandler
	var err error

	switch cmd.String("source") {
	case "none":
		for _, src = range source.All {
			remoteData, err = src.Fetch(id)
			if err != nil {
				logger.InfoNow(err)
				err = nil // prevent error got printed twice in the last iteration
				continue
			}
			if remoteData != nil {
				// found the package, exit loop
				break
			}
		}
	case source.Mcdr.Name().String():
		if id.Platform != types.Mcdr && id.Platform != types.AllPlatform {
			return fmt.Errorf("source 'mcdr' only supports mcdr platform")
		}
		remoteData, err = source.Mcdr.Fetch(id)
	case source.Modrinth.Name().String():
		if id.Platform == types.Mcdr {
			return fmt.Errorf("source 'modrinth' does not support mcdr platform")
		}
		remoteData, err = source.Modrinth.Fetch(id)
	default:
		return fmt.Errorf("unknown source: %s", cmd.String("source"))
	}
	if err != nil {
		logger.WarnNow(err)
	}
	if remoteData != nil {
		r := remoteData.ToPackageRemote()
		p.Remote = &r
	} else {
		return errors.New("package not found in any source")
	}

	// let's try to get the correct dependency info first
	// for sources like modrinth, the dependency info from remote is not reliable
	if id.Platform == types.Mcdr {
		depsData, err := src.Dependencies(id)
		if err != nil {
			logger.Debug(err)
		}
		tools.PrintAsJson(depsData.ToPackageDependencies())
		return nil
	}

	// TODO: util.DownloadFile is a temporary solution
	_, _, err = util.DownloadFile(p.Remote.FileUrl, dir)
	if err != nil {
		logger.ErrorNow(fmt.Errorf("download failed: %w", err))
	}
	return nil
}

func cmdAddInstall() {

}
