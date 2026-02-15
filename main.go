// TODO: REPLACE ALL io.ReadAll WITH STREAMING METHODS

package main

import (
	"context"
	"os"

	"lucy/cmd"
	"lucy/logger"
)

func main() {
	defer logger.DumpHistory()
	if err := cmd.Cli.Run(context.Background(), os.Args); err != nil {
		logger.ReportError(err)
	}
}
