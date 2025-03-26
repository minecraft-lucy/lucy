package remote

import "lucy/lucytypes"

type SearchOptions struct {
	ShowClientPackage bool
	IndexBy           SearchIndex
	Platform          lucytypes.Platform
}

type SearchIndex string

const (
	ByRelevance = "relevance"
	ByDownloads = "downloads"
	ByNewest    = "newest"
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
	Source  lucytypes.Source
	Results []lucytypes.ProjectName
}
