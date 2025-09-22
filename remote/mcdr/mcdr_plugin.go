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

package mcdr

import (
	"lucy/remote"
)

func getPlugin(id string) (*plugin, error) {
	everything, err := getEverything()
	if err != nil {
		return nil, err
	}
	p, ok := everything.Plugins[id]
	if !ok {
		return nil, remote.FormatRemoteError(remote.ErrorNoPackage, id)
	}
	return &p, nil
}

func getAuthor(name string) *author {
	everything, err := getEverything()
	if err != nil {
		return nil
	}
	a, ok := everything.Authors.Authors[name]
	if !ok {
		return nil
	}
	return &a
}

func getRelease(plugin *plugin, version string) (release *release, err error) {
	if plugin == nil {
		return nil, remote.FormatRemoteError(remote.ErrorNoPackage)
	}

	// treat "" (empty), "any", and "latest" as the latest version
	if version == "" || version == "any" || version == "latest" {
		if len(plugin.Release.Releases) == 0 {
			return nil, remote.FormatRemoteError(
				remote.ErrorNoVersion,
				plugin.Meta.Name,
				version,
			)
		}
		return &plugin.Release.Releases[0], nil
	}

	// search for the specific version
	for _, release := range plugin.Release.Releases {
		if release.Meta.Version == version {
			return &release, nil
		}
	}
	return nil, remote.FormatRemoteError(
		remote.ErrorNoVersion,
		plugin.Meta.Name,
		version,
	)
}
