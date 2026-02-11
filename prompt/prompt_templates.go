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
		return tools.Dim(`(Minecraft {{ .GameVersion }}, {{ if eq .ModLoader "minecraft" }}Vanilla{{ else }}{{ .ModLoader }} {{ .LoaderVersion }}{{ end }})`)
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
