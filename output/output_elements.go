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

package output

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
	fmt.Fprintf(keyValueWriter, "%s\t", tools.Bold(tools.Magenta(title)))
}

func value(value string) {
	fmt.Fprintf(keyValueWriter, "%s", value)
}

func inlineAnnot(annotation string) {
	fmt.Fprintf(keyValueWriter, "\t%s", tools.Dim(annotation))
}

func annot(value string) {
	fmt.Fprintf(keyValueWriter, "%s", tools.Dim(value))
}

func newLine() {
	fmt.Fprintf(keyValueWriter, "\n")
}

func tab() {
	fmt.Fprintf(keyValueWriter, "%s\t", tools.Bold(tools.Magenta("")))
}

func flush() {
	keyValueWriter.Flush()
}
