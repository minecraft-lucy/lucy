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

package mojang

import (
	"encoding/json"
	"io"
	"net/http"

	"lucy/datatype"
)

const VersionManifestURL = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

func getVersionManifest() (manifest *datatype.VersionManifest, err error) {
	manifest = &datatype.VersionManifest{}

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
