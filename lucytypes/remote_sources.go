package lucytypes

type Source uint8

const (
	CurseForge Source = iota
	Modrinth
	GitHub
	McdrWebsite
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

var stringToSourceMap = map[string]Source{
	"curseforge": CurseForge,
	"modrinth":   Modrinth,
	"github":     GitHub,
	"mcdr":       McdrWebsite,
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
