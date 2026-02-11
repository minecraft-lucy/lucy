package cmd

import (
	"context"
	"fmt"

	"lucy/probe"
	"lucy/tools"
	"lucy/tui"
	"lucy/types"

	"github.com/urfave/cli/v3"
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
	serverInfo := probe.ServerInfo()
	if cmd.Bool("json") {
		tools.PrintAsJson(serverInfo)
	} else {
		tui.Flush(serverInfoToStatus(&serverInfo, cmd.Bool("long")))
	}
	return nil
}

func serverInfoToStatus(
	data *types.ServerInfo,
	longOutput bool,
) (status *tui.Data) {
	if data.Executable == nil {
		return &tui.Data{
			Fields: []tui.Field{
				&tui.FieldAnnotation{
					Annotation: "(No server found)",
				},
			},
		}
	}

	status = &tui.Data{
		Fields: []tui.Field{},
	}

	status.Fields = append(
		status.Fields, &tui.FieldAnnotatedShortText{
			Title:      "Game",
			Text:       data.Executable.GameVersion.String(),
			Annotation: data.Executable.Path,
			NoTab:      true,
		},
	)

	if data.Executable.LoaderPlatform != types.Minecraft {
		status.Fields = append(
			status.Fields, &tui.FieldAnnotatedShortText{
				Title:      "Modding",
				Text:       data.Executable.LoaderPlatform.Title(),
				Annotation: data.Executable.LoaderVersion.String(),
				NoTab:      true,
			},
		)
	}

	if data.Activity != nil {
		status.Fields = append(
			status.Fields, &tui.FieldAnnotatedShortText{
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
			status.Fields, &tui.FieldShortText{
				Title: "Activity",
				Text:  tools.Dim("(Unknown)"),
			},
		)
	}

	// Modding related fields only shown when modding platform detected
	if data.Executable.LoaderPlatform != types.Minecraft {
		if len(data.Packages) > 0 {
			modNames := make([]string, 0, len(data.Packages))
			modPaths := make([]string, 0, len(modNames))
			for _, mod := range data.Packages {
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
				tools.Ternary[tui.Field](
					longOutput,
					&tui.FieldMultiAnnotatedShortText{
						Title:     "Packages",
						Texts:     modNames,
						Annots:    modPaths,
						ShowTotal: true,
					},
					&tui.FieldDynamicColumnLabels{
						Title:     "Packages",
						Labels:    modNames,
						MaxLines:  0,
						ShowTotal: true,
					},
				),
			)
		} else {
			status.Fields = append(
				status.Fields, &tui.FieldMultiAnnotatedShortText{
					Title:     "Packages",
					Texts:     []string{tools.Dim("(None)")},
					Annots:    nil,
					ShowTotal: false,
				},
			)
		}
	}

	return status
}
