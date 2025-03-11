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

package structout

import (
	"github.com/manifoldco/promptui"
	"lucy/lnout"
	"lucy/lucytypes"
	"lucy/tools"
)

// Prompt strings are lazy-loaded, so they are not initialized until they are used.
// This is to support the `--no-style` flag, which disables all styling by reloading
// the styling functions in tools with a no-op.

type PromptNote string

var (
	PromptNotePrefix                    = func() string { return tools.Cyan("*") }
	SuspectPrePackagedServer PromptNote = "This is likely a pre-packaged server. Therefore, you might want to ignore the paths, and only look for the executable with your expected game version and mod loader."
)

var (
	promptSelectedIcon       = func() string { return tools.Green('✔') }
	promptActiveIcon         = func() string { return tools.Blue('●') }
	promptInactiveIcon       = func() string { return tools.Blue('○') }
	promptExecutablePath     = func() string { return tools.Bold("{{ .Path }}") }
	promptGameInfoAnnotation = func() string {
		return tools.Dim(`(Minecraft {{ .GameVersion }}, {{ if eq .Platform "minecraft" }}Vanilla{{ else }}{{ .Platform }} {{ .LoaderVersion }}{{ end }})`)
	}
	promptSelectContent = func() string { return promptExecutablePath() + " " + promptGameInfoAnnotation() }
)

var selectExecutableTemplate *promptui.SelectTemplates

func renewPromptTemplate() {
	selectExecutableTemplate = &promptui.SelectTemplates{
		Active:   promptActiveIcon() + " " + promptSelectContent(),
		Inactive: promptInactiveIcon() + " " + promptSelectContent(),
		Selected: promptSelectedIcon() + " " + promptSelectContent(),
	}
}

func PromptSelectExecutable(
executables []*lucytypes.ExecutableInfo,
note []PromptNote,
) int {
	renewPromptTemplate()
	prompt := selectExecutableTemplate
	if note != nil {
		for _, n := range note {
			prompt.Details += PromptNotePrefix() + " " + string(n) + "\n"
		}
	}
	selectExecutable := promptui.Select{
		Label:     "Multiple possible executables detected, select one",
		Items:     executables,
		Templates: selectExecutableTemplate,
		Size:      max(8, len(executables)/3),
	}
	index, _, err := selectExecutable.Run()
	if err != nil {
		lnout.WarnNow(err)
	}
	return index
}

func PromptRememberExecutable() bool {
	confirmRememberExecutable := promptui.Prompt{
		Label:     "Remember this executable",
		IsConfirm: true,
	}
	result, _ := confirmRememberExecutable.Run()
	return result == "true"
}
