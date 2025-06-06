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

	"github.com/urfave/cli/v3"
	"lucy/local"
	"lucy/lucytypes"
	"lucy/structout"
	"lucy/tools"
)

var subcmdStatus = &cli.Command{
	Name:  "status",
	Usage: "Display basic information of the current server",
	Action: tools.Decorate(
		actionStatus,
		decoratorGlobalFlags,
	),
	Flags: []cli.Flag{
		flagJsonOutput,
		flagLongOutput,
	},
}

var actionStatus cli.ActionFunc = func(
	_ context.Context,
	cmd *cli.Command,
) error {
	serverInfo := local.GetServerInfo()
	if cmd.Bool("json") {
		tools.PrintAsJson(serverInfo)
	} else {
		structout.Flush(serverInfoToStatus(&serverInfo, cmd.Bool("long")))
	}
	return nil
}

func serverInfoToStatus(
	data *lucytypes.ServerInfo,
	longOutput bool,
) (status *structout.Data) {
	if data.Executable == nil {
		return &structout.Data{
			Fields: []structout.Field{
				&structout.FieldAnnotation{
					Annotation: "(No server found)",
				},
			},
		}
	}

	status = &structout.Data{
		Fields: []structout.Field{},
	}

	status.Fields = append(
		status.Fields, &structout.FieldAnnotatedShortText{
			Title:      "Game",
			Text:       data.Executable.GameVersion.String(),
			Annotation: data.Executable.Path,
			NoTab:      true,
		},
	)

	if data.Executable.Platform != lucytypes.Minecraft {
		status.Fields = append(
			status.Fields, &structout.FieldAnnotatedShortText{
				Title:      "Modding",
				Text:       data.Executable.Platform.Title(),
				Annotation: data.Executable.LoaderVersion.String(),
				NoTab:      true,
			},
		)
	}

	if data.Activity != nil {
		status.Fields = append(
			status.Fields, &structout.FieldAnnotatedShortText{
				Title: "Activity",
				Text: tools.Ternary(
					data.Activity.Active,
					"Active",
					"Inactive",
				),
				Annotation: tools.Ternary(
					data.Activity.Active,
					fmt.Sprintf("PID %d", data.Activity.Pid),
					"",
				),
				NoTab: true,
			},
		)
	} else {
		status.Fields = append(
			status.Fields, &structout.FieldShortText{
				Title: "Activity",
				Text:  tools.Dim("(Unknown)"),
			},
		)
	}

	// Modding related fields only shown when modding platform detected
	if data.Executable.Platform != lucytypes.Minecraft {
		if len(data.Mods) > 0 {
			modNames := make([]string, 0, len(data.Mods))
			modPaths := make([]string, 0, len(modNames))
			for _, mod := range data.Mods {
				modNames = append(
					modNames,
					tools.Ternary(
						longOutput,
						mod.Id.StringFull(),
						mod.Id.StringNameVersion(),
					),
				)
				modPaths = append(modPaths, mod.Local.Path)
			}
			status.Fields = append(
				status.Fields,
				tools.Ternary[structout.Field](
					longOutput,
					&structout.FieldMultiAnnotatedShortText{
						Title:     "Mods",
						Texts:     modNames,
						Annots:    modPaths,
						ShowTotal: true,
					},
					&structout.FieldDynamicColumnLabels{
						Title:     "Mods",
						Labels:    modNames,
						MaxLines:  0,
						ShowTotal: true,
					},
				),
			)
		} else {
			status.Fields = append(
				status.Fields, &structout.FieldMultiAnnotatedShortText{
					Title:     "Mods",
					Texts:     []string{tools.Dim("(None)")},
					Annots:    nil,
					ShowTotal: false,
				},
			)
		}
	}

	if data.Mcdr != nil {
		pluginNames := make([]string, 0, len(data.Mcdr.PluginList))
		pluginPaths := make([]string, 0, len(data.Mcdr.PluginList))
		for _, plugin := range data.Mcdr.PluginList {
			pluginNames = append(
				pluginNames,
				tools.Ternary(
					longOutput,
					plugin.Id.StringFull(),
					plugin.Id.StringNameVersion(),
				),
			)
			pluginPaths = append(
				pluginPaths,
				tools.Ternary(longOutput, plugin.Local.Path, ""),
			)
		}
		status.Fields = append(
			status.Fields, &structout.FieldMultiAnnotatedShortText{
				Title:     "MCDR Plugins",
				Texts:     pluginNames,
				Annots:    pluginPaths,
				ShowTotal: true,
			},
		)

	}

	return status
}
