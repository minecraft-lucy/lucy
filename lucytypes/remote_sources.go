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

import "strings"

type Source uint8

const (
	AutoSource Source = iota
	CurseForge
	Modrinth
	GitHub
	McdrCatalogue
	UnknownSource
)

func (s Source) String() string {
	switch s {
	case CurseForge:
		return "curseforge"
	case Modrinth:
		return "modrinth"
	case GitHub:
		return "github"
	case McdrCatalogue:
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
	case McdrCatalogue:
		return "MCDR"
	default:
		return "Unknown"
	}
}

var stringToSourceMap = map[string]Source{
	"auto":       AutoSource,
	"":           AutoSource,
	"curseforge": CurseForge,
	"modrinth":   Modrinth,
	"github":     GitHub,
	"mcdr":       McdrCatalogue,
	"unknown":    UnknownSource,
}

func StringToSource(s string) Source {
	s = strings.ToLower(s)
	if v, ok := stringToSourceMap[s]; ok {
		return v
	}
	return UnknownSource
}

type SearchOptions struct {
	ShowClientPackage bool
	IndexBy           SearchIndex
	Platform          Platform
}

type SearchIndex string

const (
	ByRelevance = "relevance"
	ByDownloads = "downloads"
	ByNewest    = "newest"
	ByName      = "name"
)

func (i SearchIndex) Valid() bool {
	switch i {
	case ByRelevance, ByDownloads, ByNewest:
		return true
	default:
		return false
	}
}

func (i SearchIndex) ToModrinth() string {
	switch i {
	case ByRelevance:
		return "relevance"
	case ByDownloads:
		return "downloads"
	case ByNewest:
		return "newest"
	default:
		return "relevance"
	}
}

type SearchResults struct {
	Source  Source
	Results []ProjectName
}
