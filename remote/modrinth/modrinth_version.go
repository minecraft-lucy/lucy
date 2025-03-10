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

package modrinth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"lucy/dependency"

	"lucy/lnout"

	"lucy/datatypes"
	"lucy/local"
	"lucy/lucytypes"
)

// TODO: Refactor to separate all API functions to accept an url. While the urls
// are generated by other functions. This will make the code more modular and
// easier to test.

var (
	ENoVersion = errors.New("modrinth version not found")
	ENoProject = errors.New("modrinth project not found")
	ENoMember  = errors.New("modrinth project member not found")
)

func listVersions(slug lucytypes.ProjectName) (
	versions []*datatypes.ModrinthVersion,
	err error,
) {
	res, _ := http.Get(versionsUrl(slug))
	data, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(data, &versions)
	if err != nil {
		return nil, ENoProject
	}
	return
}

// getVersion is named as so because a Package in lucy is equivalent to a version
// in Modrinth.
func getVersion(id lucytypes.PackageId) (
	v *datatypes.ModrinthVersion,
	err error,
) {
	versions, err := listVersions(id.Name)
	if err != nil {
		return nil, ENoVersion
	}
	if id.Version == dependency.LatestVersion {
		v, err = latestVersion(id.Name)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	for _, version := range versions {
		if version.VersionNumber == id.Version &&
			versionSupportsLoader(version, id.Platform) {
			return version, nil
		}
	}
	return nil, ENoVersion
}

func getVersionById(id string) (v *datatypes.ModrinthVersion, err error) {
	res, _ := http.Get(versionUrl(id))
	data, _ := io.ReadAll(res.Body)
	v = &datatypes.ModrinthVersion{}
	err = json.Unmarshal(data, v)
	if err != nil {
		return nil, ENoVersion
	}
	return
}

func versionSupportsLoader(
	version *datatypes.ModrinthVersion,
	loader lucytypes.Platform,
) bool {
	for _, l := range version.Loaders {
		if lucytypes.Platform(l).Eq(loader) {
			return true
		}
	}
	return false
}

func latestVersion(slug lucytypes.ProjectName) (
	v *datatypes.ModrinthVersion,
	err error,
) {
	versions, err := listVersions(slug)
	if err != nil {
		return nil, err
	}
	for _, version := range versions {
		if version.VersionType == "release" &&
			(v == nil || version.DatePublished.After(v.DatePublished)) {
			v = version
		}
	}
	if v == nil {
		lnout.Info("no release version found for " + slug.Title())
		return nil, ENoVersion
	} else {
		lnout.Debug("latest version of " + slug.String() + ": " + v.VersionNumber.String())
	}
	return v, nil
}

func LatestCompatibleVersion(slug lucytypes.ProjectName) (
	v *datatypes.ModrinthVersion,
	err error,
) {
	versions, err := listVersions(slug)
	if err != nil {
		return nil, err
	}
	serverInfo := local.GetServerInfo()
	if serverInfo.Executable == local.UnknownExecutable {
		lnout.Info("no executable found, unable to infer a compatible version. falling back to latest version")
		v, err := latestVersion(slug)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	for _, version := range versions {
		for _, gameVersion := range version.GameVersions {
			if gameVersion == serverInfo.Executable.GameVersion.String() &&
				version.VersionType == "release" &&
				(v == nil || version.DatePublished.After(v.DatePublished)) {
				v = version
			}
		}
	}
	return v, nil
}
