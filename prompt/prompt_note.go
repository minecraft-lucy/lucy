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
