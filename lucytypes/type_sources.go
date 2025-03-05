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

package lucytypes

type Source uint8

const (
	CurseForge Source = iota
	Modrinth
	GitHub
	McdrWebsite
	UnknownSource
	// Auto is temporarily removed.
)

func (s Source) String() string {
	switch s {
	case CurseForge:
		return "curseforge"
	case Modrinth:
		return "modrinth"
	case GitHub:
		return "github"
	case McdrWebsite:
		return "mcdr"
	default:
		return "unknown"
	}
}

func (s Source) Title() string {
	switch s {
	case CurseForge:
		return "CurseForge"
	case Modrinth:
		return "Modrinth"
	case GitHub:
		return "GitHub"
	case McdrWebsite:
		return "MCDR"
	default:
		return "Unknown"
	}
}
