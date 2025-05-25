package lucytypes

type Source uint8

const (
	CurseForge Source = iota
	Modrinth
	GitHub
	McdrCatalogue
	UnknownSource
	Auto
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
	"curseforge": CurseForge,
	"modrinth":   Modrinth,
	"github":     GitHub,
	"mcdr":       McdrCatalogue,
	"auto":       Auto,
	"":           Auto,
	"unknown":    UnknownSource,
}

func StringToSource(s string) Source {
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
