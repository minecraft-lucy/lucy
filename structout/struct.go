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

// Package structout is a key-value based commandline output framework.
//
// The core of this package it the lucytypes.OutputData struct. It is an array
// of different types of fields that defines different types of output formats.
// The OutputData struct can be simply passed to the Flush function to output
// to the commandline.
//
// Note the field will not show if its content is empty
package structout

import (
	"strconv"
	"strings"

	"lucy/tools"
)

type Data struct {
	Fields []Field
}

type Field interface {
	Output()
}

func SourceInfo(source string) {
	annot("(Source: " + tools.Underline(source) + ")")
	newLine()
}

// Separator prints a separator line. A length of 0 will print a line of 66%
// terminal width.
//
// Separator also adjusts itself so it does not exceed the terminal width.
//
// Use dim to control whether the separator is dimmed.
func Separator(len int, dim bool) {
	if len == 0 {
		len = tools.TermWidth() * 2 / 3
	} else if len > tools.TermWidth() {
		len = tools.TermWidth()
	}

	sep := strings.Repeat("-", len)
	if dim {
		annot(sep)
	} else {
		value(sep)
	}
	newLine()
}

type FieldSeparator struct {
	Length int
	Dim    bool
}

func (f *FieldSeparator) Output() {
	Separator(f.Length, f.Dim)
}

type FieldAnnotation struct {
	Annotation string
}

func (f *FieldAnnotation) Output() {
	annot(f.Annotation)
	newLine()
}

type FieldShortText struct {
	Title string
	Text  string
}

func (f *FieldShortText) Output() {
	key(f.Title)
	value(f.Text)
	newLine()
}

type FieldAnnotatedShortText struct {
	Title      string
	Text       string
	Annotation string
	NoTab      bool
}

func (f *FieldAnnotatedShortText) Output() {
	key(f.Title)
	value(f.Text)
	if f.NoTab {
		value("  " + tools.Dim(f.Annotation))
	} else {
		inlineAnnot(f.Annotation)
	}
	newLine()
}

var FieldNil = &fieldNil{}

type fieldNil struct{}

func (f *fieldNil) Output() {}

// FieldLabels is a field that contains a title and a list of labels. If the
// maxWidth is 0, it defaults to max(33% of terminal width, 40)
type FieldLabels struct {
	Title    string
	Labels   []string
	MaxWidth int
	MaxLines int
}

func (f *FieldLabels) Output() {
	if len(f.Labels) == 0 {
		return
	}

	key(f.Title)
	if f.MaxWidth == 0 {
		f.MaxWidth = max(33*tools.TermWidth()/100, 40)
	}

	width := 0
	lines := 1
	for i, label := range f.Labels {
		value(label)
		if i != len(f.Labels)-1 {
			value(", ")
		}
		width += len(label) + 2
		if width >= f.MaxWidth && i != len(f.Labels)-1 {
			newLine()
			tab()
			width = 0
			lines++
			if f.MaxLines != 0 && lines > f.MaxLines {
				annot("(" + strconv.Itoa(len(f.Labels)-i-1) + " more, use -e to show all)")
				newLine()
				break
			}
		}
	}

	if width != 0 {
		newLine()
	}
}

type FieldDynamicColumnLabels struct {
	Title     string
	Labels    []string
	MaxLines  int
	ShowTotal bool
}

func (f *FieldDynamicColumnLabels) Output() {
	if len(f.Labels) == 0 {
		return
	}

	// This field should have a unique indent size so it's predictable. Therefore
	// we call flush() before output.
	flush()
	key(f.Title)

	lines := 1
	maxLabelLen := 0
	for _, label := range f.Labels {
		if len(label) > maxLabelLen {
			maxLabelLen = len(label)
		}
	}

	columns := (tools.TermWidth() - 4) / (maxLabelLen + 2)
	if columns <= 0 {
		columns = 1
	}

	for i, label := range f.Labels {
		lastInRow := (i+1)%columns == 0
		lastAmongAll := i == len(f.Labels)-1
		value(label)

		if f.MaxLines != 0 && lines == f.MaxLines && lastInRow {
			newLine()
			tab()
			if f.ShowTotal {
				annot("(" + strconv.Itoa(len(f.Labels)) + "in total, " + strconv.Itoa(len(f.Labels)-i-1) + " more)")
			} else {
				annot("(" + strconv.Itoa(len(f.Labels)-i-1) + " more)")
			}
			break
		}

		if lastAmongAll && f.ShowTotal {
			if lastInRow {
				newLine()
			}
			tab()
			annot("(" + strconv.Itoa(len(f.Labels)) + " total)")
			break
		}

		if lastInRow {
			newLine()
			lines++
		}
		tab()
	}

	// After output, we call flush() to reset the indent size.
	flush()
}

// FieldMultiShortTextWithAnnot accepts 2 arrays, Texts and Annots. len(Texts) determines
// the length of the output. Any content in Annots after len(Texts) will be omitted.
type FieldMultiShortTextWithAnnot struct {
	Title     string
	Texts     []string
	Annots    []string
	ShowTotal bool
}

func (f *FieldMultiShortTextWithAnnot) Output() {
	if len(f.Texts) == 0 {
		return
	}

	for i, t := range f.Texts {
		// key() has a tab() at the beginning, so we skip the first tab
		if i == 0 {
			key(f.Title)
		} else {
			tab()
		}
		value(t)
		if f.Annots != nil && i < len(f.Annots) {
			inlineAnnot(f.Annots[i])
		}
		newLine()
	}

	if f.ShowTotal {
		tab()
		annot("(" + strconv.Itoa(len(f.Texts)) + " total)")
		newLine()
	}
}

type FieldMultiShortText struct {
	Title     string
	Texts     []string
	ShowTotal bool
}

func (f *FieldMultiShortText) Output() {
	if len(f.Texts) == 0 {
		return
	}

	for i, t := range f.Texts {
		if i == 0 {
			key(f.Title)
		} else {
			tab()
		}
		value(t)
		newLine()
	}

	if f.ShowTotal {
		tab()
		annot("(" + strconv.Itoa(len(f.Texts)) + " total)")
		newLine()
	}
}

// FieldCheckBox defaults to a red cross and green check when TrueText and
// FalseText is not specified.
type FieldCheckBox struct {
	Title     string
	Boolean   bool
	TrueText  string
	FalseText string
}

func (f *FieldCheckBox) Output() {
	key(f.Title)

	if f.TrueText == "" {
		f.TrueText = tools.Green("\u2713") // Check
	}
	if f.FalseText == "" {
		f.FalseText = tools.Red("\u2717") // X
	}

	if f.Boolean {
		value(f.TrueText)
	} else {
		value(f.FalseText)
	}
}

func Flush(data *Data) {
	for _, field := range data.Fields {
		if field != nil {
			field.Output()
		}
	}
	newLine()
	flush()
}
