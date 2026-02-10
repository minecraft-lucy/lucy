package prompt

import (
	"github.com/manifoldco/promptui"
	"lucy/logger"
	"lucy/types"
)

func SelectExecutable(
	executables []*types.ExecutableInfo,
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
