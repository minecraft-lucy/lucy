package prompt

import "lucy/tools"

type Note string

var (
	NotePrefix                        = func() string { return tools.Cyan("*") }
	NoteSuspectPrePackagedServer Note = "This is likely a pre-packaged server. Therefore, you might want to ignore the paths, and only look for the executable with your expected game version and mod loader."
)

func generateNote(notes ...Note) string {
	var note string
	for _, n := range notes {
		note += NotePrefix() + " " + string(n) + "\n"
	}
	return note
}
