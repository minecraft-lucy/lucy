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

package dependency

import (
	"strconv"
	"strings"

	"lucy/lucytypes"
)

// Parse is the main function to parse a RawVersion into a ComparableVersion.
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
	scheme lucytypes.VersionScheme,
) lucytypes.ComparableVersion {
	switch raw {
	case lucytypes.LatestVersion, lucytypes.LatestCompatibleVersion, lucytypes.NoVersion, lucytypes.AllVersion, lucytypes.UnknownVersion:
		return lucytypes.InvalidVersion
	}
	switch scheme {
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

func parseSemver(s string) (v lucytypes.ComparableVersion) {
	return operatorPlus(s)
}

// These two are equivalent, for now.
func parseMinecraftRelease(s string) (v lucytypes.ComparableVersion) {
	return parseSemver(s)
}

func operatorPlus(s string) (v lucytypes.ComparableVersion) {
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

func operatorDash(s string) (v lucytypes.ComparableVersion) {
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

func operatorDot(s string) (v lucytypes.ComparableVersion) {
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

func parseMinecraftSnapshot(s string) lucytypes.ComparableVersion {
	return operatorInWeekIndex(s)
}

func operatorWeek(s string) (v lucytypes.ComparableVersion) {
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

func operatorInWeekIndex(s string) (v lucytypes.ComparableVersion) {
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
