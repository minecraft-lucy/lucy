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
	"strconv"

	"lucy/remote"
	"lucy/remote/source"
	"lucy/types"

	"github.com/urfave/cli/v3"
	"lucy/logger"
	"lucy/syntax"
	"lucy/tools"
	"lucy/tui"
)

var subcmdSearch = &cli.Command{
	Name:  "search",
	Usage: "Search for mods and plugins",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "index",
			Aliases: []string{"i"},
			Usage:   "Index search results by `INDEX`",
			Value:   "relevance",
			Validator: func(s string) error {
				if types.SearchIndex(s).Valid() {
					return nil
				}
				return errors.New("must be one of \"relevance\", \"downloads\",\"newest\"")
			},
		},
		&cli.BoolFlag{
			Name:    "client",
			Aliases: []string{"c"},
			Usage:   "Also show client-only mods in results",
			Value:   false,
		},
		flagJsonOutput,
		flagLongOutput,
		flagNoStyle,
		flagSource,
	},
	ArgsUsage: "<platform/name>",
	Action: tools.Decorate(
		actionSearch,
		decoratorGlobalFlags,
		decoratorHelpAndExitOnNoArg,
		decoratorLogAndExitOnError,
	),
}

var (
	errorUnknownSource     = errors.New("unknown specified source")
	errorUnsupportedSource = errors.New("unsupported source")
	errorInvalidPlatform   = errors.New("invalid platform")
)

var actionSearch cli.ActionFunc = func(
	_ context.Context,
	cmd *cli.Command,
) error {
	p := syntax.Parse(cmd.Args().First())

	showClientPackage := cmd.Bool("client")
	indexBy := types.SearchIndex(cmd.String("index"))
	options := types.SearchOptions{
		ShowClientPackage: showClientPackage,
		IndexBy:           indexBy,
	}
	sourceStr := cmd.String("source")
	source := types.StringToSource(sourceStr)

	out := &tui.Data{}
	res := types.SearchResults{}
	var err error

	if source == types.AutoSource {
		switch p.Platform {
		case types.AllPlatform:
			for _, sourceHandler := range source.All {
				res, err = remote.Search(sourceHandler, p.Name, options)
				if err != nil {
					logger.WarnNow(
						fmt.Errorf(
							"search on %s failed: %w",
							sourceHandler.Name().Title(),
							err,
						),
					)
					continue
				}
				appendToSearchOutput(out, cmd.Bool("long"), res)
			}
		case types.Forge, types.Fabric, types.Neoforge:
			res, err = remote.Search(source.Modrinth, p.Name, options)
			if err != nil && !errors.Is(err, remote.ErrorNoResults) {
				logger.Fatal(err)
			}
			appendToSearchOutput(out, cmd.Bool("long"), res)
		case types.Mcdr:
			res, err = remote.Search(source.Mcdr, p.Name, options)
			if err != nil && !errors.Is(err, remote.ErrorNoResults) {
				logger.Fatal(err)
			}
			appendToSearchOutput(out, cmd.Bool("long"), res)
		case types.UnknownPlatform:
			logger.Fatal(
				fmt.Errorf(
					"%w: %s",
					errorInvalidPlatform,
					p.Platform,
				),
			)
		}
	} else {
		if source == types.UnknownSource {
			logger.Fatal(
				fmt.Errorf(
					"%w: %s",
					errorUnknownSource,
					sourceStr,
				),
			)
		}
		sourceHandler, ok := source.Map[source]
		if !ok {
			logger.Fatal(
				fmt.Errorf(
					"%w: %s",
					errorUnsupportedSource,
					source.Title(),
				),
			)
		}
		res, err = remote.Search(
			sourceHandler,
			p.Name,
			options,
		)
		if err != nil && !errors.Is(err, remote.ErrorNoResults) {
			logger.Fatal(err)
		}
	}

	if errors.Is(err, remote.ErrorNoResults) {
		logger.InfoNow("no results found")
	}

	tui.Flush(out)
	return nil
}

func appendToSearchOutput(
	out *tui.Data,
	showAll bool,
	res types.SearchResults,
) {
	var results []string
	for _, r := range res.Results {
		results = append(results, r.String())
	}

	if len(out.Fields) != 0 {
		out.Fields = append(
			out.Fields, &tui.FieldSeparator{
				Length: 0,
				Dim:    false,
			},
		)
	}

	out.Fields = append(
		out.Fields,
		&tui.FieldAnnotation{
			Annotation: "Results from " + res.Source.Title(),
		},
	)

	if res.Source == types.Modrinth && len(res.Results) == 100 {
		out.Fields = append(
			out.Fields,
			&tui.FieldAnnotation{
				Annotation: "* only showing the top 100",
			},
		)
	}

	out.Fields = append(
		out.Fields,
		&tui.FieldShortText{
			Title: "#  ",
			Text:  strconv.Itoa(len(res.Results)),
		},
		&tui.FieldDynamicColumnLabels{
			Title:  ">>>",
			Labels: results,
			MaxLines: tools.Ternary(
				showAll,
				0,
				tools.TermHeight()-6,
			),
		},
	)
}
