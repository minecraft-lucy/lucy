package cli

import (
	"context"
	"github.com/urfave/cli/v3"
	"strings"
)

// Frontend
// This changes when user runs the web interface
var Frontend = "cli"

var Cli = &cli.Command{

	Name:   "lucy",
	Usage:  "The Minecraft server-side package manager",
	Action: noArgAction,
	Commands: []*cli.Command{
		SubcmdStatus,
		SubcmdInfo,
		SubcmdSearch,
	},
}

// This shows the help message of the called command
func noArgAction(_ context.Context, cmd *cli.Command) error {
	cli.ShowAppHelpAndExit(cmd, 0)
	return nil
}

// Parse the platform/package syntax
func parsePackageSyntax(query string) (platform string, packageName string) {
	split := strings.Split(query, "/")
	if len(split) == 1 {
		return "all", split[0]
	} else if len(split) == 2 {
		return split[0], split[1]
	} else {
		return "", ""
	}
}
