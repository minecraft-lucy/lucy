package cmd

import (
	"context"
	"fmt"
	"slices"

	"lucy/remote/source"

	"lucy/logger"
	"lucy/remote"
	"lucy/syntax"
	"lucy/tools"
	"lucy/tui"
	"lucy/types"

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
		decoratorGlobalFlags,
		decoratorHelpAndExitOnNoArg,
		decoratorLogAndExitOnError,
	),
}

var actionInfo cli.ActionFunc = func(
	ctx context.Context,
	cmd *cli.Command,
) error {
	id := syntax.Parse(cmd.Args().First())
	p := id.NewPackage()

	var out *tui.Data
	var err error

	switch id.Platform {
	case types.AllPlatform:
		for _, source := range source.All {
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

	case types.Fabric, types.Forge:
		info, err := remote.Information(source.Modrinth, id.Name)
		if err != nil {
			logger.ErrorNow(err)
		}
		p.Information = &info

		remote, err := remote.Fetch(source.Modrinth, id)
		p.Remote = &remote
		if err != nil {
			logger.ErrorNow(err)
			return err
		}
		out = infoOutput(p)
	case types.Mcdr:
		info, err := remote.Information(
			source.Mcdr,
			id.Name,
		)
		if err != nil {
			logger.WarnNow(err)
			break
		}
		remote, err := remote.Fetch(source.Mcdr, id)
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
	tui.Flush(out)
	return nil
}

// TODO: Link to newest version
// TODO: Link to latest compatible version
// TODO: Generate `lucy add` command

func infoOutput(p *types.Package) *tui.Data {
	o := &tui.Data{
		Fields: []tui.Field{
			&tui.FieldAnnotation{
				Annotation: "(from " + p.Remote.Source.Title() + ")",
			},
			&tui.FieldShortText{
				Title: "Name",
				Text:  p.Information.Title,
			},
			&tui.FieldShortText{
				Title: "Description",
				Text:  p.Information.Brief,
			},
			tools.Ternary[tui.Field](
				p.Information.DescriptionIsMarkdown,
				&tui.FieldMarkdown{
					Title:         "Information",
					Text:          p.Information.Description,
					Padding:       true,
					LineWrap:      true,
					MaxColumns:    tools.TermWidth() * 8 / 10,
					MaxLines:      tools.TermHeight() * 3 / 2,
					UseAlternate:  true,
					AlternateText: tools.Underline(p.Information.DescriptionUrl),
				},
				&tui.FieldLongText{
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
	for _, author := range p.Information.Authors {
		authorNames = append(authorNames, author.Name)
		authorLinks = append(authorLinks, author.Url)
	}

	o.Fields = append(
		o.Fields,
		&tui.FieldMultiAnnotatedShortText{
			Title:     "Authors",
			Texts:     authorNames,
			Annots:    authorLinks,
			ShowTotal: false,
		},
	)

	if p.Information != nil {
		o.Fields = append(
			o.Fields,
			&tui.FieldShortText{
				Title: "License",
				Text:  p.Information.License,
			},
		)
	}

	for _, url := range p.Information.Urls {
		o.Fields = append(
			o.Fields, &tui.FieldShortText{
				Title: url.Name,
				Text:  tools.Underline(url.Url),
			},
		)
	}

	o.Fields = append(
		o.Fields, &tui.FieldAnnotatedShortText{
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
		!slices.Contains(p.Supports.Platforms, types.Mcdr) {
		f := &tui.FieldLabels{
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
