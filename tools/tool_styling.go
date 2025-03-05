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

package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/term"
)

const (
	StyleReset = iota
	StyleBold
	StyleDim
	StyleItalic
	StyleUnderline
	StyleBlackText = iota + 25
	StyleRedText
	StyleGreenText
	StyleYellowText
	StyleBlueText
	StyleMagentaText
	StyleCyanText
	StyleWhiteText
)

const esc = '\u001B'

func styleFactory(i int) func(any) string {
	return func(v any) string {
		s := v.(string)
		return fmt.Sprintf("%c[%dm%s%c[%dm", esc, i, s, esc, StyleReset)
	}
}

func Capitalize(v any) string {
	s, ok := v.(string)
	if !ok {
		s = fmt.Sprintf("%v", v)
	}
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

var (
	Bold      = styleFactory(StyleBold)
	Dim       = styleFactory(StyleDim)
	Italic    = styleFactory(StyleItalic)
	Underline = styleFactory(StyleUnderline)
	Red       = styleFactory(StyleRedText)
	Green     = styleFactory(StyleGreenText)
	Yellow    = styleFactory(StyleYellowText)
	Blue      = styleFactory(StyleBlueText)
	Mangeta   = styleFactory(StyleMagentaText)
	Cyan      = styleFactory(StyleCyanText)
)

// PrintJson is usually used for debugging purposes
func PrintJson(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
}

func TermWidth() int {
	width, _, _ := term.GetSize(0)
	return width
}

func TermHeight() int {
	_, height, _ := term.GetSize(0)
	return height
}
