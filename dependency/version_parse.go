package dependency

import (
	"strconv"
	"strings"

	"lucy/lucytypes"
)

// Parse is the main function to parse a RawVersion into a SemanticVersion.
//
// If the raw version is one of the special constants (which should be inferred
// before passing to this function), it returns InvalidVersion.
//
// It will attempt each type of version parsing, in order of specificity.
//
// If the label is not compatible with the version, return a semantic version
// that is labeled as Raw and contains the raw version as it is. This is to
// ensure a basic support to some non-standard versions numbers by only supporting
// Eq and Neq comparisons.
func Parse(
	raw lucytypes.RawVersion,
	label lucytypes.VersionLabel,
) lucytypes.SemanticVersion {
	switch raw {
	case lucytypes.LatestVersion, lucytypes.LatestCompatibleVersion, lucytypes.NoVersion, lucytypes.AllVersion, lucytypes.UnknownVersion:
		return lucytypes.InvalidVersion
	}
	switch label {
	case lucytypes.Semver:
		return parseSemver(string(raw))
	case lucytypes.MinecraftRelease:
		return parseMinecraftRelease(string(raw))
	case lucytypes.MinecraftSnapshot:
		return parseMinecraftSnapshot(string(raw))
	default:
		return lucytypes.InvalidVersion
	}
}

func parseSemver(s string) (v lucytypes.SemanticVersion) {
	return operatorPlus(s)
}

// These two are equivalent, for now.
func parseMinecraftRelease(s string) (v lucytypes.SemanticVersion) {
	return parseSemver(s)
}

func operatorPlus(s string) (v lucytypes.SemanticVersion) {
	tokens := strings.Split(s, "+")
	if len(tokens) >= 2 {
		s = strings.Join(tokens[:len(tokens)-1], "")
	}
	v = operatorDash(s)
	if v == lucytypes.InvalidVersion {
		return v
	}
	if len(tokens) >= 2 {
		v.Build = tokens[len(tokens)-1]
	}
	return v
}

func operatorDash(s string) (v lucytypes.SemanticVersion) {
	tokens := strings.Split(s, "-")
	if len(tokens) >= 2 {
		s = strings.Join(tokens[:len(tokens)-1], "")
	}
	v = operatorDot(s)
	if v == lucytypes.InvalidVersion {
		return v
	}
	if len(tokens) >= 2 {
		v.Prerelease = tokens[len(tokens)-1]
	}
	return v
}

func operatorDot(s string) (v lucytypes.SemanticVersion) {
	tokens := strings.Split(s, ".")
	if len(tokens) >= 2 {
		major, err := strconv.Atoi(tokens[0])
		if err != nil {
			return lucytypes.InvalidVersion
		}
		minor, err := strconv.Atoi(tokens[1])
		if err != nil {
			return lucytypes.InvalidVersion
		}
		v.Major = uint16(major)
		v.Minor = uint16(minor)
	}
	if len(tokens) == 3 {
		patch, err := strconv.Atoi(tokens[2])
		if err != nil {
			return lucytypes.InvalidVersion
		}
		v.Patch = uint16(patch)
	}
	return v
}

func parseMinecraftSnapshot(s string) lucytypes.SemanticVersion {
	return operatorInWeekIndex(s)
}

func operatorWeek(s string) (v lucytypes.SemanticVersion) {
	tokens := strings.Split(s, "w")
	if len(tokens) != 2 {
		return lucytypes.InvalidVersion
	}
	major, err := strconv.Atoi(tokens[0])
	if err != nil {
		return lucytypes.InvalidVersion
	}
	minor, err := strconv.Atoi(tokens[1])
	if err != nil {
		return lucytypes.InvalidVersion
	}
	v.Major = uint16(major)
	v.Minor = uint16(minor)
	return v
}

func operatorInWeekIndex(s string) (v lucytypes.SemanticVersion) {
	tokens := s[len(s)-1]
	v = operatorWeek(s[:len(s)-1])
	if v == lucytypes.InvalidVersion {
		return v
	}
	switch tokens {
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h':
		v.Patch = uint16(tokens)
	default:
		return lucytypes.InvalidVersion
	}
	return v
}
