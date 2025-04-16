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

func init() {
	renewStyleFunctions()
}

const (
	styleReset = iota
	styleBold
	styleDim
	styleItalic
	styleUnderline
	styleBlackText = iota + 25
	styleRedText
	styleGreenText
	styleYellowText
	styleBlueText
	styleMagentaText
	styleCyanText
	styleWhiteText
)

const esc = '\u001B'

var (
	Bold      func(any) string
	Dim       func(any) string
	Italic    func(any) string
	Underline func(any) string
	Red       func(any) string
	Green     func(any) string
	Yellow    func(any) string
	Blue      func(any) string
	Magenta   func(any) string
	Cyan      func(any) string
)

func renewStyleFunctions() {
	Bold = styleFactory(styleBold)
	Dim = styleFactory(styleDim)
	Italic = styleFactory(styleItalic)
	Underline = styleFactory(styleUnderline)
	Red = styleFactory(styleRedText)
	Green = styleFactory(styleGreenText)
	Yellow = styleFactory(styleYellowText)
	Blue = styleFactory(styleBlueText)
	Magenta = styleFactory(styleMagentaText)
	Cyan = styleFactory(styleCyanText)
}

func TurnOffStyles() {
	styleFactory = func(i int) func(any) string {
		return func(v any) string {
			switch v := v.(type) {
			case rune:
				return string(v)
			default:
				return fmt.Sprintf("%v", v)
			}
		}
	}
	renewStyleFunctions()
}

var styleFactory = func(i int) func(any) string {
	return func(v any) string {
		var s string
		switch v := v.(type) {
		case rune:
			s = string(v)
		default:
			s = fmt.Sprintf("%v", v)
		}
		return fmt.Sprintf("%c[%dm%s%c[%dm", esc, i, s, esc, styleReset)
	}
}

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
