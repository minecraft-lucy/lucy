package lucyversion

import (
	"lucy/lucytypes"
	"strconv"
	"strings"
)

func Parse(raw lucytypes.RawVersion, label VersionLabel) SemanticVersion {
	switch raw {
	case lucytypes.LatestVersion, lucytypes.LatestCompatibleVersion, lucytypes.NoVersion, lucytypes.AllVersion, lucytypes.UnknownVersion:
		return InvalidVersion
	}
	switch label {
	case Semver:
		return parseSemver(string(raw))
	case MinecraftRelease:
		return parseMinecraftRelease(string(raw))
	case MinecraftSnapshot:
		return parseMinecraftSnapshot(string(raw))
	default:
		return InvalidVersion
	}
}

func parseSemver(s string) (v SemanticVersion) {
	return operatorPlus(s)
}

// These two are equivalent, for now.
func parseMinecraftRelease(s string) (v SemanticVersion) {
	return parseSemver(s)
}

func operatorPlus(s string) (v SemanticVersion) {
	tokens := strings.Split(s, "+")
	if len(tokens) >= 2 {
		s = strings.Join(tokens[:len(tokens)-1], "")
	}
	v = operatorDash(s)
	if v == InvalidVersion {
		return v
	}
	if len(tokens) >= 2 {
		v.Build = tokens[len(tokens)-1]
	}
	return v
}

func operatorDash(s string) (v SemanticVersion) {
	tokens := strings.Split(s, "-")
	if len(tokens) >= 2 {
		s = strings.Join(tokens[:len(tokens)-1], "")
	}
	v = operatorDot(s)
	if v == InvalidVersion {
		return v
	}
	if len(tokens) >= 2 {
		v.Prerelease = tokens[len(tokens)-1]
	}
	return v
}

func operatorDot(s string) (v SemanticVersion) {
	tokens := strings.Split(s, ".")
	if len(tokens) >= 2 {
		major, err := strconv.Atoi(tokens[0])
		minor, err := strconv.Atoi(tokens[1])
		if err != nil {
			return InvalidVersion
		}
		v.Major = uint16(major)
		v.Minor = uint16(minor)
	}
	if len(tokens) == 3 {
		patch, err := strconv.Atoi(tokens[2])
		if err != nil {
			return InvalidVersion
		}
		v.Patch = uint16(patch)
	}
	return v
}

func parseMinecraftSnapshot(s string) SemanticVersion {
	return operatorInWeekIndex(s)
}

func operatorWeek(s string) (v SemanticVersion) {
	tokens := strings.Split(s, "w")
	if len(tokens) != 2 {
		return InvalidVersion
	}
	major, err := strconv.Atoi(tokens[0])
	minor, err := strconv.Atoi(tokens[1])
	if err != nil {
		return InvalidVersion
	}
	v.Major = uint16(major)
	v.Minor = uint16(minor)
	return v
}

func operatorInWeekIndex(s string) (v SemanticVersion) {
	tokens := s[len(s)-1]
	v = operatorWeek(s[:len(s)-1])
	if v == InvalidVersion {
		return v
	}
	switch tokens {
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h':
		v.Patch = uint16(tokens)
	default:
		return InvalidVersion
	}
	return v
}
