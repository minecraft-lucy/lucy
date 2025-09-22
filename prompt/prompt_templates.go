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

package prompt

import (
	"github.com/manifoldco/promptui"
	"lucy/tools"
)

// Prompt strings are lazy-loaded, so they are not initialized until they are used.
// This is to support the `--no-style` flag, which disables all styling by reloading
// the styling functions in tools with a no-op.

var (
	promptSelectedIcon       = func() string { return tools.Green('✔') }
	promptActiveIcon         = func() string { return tools.Blue('●') }
	promptInactiveIcon       = func() string { return tools.Blue('○') }
	promptExecutablePath     = func() string { return tools.Bold("{{ .Path }}") }
	promptGameInfoAnnotation = func() string {
		return tools.Dim(`(Minecraft {{ .GameVersion }}, {{ if eq .LoaderPlatform "minecraft" }}Vanilla{{ else }}{{ .LoaderPlatform }} {{ .LoaderVersion }}{{ end }})`)
	}
	promptSelectContent = func() string { return promptExecutablePath() + " " + promptGameInfoAnnotation() }
)

var selectExecutableTemplate promptui.SelectTemplates

func renewPromptTemplate() {
	selectExecutableTemplate = promptui.SelectTemplates{
		Active:   promptActiveIcon() + " " + promptSelectContent(),
		Inactive: promptInactiveIcon() + " " + promptSelectContent(),
		Selected: promptSelectedIcon() + " " + promptSelectContent(),
	}
}
