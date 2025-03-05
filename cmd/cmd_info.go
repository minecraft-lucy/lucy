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
	"lucy/lucyerrors"
	"slices"

	"lucy/datatypes"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/output"
	"lucy/remote"
	"lucy/remote/mcdr"
	"lucy/syntax"
	"lucy/tools"

	"github.com/urfave/cli/v3"
)

var subcmdInfo = &cli.Command{
	Name:  "info",
	Usage: "Display information of a mod or plugin",
	Flags: []cli.Flag{
		sourceFlag(lucytypes.Modrinth),
		flagJsonOutput,
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

	var out *output.Data
	var err error

	switch id.Platform {
	case lucytypes.AllPlatform:
		p.Information, err = remote.Information(lucytypes.Modrinth, id.Name)
		p.Remote, err = remote.Fetch(lucytypes.Modrinth, id)
		if err == nil {
			out = infoOutput(p)
			break
		}
		p.Information, err = remote.Information(lucytypes.McdrWebsite, id.Name)
		p.Remote, err = remote.Fetch(lucytypes.McdrWebsite, id)
		if err == nil {
			out = infoOutput(p)
			break
		}
		err = fmt.Errorf("%w: %s", lucyerrors.ENotFound, id.StringFull())
		logger.ErrorNow(err)
		return err
	case lucytypes.Fabric, lucytypes.Forge:
		p.Information, err = remote.Information(lucytypes.Modrinth, id.Name)
		if err != nil {
			logger.ErrorNow(err)
		}
		p.Remote, err = remote.Fetch(lucytypes.Modrinth, id)
		if err != nil {
			logger.ErrorNow(err)
			return err
		}
		out = infoOutput(p)
	case lucytypes.Mcdr:
		mcdrPlugin, err := mcdr.SearchMcdrPluginCatalogue(id.Name)
		if err != nil {
			logger.Warn(err)
			break
		}
		out = mcdrPluginInfoToInfo(mcdrPlugin)
	}
	if err != nil {
		logger.Warn(err)
		return err
	}
	if cmd.Bool(flagJsonOutput.Name) {
		tools.PrintJson(p)
		return nil
	}
	output.Flush(out)
	return nil
}

// TODO: Link to newest version
// TODO: Link to latest compatible version
// TODO: Generate `lucy add` command

func mcdrPluginInfoToInfo(source *datatypes.McdrPluginInfo) *output.Data {
	info := &output.Data{
		Fields: []output.Field{
			&output.FieldShortText{
				Title: "Name",
				Text:  source.Id,
			},
			&output.FieldShortText{
				Title: "Description",
				Text:  source.Introduction.EnUs,
			},
			&output.FieldMultiShortTextWithAnnot{
				Title:  "Authors",
				Texts:  []string{},
				Annots: []string{},
			},
			&output.FieldShortText{
				Title: "Source Code",
				Text:  tools.Underline(source.Repository),
			},
		},
	}

	// This is temporary TODO: Use iota for fields instead
	const authorsField = 2
	a := info.Fields[authorsField].(*output.FieldMultiShortTextWithAnnot)

	for _, p := range source.Authors {
		a.Texts = append(a.Texts, p.Name)
		a.Annots = append(a.Annots, tools.Underline(p.Link))
	}

	return info
}

func infoOutput(p *lucytypes.Package) *output.Data {
	o := &output.Data{
		Fields: []output.Field{
			&output.FieldAnnotation{
				Annotation: "(from " + p.Remote.Source.Title() + ")",
			},
			&output.FieldShortText{
				Title: "Name",
				Text:  p.Information.Title,
			},
			&output.FieldShortText{
				Title: "Description",
				Text:  p.Information.Brief,
			},
		},
	}

	var authorNames []string
	var authorLinks []string
	for _, author := range p.Information.Author {
		authorNames = append(authorNames, author.Name)
		// TODO: Improve author info annotation format
		authorLinks = append(authorLinks, author.Url)
	}

	for _, url := range p.Information.Urls {
		o.Fields = append(
			o.Fields, &output.FieldShortText{
				Title: url.Name,
				Text:  tools.Underline(url.Url),
			},
		)
	}

	o.Fields = append(
		o.Fields, &output.FieldAnnotatedShortText{
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
		f := &output.FieldLabels{
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
