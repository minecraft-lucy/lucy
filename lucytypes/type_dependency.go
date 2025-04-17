package lucytypes

// RawVersion is the version of a package. Here we expect mods and plugins
// use semver (which they should). A known exception is Minecraft snapshots.
//
// There are several special constant values for RawVersion. You MUST call
// remote.InferVersion() before parsing them to SemanticVersion.
type RawVersion string

func (v RawVersion) String() string {
	if v == AllVersion || v == "" {
		return "any"
	}
	if v == NoVersion {
		return "none"
	}
	if v == UnknownVersion {
		return "unknown"
	}
	if v == LatestVersion {
		return "latest"
	}
	if v == LatestCompatibleVersion {
		return "compatible"
	}
	return string(v)
}

func (v RawVersion) NeedsInfer() bool {
	if v == AllVersion || v == NoVersion || v == UnknownVersion ||
		v == LatestVersion || v == LatestCompatibleVersion {
		return true
	}
	return false
}

var (
	AllVersion              RawVersion = "all"
	NoVersion               RawVersion = "none"
	UnknownVersion          RawVersion = "unknown"
	LatestVersion           RawVersion = "latest"
	LatestCompatibleVersion RawVersion = "compatible"
)

func (p1 SemanticVersion) Eq(p2 SemanticVersion) bool {
	// If the labels are different, the versions are not comparable.
	if p1.Label != p2.Label {
		return false
	}
	return p1.Major == p2.Major && p1.Minor == p2.Minor && p1.Patch == p2.Patch
}

func (p1 SemanticVersion) Neq(p2 SemanticVersion) bool {
	return !p1.Eq(p2)
}

func (p1 SemanticVersion) StrictEq(p2 SemanticVersion) bool {
	// Even in strict equality, we ignore the build.
	if p1.Label != p2.Label {
		return false
	}
	return p1.Major == p2.Major && p1.Minor == p2.Minor && p1.Patch == p2.Patch && p1.Prerelease == p2.Prerelease
}

func (p1 SemanticVersion) Lt(p2 SemanticVersion) bool {
	if p1.Major < p2.Major {
		return true
	}
	if p1.Major > p2.Major {
		return false
	}
	if p1.Minor < p2.Minor {
		return true
	}
	if p1.Minor > p2.Minor {
		return false
	}
	return p1.Patch < p2.Patch
}

func (p1 SemanticVersion) Gt(p2 SemanticVersion) bool {
	if p1.Major > p2.Major {
		return true
	}
	if p1.Major < p2.Major {
		return false
	}
	if p1.Minor > p2.Minor {
		return true
	}
	if p1.Minor < p2.Minor {
		return false
	}
	return p1.Patch > p2.Patch
}

func (p1 SemanticVersion) Lte(p2 SemanticVersion) bool {
	return p1.Lt(p2) || p1.Eq(p2)
}

func (p1 SemanticVersion) Gte(p2 SemanticVersion) bool {
	return p1.Gt(p2) || p1.Eq(p2)
}

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

const (
	maxWeek        uint16 = 52 + 2
	maxInWeekIndex        = uint16('h')
	minInWeekIndex        = uint16('a')
)

// Dependency can describe a dependency relationship. You MUST NOT use the
// Id's PackageId.Version field. Instead, you should use the Value and Operator.
type Dependency struct {
	Id       PackageId
	Value    SemanticVersion
	Operator VersionOperator
}

func (d Dependency) Satisfy(
	id PackageId,
	v SemanticVersion,
) bool {
	if (id.Platform != d.Id.Platform) || (id.Name != d.Id.Name) {
		return false
	}
	switch d.Operator {
	case Equal:
		return v.Eq(d.Value)
	case NotEqual:
		return v.Neq(d.Value)
	case GreaterThan:
		return v.Gt(d.Value)
	case GreaterThanOrEqual:
		return v.Gte(d.Value)
	case LessThan:
		return v.Lt(d.Value)
	case LessThanOrEqual:
		return v.Lte(d.Value)
	default:
		return false
	}
}

type VersionOperator uint8

const (
	Equal VersionOperator = iota
	NotEqual
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
)
