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
	"lucy/logger"
	"lucy/lucytypes"
)

func SelectExecutable(
	executables []*lucytypes.ExecutableInfo,
	notes []Note,
) int {
	renewPromptTemplate()
	prompt := selectExecutableTemplate
	prompt.Details = generateNote(notes...)
	selectExecutable := promptui.Select{
		Label:     "Multiple possible executables detected, select one",
		Items:     executables,
		Templates: &prompt,
		Size:      max(8, len(executables)/3),
	}
	selection, _, err := selectExecutable.Run()
	if err != nil {
		logger.WarnNow(err)
	}
	return selection
}

func RememberExecutable() bool {
	confirmRememberExecutable := promptui.Prompt{
		Label:     "Remember this executable",
		IsConfirm: true,
	}
	result, _ := confirmRememberExecutable.Run()
	return result == "true"
}
