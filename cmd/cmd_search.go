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

	"lucy/lucytypes"
	"lucy/remote"
	"lucy/remote/sources"

	"github.com/urfave/cli/v3"
	"lucy/logger"
	"lucy/structout"
	"lucy/syntax"
	"lucy/tools"
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
				if lucytypes.SearchIndex(s).Valid() {
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
	Action: tools.Decorate(
		actionSearch,
		decoratorGlobalFlags,
		decoratorHelpAndExitOnNoInput,
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
	indexBy := lucytypes.SearchIndex(cmd.String("index"))
	options := lucytypes.SearchOptions{
		ShowClientPackage: showClientPackage,
		IndexBy:           indexBy,
	}
	sourceStr := cmd.String("source")
	source := lucytypes.StringToSource(sourceStr)

	out := &structout.Data{}
	res := lucytypes.SearchResults{}
	var err error

	if source == lucytypes.AutoSource {
		switch p.Platform {
		case lucytypes.AllPlatform:
			for _, sourceHandler := range sources.All {
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
		case lucytypes.Forge, lucytypes.Fabric, lucytypes.Neoforge:
			res, err = remote.Search(sources.Modrinth, p.Name, options)
			if err != nil && !errors.Is(err, remote.ErrorNoResults) {
				logger.Fatal(err)
			}
			appendToSearchOutput(out, cmd.Bool("long"), res)
		case lucytypes.Mcdr:
			res, err = remote.Search(sources.Mcdr, p.Name, options)
			if err != nil && !errors.Is(err, remote.ErrorNoResults) {
				logger.Fatal(err)
			}
			appendToSearchOutput(out, cmd.Bool("long"), res)
		case lucytypes.UnknownPlatform:
			logger.Fatal(
				fmt.Errorf(
					"%w: %s",
					errorInvalidPlatform,
					p.Platform,
				),
			)
		}
	} else {
		if source == lucytypes.UnknownSource {
			logger.Fatal(
				fmt.Errorf(
					"%w: %s",
					errorUnknownSource,
					sourceStr,
				),
			)
		}
		sourceHandler, ok := sources.Map[source]
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

	structout.Flush(out)
	return nil
}

func appendToSearchOutput(
	out *structout.Data,
	showAll bool,
	res lucytypes.SearchResults,
) {
	var results []string
	for _, r := range res.Results {
		results = append(results, r.String())
	}

	if len(out.Fields) != 0 {
		out.Fields = append(
			out.Fields, &structout.FieldSeparator{
				Length: 0,
				Dim:    false,
			},
		)
	}

	out.Fields = append(
		out.Fields,
		&structout.FieldAnnotation{
			Annotation: "Results from " + res.Source.Title(),
		},
	)

	if res.Source == lucytypes.Modrinth && len(res.Results) == 100 {
		out.Fields = append(
			out.Fields,
			&structout.FieldAnnotation{
				Annotation: "* only showing the top 100",
			},
		)
	}

	out.Fields = append(
		out.Fields,
		&structout.FieldShortText{
			Title: "#  ",
			Text:  strconv.Itoa(len(res.Results)),
		},
		&structout.FieldDynamicColumnLabels{
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
