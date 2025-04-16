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
	"log"
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
		flagSource(lucytypes.AutoSource),
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
	},
	Action: tools.Decorate(
		actionSearch,
		globalFlagsDecorator,
		helpOnNoInputDecorator,
	),
}

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
	source := cmd.String("source")

	var res = lucytypes.SearchResults{}
	var err error
	switch source {
	case lucytypes.AutoSource.String():
		switch p.Platform {
		case lucytypes.AllPlatform:
			res, err = remote.Search(sources.Modrinth, p.Name, options)
		case lucytypes.Forge, lucytypes.Fabric, lucytypes.Neoforge:
			res, err = remote.Search(sources.Modrinth, p.Name, options)
		case lucytypes.Mcdr:
			res, err = remote.Search(sources.Mcdr, p.Name, options)
		case lucytypes.UnknownPlatform:
			log.Fatal("unknown platform")
		}
	default:
		res, err = remote.Search(
			sources.Map[lucytypes.StringToSource(cmd.String("source"))],
			p.Name,
			options,
		)

	}
	if err != nil {
		logger.Fatal(err)
	}
	structout.Flush(generateSearchOutput(res, cmd.Bool("long")))
	return nil
}

func generateSearchOutput(
	res lucytypes.SearchResults,
	showAll bool,
) *structout.Data {
	var results []string
	for _, r := range res.Results {
		results = append(results, r.String())
	}

	return &structout.Data{
		Fields: []structout.Field{
			&structout.FieldShortText{
				Title: "#  ",
				Text:  strconv.Itoa(len(res.Results)),
			},
			&structout.FieldDynamicColumnLabels{
				Title:    ">>>",
				Labels:   results,
				MaxLines: tools.Ternary(showAll, 0, tools.TermHeight()-6),
			},
		},
	}
}
