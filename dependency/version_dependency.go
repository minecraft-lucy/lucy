package dependency

type VersionOperator uint8

const (
	Equal VersionOperator = iota
	NotEqual
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
)

// RawVersion is the version of a package. Here we expect mods and plugins
// use semver (which they should). A known exception is Minecraft snapshots.
//
// There are several special constant values for RawVersion. You MUST call
// remote.InferVersion() before parsing them to SemanticVersion.
type RawVersion string

func (v RawVersion) String() string {
	if v == AllVersion {
		return "any"
	}
	if v == NoVersion || v == "" {
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

var (
	AllVersion              RawVersion = "all"
	NoVersion               RawVersion = "none"
	UnknownVersion          RawVersion = "unknown"
	LatestVersion           RawVersion = "latest"
	LatestCompatibleVersion RawVersion = "compatible"
)
