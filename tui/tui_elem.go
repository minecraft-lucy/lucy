package tui

import (
	"fmt"
	"os"
	"text/tabwriter"

	"lucy/tools"
)

const tabWriterDebug = false

var keyValueWriter = tabwriter.NewWriter(
	os.Stdout,
	0,
	0,
	2,
	' ',
	tools.Ternary(tabWriterDebug, tabwriter.Debug, 0),
)

func key(title string) {
	_, _ = fmt.Fprintf(keyValueWriter, "%s\t", tools.Bold(tools.Magenta(title)))
}

func value(value string) {
	_, _ = fmt.Fprintf(keyValueWriter, "%s", value)
}

func inlineAnnot(annotation string) {
	_, _ = fmt.Fprintf(keyValueWriter, "\t%s", tools.Dim(annotation))
}

func annot(value string) {
	_, _ = fmt.Fprintf(keyValueWriter, "%s", tools.Dim(value))
}

func newLine() {
	_, _ = fmt.Fprintf(keyValueWriter, "\n")
}

func tab() {
	_, _ = fmt.Fprintf(keyValueWriter, "%s\t", tools.Bold(tools.Magenta("")))
}

func flush() {
	_ = keyValueWriter.Flush()
}
