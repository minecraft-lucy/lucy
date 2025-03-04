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

package syntax

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/mod/semver"

	"lucy/datatypes"
	"lucy/lucytypes"
)

var (
	EInvalidVersionComparison = errors.New("invalid version comparison")
	EVersionNotFound          = errors.New("version do not exist")
)

// TODO: Remove the err return value
// TODO: Use tools.Memoize to cache the result

const VersionManifestURL = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

func getVersionManifest() (manifest *datatypes.VersionManifest, err error) {
	manifest = &datatypes.VersionManifest{}

	// TODO: Add cache mechanism if http call is too slow or fails
	resp, err := http.Get(VersionManifestURL)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

// ComparePackageVersions gives -1 when v1 is older than v2, 0 when they are
// the same (or an error occurred), and 1 when v1 is newer than v2. 0 is returned
// when either v1 or v2 is AllVersion
func ComparePackageVersions(p1, p2 *lucytypes.PackageId) (c int8, err error) {
	v1, v2 := p1.Version, p2.Version

	if v1 == lucytypes.AllVersion || v2 == lucytypes.AllVersion {
		return 0, nil
	}

	if p1.Platform != p2.Platform {
		return 0, fmt.Errorf(
			"%w between package on platform %s and %s",
			EInvalidVersionComparison,
			p1.Platform,
			p2.Platform,
		)
	}

	if v1 == v2 {
		return 0, nil
	}

	if p1.Platform == lucytypes.Minecraft {
		return compareMinecraftVersions(v1, v2)
	}
	return int8(semver.Compare("v"+string(v1), "v"+string(v2))), nil
}

func compareMinecraftVersions(v1, v2 lucytypes.RawVersion) (
	c int8,
	err error,
) {
	manifest, err := getVersionManifest()
	if err != nil {
		return 0, err
	}

	i1, i2 := -1, -1
	for i, v := range manifest.Versions {
		if v1 == lucytypes.RawVersion(v.Id) {
			i1 = i
		}
		if v2 == lucytypes.RawVersion(v.Id) {
			i2 = i
		}
	}

	if i1 == i2 {
		return 0, nil
	}
	if i1 < i2 {
		return 1, nil
	}
	if i1 > i2 {
		return -1, nil
	}

	return 0, EInvalidVersionComparison
}
