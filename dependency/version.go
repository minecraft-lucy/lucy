package dependency

// SemanticVersion is not exactly semver. However, I cannot find a better name
// for a structural version type. This will be mainly used in the construction of
// the dependency graph.
//
// For Minecraft Snapshots, Major is the year, Minor is the week of the year,
// and Patch is the rune at the end of the version string (to ascii code).
//
// In principle, you cannot compare two versions with different labels. This type
// of comparison always returns false.
//
// The StrictEq method is checks for Prerelease.
//
// Build is for display purposes only. It is not used in any conditional expressions.
//
// Patch is allowed to be zero for Minecraft releases, by this I mean the first
// release of each Minor, such as 1.19.
type SemanticVersion struct {
	Label      VersionLabel
	Major      uint16
	Minor      uint16
	Patch      uint16
	Prerelease string
	Build      string
}

type VersionLabel uint8

const (
	Semver VersionLabel = iota
	MinecraftSnapshot
	MinecraftRelease
	Invalid
)

var InvalidVersion = SemanticVersion{
	Label: Invalid,
}

func (v SemanticVersion) Validate() bool {
	switch v.Label {
	case Semver:
		return v.Major != 0 || v.Minor != 0 || v.Patch != 0
	case MinecraftSnapshot:
		return v.Major != 0 &&
			v.Minor != 0 && v.Minor <= maxWeek &&
			v.Patch >= minInWeekIndex && v.Patch <= maxInWeekIndex
	case MinecraftRelease:
		return v.Major != 0 && v.Minor != 0
	case Invalid:
		return false
	default:
		return false

	}
}

const maxWeek uint16 = 52 + 2
const maxInWeekIndex = uint16('h')
const minInWeekIndex = uint16('a')
