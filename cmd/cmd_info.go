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
	"slices"

	"lucy/remote/sources"

	"lucy/logger"
	"lucy/lucytypes"
	"lucy/remote"
	"lucy/structout"
	"lucy/syntax"
	"lucy/tools"

	"github.com/urfave/cli/v3"
)

var subcmdInfo = &cli.Command{
	Name:  "info",
	Usage: "Display information of a mod or plugin",
	Flags: []cli.Flag{
		flagSource,
		flagJsonOutput,
		flagNoStyle,
	},
	Action: tools.Decorate(
		actionInfo,
		globalFlagsDecorator,
		helpOnNoInputDecorator,
	),
}

var actionInfo cli.ActionFunc = func(
	ctx context.Context,
	cmd *cli.Command,
) error {
	id := syntax.Parse(cmd.Args().First())
	p := id.NewPackage()

	var out *structout.Data
	var err error

	switch id.Platform {
	case lucytypes.AllPlatform:
		for _, source := range sources.All {
			info, err := remote.Information(source, id.Name)
			if err != nil {
				continue
			}
			p.Information = &info

			remote, err := remote.Fetch(source, id)
			if err != nil {
				continue
			}
			p.Remote = &remote
			out = infoOutput(p)
			break
		}

	case lucytypes.Fabric, lucytypes.Forge:
		info, err := remote.Information(sources.Modrinth, id.Name)
		if err != nil {
			logger.ErrorNow(err)
		}
		p.Information = &info

		remote, err := remote.Fetch(sources.Modrinth, id)
		p.Remote = &remote
		if err != nil {
			logger.ErrorNow(err)
			return err
		}
		out = infoOutput(p)
	case lucytypes.Mcdr:
		info, err := remote.Information(
			sources.Mcdr,
			id.Name,
		)
		if err != nil {
			logger.WarnNow(err)
			break
		}
		remote, err := remote.Fetch(sources.Mcdr, id)
		if err != nil {
			logger.WarnNow(err)
			break
		}
		p.Information, p.Remote = &info, &remote
		out = infoOutput(p)
	}

	if err != nil {
		logger.Warn(err)
		return err
	}
	if out == nil {
		err = fmt.Errorf("%w: %s", remote.ErrorNoPackage, id.StringFull())
		logger.ErrorNow(err)
		return err
	}
	if cmd.Bool(flagJsonOutput.Name) {
		tools.PrintAsJson(p)
		return nil
	}
	structout.Flush(out)
	return nil
}

// TODO: Link to newest version
// TODO: Link to latest compatible version
// TODO: Generate `lucy add` command

func infoOutput(p *lucytypes.Package) *structout.Data {
	o := &structout.Data{
		Fields: []structout.Field{
			&structout.FieldAnnotation{
				Annotation: "(from " + p.Remote.Source.Title() + ")",
			},
			&structout.FieldShortText{
				Title: "Name",
				Text:  p.Information.Title,
			},
			&structout.FieldShortText{
				Title: "Description",
				Text:  p.Information.Brief,
			},
			tools.Ternary[structout.Field](
				p.Information.MarkdownDescription,
				&structout.FieldMarkdown{
					Title:         "Information",
					Text:          p.Information.Description,
					Padding:       true,
					LineWrap:      true,
					MaxColumns:    tools.TermWidth() * 8 / 10,
					MaxLines:      tools.TermHeight() * 3 / 2,
					UseAlternate:  true,
					AlternateText: tools.Underline(p.Information.DescriptionUrl),
				},
				&structout.FieldLongText{
					Title:         "Information",
					Text:          p.Information.Description,
					Padding:       true,
					LineWrap:      true,
					MaxColumns:    tools.TermWidth() * 8 / 10,
					MaxLines:      tools.TermHeight() * 3 / 2,
					UseAlternate:  true,
					AlternateText: tools.Underline(p.Information.DescriptionUrl),
				},
			),
		},
	}

	var authorNames []string
	var authorLinks []string
	for _, author := range p.Information.Author {
		authorNames = append(authorNames, author.Name)
		authorLinks = append(authorLinks, author.Url)
	}

	o.Fields = append(
		o.Fields,
		&structout.FieldMultiAnnotatedShortText{
			Title:     "Authors",
			Texts:     authorNames,
			Annots:    authorLinks,
			ShowTotal: false,
		},
	)

	if p.Information != nil {
		o.Fields = append(
			o.Fields,
			&structout.FieldShortText{
				Title: "License",
				Text:  p.Information.License,
			},
		)
	}

	for _, url := range p.Information.Urls {
		o.Fields = append(
			o.Fields, &structout.FieldShortText{
				Title: url.Name,
				Text:  tools.Underline(url.Url),
			},
		)
	}

	o.Fields = append(
		o.Fields, &structout.FieldAnnotatedShortText{
			Title:      "Download",
			Text:       tools.Underline(p.Remote.FileUrl),
			Annotation: p.Remote.Filename,
			NoTab:      true,
		},
	)

	// TODO: Put current server version on the top
	// TODO: Hide snapshot versions, except if the current server is using it
	if p.Supports != nil &&
		p.Supports.Platforms != nil &&
		!slices.Contains(p.Supports.Platforms, lucytypes.Mcdr) {
		f := &structout.FieldLabels{
			Title:    "Game Versions",
			Labels:   []string{},
			MaxWidth: 0,
			MaxLines: tools.TermHeight() / 2,
		}
		for _, version := range p.Supports.MinecraftVersions {
			f.Labels = append(f.Labels, version.String())
		}
		o.Fields = append(o.Fields, f)
	}

	return o
}
