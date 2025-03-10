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

// Package syntax defines the syntax for specifying packages and platforms.
//
// A package can either be specified by a string in the format of
// "platform/name@version". Only the name is required, both platform and version
// can be omitted.
//
// Valid Examples:
//   - carpet
//   - mcdr/prime-backup
//   - fabric/jade@1.0.0
//   - fabric@12.0
//   - minecraft@1.19 (recommended)
//   - minecraft/minecraft@1.16.5 (= minecraft@1.16.5)
//   - 1.8.9 (= minecraft@1.8.9)
package syntax

import (
	"errors"
	"strings"

	"lucy/dependency"

	"lucy/lnout"
	"lucy/lucytypes"
)

func PackageName(s string) lucytypes.ProjectName {
	return lucytypes.ProjectName(sanitize(s))
}

// sanitize tolerates some common interchangeability between characters. This
// includes underscores, chinese full stops, and backslashes. It also converts
// uppercase characters to lowercase.
func sanitize(s string) (clean string) {
	clean = ""
	for _, char := range s {
		if char == '_' {
			clean += string('-')
		} else if char == '\\' {
			clean += string('/')
		} else if char == '。' {
			clean += string('.')
		} else if char >= 'A' && char <= 'Z' {
			clean += strings.ToLower(string(char))
		} else {
			clean += string(char)
		}
	}
	return
}

var (
	ESyntax   = errors.New("invalid syntax")
	EPlatform = errors.New("invalid platform")
)

// Parse is exported to parse a string into a PackageId struct. This function
// should only be used on user inputs. Therefore, It does NOT return an
// error but instead invokes a fatal if the input is invalid.
func Parse(s string) (p lucytypes.PackageId) {
	s = sanitize(s)
	p = lucytypes.PackageId{}
	var err error
	p.Platform, p.Name, p.Version, err = parseOperatorAt(s)
	if err != nil {
		if errors.Is(err, ESyntax) {
			panic(err)
		} else {
			lnout.Fatal(err)
		}
	}
	lnout.Debug("parsed input as package: " + p.StringFull())
	return
}

// parseOperatorAt is called first since '@' operator always occur after '/' (equivalent
// to a lower priority).
func parseOperatorAt(s string) (
pl lucytypes.Platform,
n lucytypes.ProjectName,
v dependency.RawVersion,
err error,
) {
	split := strings.Split(s, "@")

	pl, n, err = parseOperatorSlash(split[0])
	if err != nil {
		return "", "", "", ESyntax
	}

	if len(split) == 1 {
		v = dependency.AllVersion
	} else if len(split) == 2 {
		v = dependency.RawVersion(split[1])
		if v == dependency.NoVersion || v == dependency.AllVersion {
			return "", "", "", ESyntax
		}
	} else {
		return "", "", "", ESyntax
	}

	return
}

func parseOperatorSlash(s string) (
pl lucytypes.Platform,
n lucytypes.ProjectName,
err error,
) {
	split := strings.Split(s, "/")

	if len(split) == 1 {
		pl = lucytypes.AllPlatform
		n = lucytypes.ProjectName(split[0])
		if lucytypes.Platform(n).Valid() {
			// Remember, all platforms are also valid packages under themselves.
			// This literal is for users to specify the platform itself. See the
			// docs for syntaxtypes.Platform for more information.
			pl = lucytypes.Platform(n)
			n = lucytypes.ProjectName(pl)
		}
	} else if len(split) == 2 {
		pl = lucytypes.Platform(split[0])
		if !pl.Valid() {
			return "", "", EPlatform
		}
		n = lucytypes.ProjectName(split[1])
	} else {
		return "", "", ESyntax
	}

	return
}
