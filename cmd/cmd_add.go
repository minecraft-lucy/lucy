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

	"lucy/tools"

	"github.com/urfave/cli/v3"
	"lucy/lnout"
	"lucy/local"
	"lucy/lucyerrors"
	"lucy/lucytypes"
	"lucy/remote/modrinth"
	"lucy/syntax"
	"lucy/util"
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
		globalFlagsDecorator,
		helpOnNoInputDecorator,
	),
}

// actionAdd
// The strategy is:
//
//   - Most up to date
//   - Compatible with the server version
//   - Release version
//
// TODO: Version specification
var actionAdd cli.ActionFunc = func(
	ctx context.Context,
	cmd *cli.Command,
) error {
	// TODO: Platform specification
	// TODO: Platform compatibility check
	// TODO: Error handling

	p := syntax.Parse(cmd.Args().First())
	serverInfo := local.GetServerInfo()

	if !serverInfo.HasLucy {
		return errors.New("lucy is not installed, run `lucy init` before downloading mods")
	}

	if serverInfo.Executable == local.UnknownExecutable {
		// Case where the server is not detected
		return errors.New("no executable found, `lucy add` requires a server in current directory")
	} else if p.Platform == lucytypes.Mcdr && serverInfo.Mcdr == nil {
		// Case where MCDR is not installed but the user wants to download MCDR plugins
		// TODO: Deal with this
		lnout.Error(errors.New("no mcdr found, while mcdr plugins requested"))
		return nil
	} else if p.Platform != lucytypes.AllPlatform && p.Platform != serverInfo.Executable.Platform {
		// Case where the platform of the mod is different from the server
		// TODO: Deal with this
		lnout.Error(errors.New("platform mismatch"))
		return nil
	}

	newestVersion, _ := modrinth.LatestCompatibleVersion(p.Name)
	downloadFile, err := util.DownloadFile(
		// Not sure how to deal with multiple files
		// As the motivation for publishers to provide multiple files is unclear
		// TODO: Maybe add a prompt to let the user choose
		newestVersion.Files[0].Url,
		"mod",
		newestVersion.Files[0].Filename,
	)
	if err != nil {
		if errors.Is(err, lucyerrors.NoLucyError) {
			lnout.Warn(err)
		} else {
			lnout.Error(errors.New("failed at downloading: " + err.Error()))
			return nil
		}
	}

	_ = util.MoveFile(downloadFile, serverInfo.ModPath)

	return nil
}
