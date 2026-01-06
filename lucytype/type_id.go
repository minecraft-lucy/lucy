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

// Package lucytypes is a general package for all types used in Lucy.
package lucytype

import (
	"strings"

	"lucy/tools"
)

// Platform is an enum of several string constants. All platform is a package under
// itself, for example, "fabric/fabric" is a valid package, and is equivalent to
// "fabric". This literal is typically used when installing/upgrading a platform
// itself.
type Platform string

const (
	AllPlatform     Platform = ""
	Minecraft       Platform = "minecraft"
	Fabric          Platform = "fabric"
	Forge           Platform = "forge"
	Neoforge        Platform = "neoforge"
	Mcdr            Platform = "mcdr"
	UnknownPlatform Platform = "unknown"
)

func (p Platform) Title() string {
	if p == AllPlatform {
		return "Any"
	}
	if p.Valid() {
		return strings.ToUpper(string(p)[0:1]) + string(p)[1:]
	}
	return "Unknown"
}

func (p Platform) String() string {
	if p == AllPlatform {
		return "any"
	}
	return string(p)
}

// Valid should be edited if you added a new platform.
func (p Platform) Valid() bool {
	switch p {
	case Minecraft, Fabric, Forge, Neoforge, Mcdr, AllPlatform, UnknownPlatform:
		return true
	}
	return false
}

func (p Platform) Eq(other Platform) bool {
	if p == AllPlatform || other == AllPlatform {
		return true
	}
	if p == UnknownPlatform || other == UnknownPlatform {
		return false
	}
	return p == other
}

// ProjectName is the slug of the package, using hyphens as separators. For example,
// "fabric-api".
//
// It is non-case-sensitive, though lowercase is recommended. Underlines '_' are
// equivalent to hyphens.
//
// A slug from an upstream API is preferred, if possible. Otherwise, the slug is
// obtained from the executable file. No exceptions since a package must either
// exist on a remote API or user's local files.
type ProjectName string

func ToProjectName(s string) ProjectName {
	var b strings.Builder
	b.Grow(len(s))

	for _, char := range s {
		switch {
		case char == '_':
			b.WriteByte('-')
		case 'A' <= char && char <= 'Z':
			b.WriteRune(char + 'a' - 'A')
		default:
			b.WriteRune(char)
		}
	}

	return ProjectName(b.String())
}

// Title Replaces underlines or hyphens with spaces, then capitalize the first
// letter.
func (n ProjectName) Title() string {
	return tools.Capitalize(strings.ReplaceAll(string(n), "-", " "))
}

func (n ProjectName) String() string {
	return string(n)
}

func (n ProjectName) ToPEP8() string {
	return strings.ReplaceAll(string(n), "-", "_")
}

type PackageId struct {
	Platform Platform
	Name     ProjectName
	Version  RawVersion
}

func (p PackageId) NewPackage() *Package {
	return &Package{
		Id: PackageId{
			Platform: p.Platform,
			Name:     p.Name,
			Version:  p.Version,
		},
	}
}

func (p PackageId) String() string {
	return tools.Ternary(
		p.Platform == AllPlatform,
		"", string(p.Platform)+"/",
	) +
		string(p.Name) +
		tools.Ternary(
			p.Version == AllVersion,
			"",
			"@"+string(p.Version),
		)
}

func (p PackageId) StringFull() string {
	return p.Platform.String() + "/" + p.StringNameVersion()
}

func (p PackageId) StringNameVersion() string {
	return string(p.Name) + "@" + p.Version.String()
}

func (p PackageId) StringPlatformName() string {
	return string(p.Platform) + "/" + string(p.Name)
}
