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
	"lucy/lucytypes"
	"lucy/tools"
)

type PromptNote string

var PromptNotePrefix = tools.Cyan("*") + " "

var SuspectPrePackagedServer PromptNote = "This is likely a pre-packaged server. Therefore, you might want to ignore the paths, and only look for the executable with your expected game version and mod loader."

var selectExecutableTemplate = &promptui.SelectTemplates{
	Active:   `{{ "●" | blue }} {{ .Path | bold }} [2m(Minecraft {{ .GameVersion }}, {{ if eq .Platform "minecraft" }}Vanilla{{ else }}{{ .Platform }} {{ .LoaderVersion }}{{ end }})[0m`,
	Inactive: `{{ "○" | blue }} {{ .Path }} [2m(Minecraft {{ .GameVersion }}, {{ if eq .Platform "minecraft" }}Vanilla{{ else }}{{ .Platform }} {{ .LoaderVersion }}{{ end }})[0m`,
	Selected: `{{ "✔︎" | green }} {{ .Path | bold }} [2m(Minecraft {{ .GameVersion }}, {{ if eq .Platform "minecraft" }}Vanilla{{ else }}{{ .Platform }} {{ .LoaderVersion }}{{ end }})[0m`,
}

func PromptSelectExecutable(
	executables []*lucytypes.ExecutableInfo,
	note []PromptNote,
) int {
	prompt := selectExecutableTemplate
	if note != nil {
		for _, n := range note {
			prompt.Details += PromptNotePrefix + string(n) + "\n"
		}
	}
	selectExecutable := promptui.Select{
		Label:     "Multiple possible executables detected, select one",
		Items:     executables,
		Templates: selectExecutableTemplate,
		Size:      max(8, len(executables)/3),
	}
	index, _, _ := selectExecutable.Run()
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
