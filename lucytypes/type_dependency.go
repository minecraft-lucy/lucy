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

// RawVersion is the version of a package. Here we expect mods and plugins
// use semver (which they should). A known exception is Minecraft snapshots.
//
// There are several special constant values for RawVersion. You MUST call
// remote.InferVersion() before parsing them to ComparableVersion.
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

func (p1 ComparableVersion) schemeMatch(p2 ComparableVersion) bool {
	return p1.Scheme == p2.Scheme
}

// Eq checks whether p1 is equal to p2.
func (p1 ComparableVersion) Eq(p2 ComparableVersion) bool {
	if !p1.schemeMatch(p2) {
		return false
	}
	return p1.Major == p2.Major && p1.Minor == p2.Minor && p1.Patch == p2.Patch
}

// Neq checks whether p1 is not equal to p2.
func (p1 ComparableVersion) Neq(p2 ComparableVersion) bool {
	if !p1.schemeMatch(p2) {
		return false
	}
	return !p1.Eq(p2)
}

// StrictEq checks whether p1 is strictly equal to p2. This includes
// the prerelease tag.
func (p1 ComparableVersion) StrictEq(p2 ComparableVersion) bool {
	if !p1.schemeMatch(p2) {
		return false
	}
	// Even in strict equality, we ignore the build.
	if p1.Scheme != p2.Scheme {
		return false
	}
	return p1.Eq(p2)
}

// Lt checks whether p1 is less than p2.
func (p1 ComparableVersion) Lt(p2 ComparableVersion) bool {
	if !p1.schemeMatch(p2) {
		return false
	}

	if p1.Major < p2.Major {
		return true
	} else if p1.Major > p2.Major {
		return false
	}
	if p1.Minor < p2.Minor {
		return true
	} else if p1.Minor > p2.Minor {
		return false
	}
	if p1.Patch < p2.Patch {
		return true
	}

	// Usually a version with a prerelease tag is considered lower than
	// the same version without a prerelease tag.
	if p1.Prerelease == "" {
		if p2.Prerelease != "" {
			return true
		}
		return false
	} else {
		if p2.Prerelease != "" {
			// lexicographical order
			return strings.Compare(p1.Prerelease, p2.Prerelease) < 0
		} else {
			// p1 has a prerelease tag, p2 does not.
			return false
		}
	}
}

// Gt checks whether p1 is greater than p2.
func (p1 ComparableVersion) Gt(p2 ComparableVersion) bool {
	if !p1.schemeMatch(p2) {
		return false
	}

	if p1.Major > p2.Major {
		return true
	} else if p1.Major < p2.Major {
		return false
	}
	if p1.Minor > p2.Minor {
		return true
	} else if p1.Minor < p2.Minor {
		return false
	}
	if p1.Patch > p2.Patch {
		return true
	}

	// Usually a version with a prerelease tag is considered lower than
	// the same version without a prerelease tag.
	if p1.Prerelease == "" {
		if p2.Prerelease != "" {
			return false
		}
		return false
	} else {
		if p2.Prerelease != "" {
			// lexicographical order
			return strings.Compare(p1.Prerelease, p2.Prerelease) > 0
		} else {
			// p1 has a prerelease tag, p2 does not.
			return true
		}
	}
}

// WeakGt is for being compatible with the '^' operator in semver. Like Gt,
// it checks whether p1 is greater than p2. However, it does not allow
// the major version to be different.
func (p1 ComparableVersion) WeakGt(p2 ComparableVersion) bool {
	if !p1.schemeMatch(p2) {
		return false
	}

	if p1.Major != p2.Major {
		return false
	}
	return p1.Gt(p2)
}

// Lte checks whether p1 is less than or equal to p2.
func (p1 ComparableVersion) Lte(p2 ComparableVersion) bool {
	if !p1.schemeMatch(p2) {
		return false
	}

	return p1.Lt(p2) || p1.Eq(p2)
}

// Gte checks whether p1 is greater than or equal to p2.
func (p1 ComparableVersion) Gte(p2 ComparableVersion) bool {
	if !p1.schemeMatch(p2) {
		return false
	}

	return p1.Gt(p2) || p1.Eq(p2)
}

// ComparableVersion is a structural representation of a version (in numbers and
// partial strings). Note that it do not include the package name.
//
// For Minecraft Snapshots, Major is the year, Minor is the week of the year,
// and Patch is the rune at the end of the version string (to ascii code).
//
// In principle, you cannot compare two versions with different schema. This
// type of comparison always returns false.
//
// The StrictEq method is checks for Prerelease.
//
// Build is for recognition purposes only. It is not used in any conditional expressions.
//
// Patch is allowed to be zero for Minecraft releases, by this I mean the first
// release of each Minor, such as 1.19.
type ComparableVersion struct {
	Scheme     VersionScheme // The type of versioning scheme used.
	Major      uint16
	Minor      uint16
	Patch      uint16
	Prerelease string
	Build      string
}

type VersionScheme uint8

const (
	Semver VersionScheme = iota
	MinecraftSnapshot
	MinecraftRelease
	Invalid
)

var InvalidVersion = ComparableVersion{
	Scheme: Invalid,
}

func (v ComparableVersion) Validate() bool {
	switch v.Scheme {
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

type Dependency struct {
	Id           PackageId
	Requirements []Requirement
}

type Requirement struct {
	Value    ComparableVersion
	Operator VersionOperator
}

func (d Dependency) Satisfy(
	id PackageId,
	v ComparableVersion,
) bool {
	if (id.Platform != d.Id.Platform) || (id.Name != d.Id.Name) {
		return false
	}
	for _, req := range d.Requirements {
		switch req.Operator {
		case Equal:
			if !v.Eq(req.Value) {
				return false
			}
		case NotEqual:
			if !v.Neq(req.Value) {
				return false
			}
		case GreaterThan:
			if !v.Gt(req.Value) {
				return false
			}
		case WeakGreaterThan:
			if !v.WeakGt(req.Value) {
				return false
			}
		case GreaterThanOrEqual:
			if !v.Gte(req.Value) {
				return false
			}
		case LessThan:
			if !v.Lt(req.Value) {
				return false
			}
		case LessThanOrEqual:
			if !v.Lte(req.Value) {
				return false
			}
		}
	}

	return true
}

type VersionOperator uint8

const (
	Equal VersionOperator = iota
	NotEqual
	GreaterThan
	WeakGreaterThan // for ^ operator in semver
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
)
