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
	"fmt"

	"lucy/lucytypes"
	"lucy/remote"
)

type self struct{}

func (s self) Name() lucytypes.Source {
	return lucytypes.McdrCatalogue
}

var Self self

func (s self) Search(
	query string,
	options lucytypes.SearchOptions,
) (res remote.RawSearchResults, err error) {
	if options.Platform != lucytypes.Mcdr && options.Platform != lucytypes.AllPlatform {
		return nil, remote.FormatError(
			remote.ErrorUnsupportedPlatform,
			lucytypes.McdrCatalogue,
			options.Platform,
		)
	}
	everything, err := getEverything()
	if err != nil {
		return nil, err
	}
	res = &queriedEverything{
		Everything: *everything,
		IndexBy:    options.IndexBy,
		Query:      query,
	}
	return res, nil
}

func (s self) Fetch(id lucytypes.PackageId) (
	rem remote.RawPackageRemote,
	err error,
) {
	if id.Platform != lucytypes.Mcdr && id.Platform != lucytypes.AllPlatform {
		return nil, remote.ErrorUnsupportedPlatform
	}
	p := getPlugin(string(id.Name))
	if p == nil {
		return nil, fmt.Errorf("%w: %s", remote.ErrorNoPackage, id)
	}

	version := id.Version
	if version.NeedsInfer() {
		parsed, err := s.ParseAmbiguousVersion(id)
		if err != nil {
			return nil, err
		}
		version = parsed.Version
	}
	release, err := getRelease(p, version.String())
	if err != nil {
		return nil, err
	}

	return release, nil
}

func (s self) Information(name lucytypes.ProjectName) (
	info remote.RawProjectInformation,
	err error,
) {
	p := getPlugin(string(name))
	if p == nil {
		return nil, fmt.Errorf("%w: %s", remote.ErrorNoPackage, name)
	}
	return p, nil
}

func (s self) Dependencies(id lucytypes.PackageId) (
	deps remote.RawPackageDependencies,
	err error,
) {
	// TODO implement me
	panic("implement me")
}

func (s self) Support(name lucytypes.ProjectName) (
	supports remote.RawProjectSupport,
	err error,
) {
	// TODO implement me
	panic("implement me")
}

func (s self) ParseAmbiguousVersion(id lucytypes.PackageId) (
	parsed lucytypes.PackageId,
	err error,
) {
	switch id.Version {
	case lucytypes.AllVersion, lucytypes.LatestVersion, lucytypes.LatestCompatibleVersion:
		p := getPlugin(id.Name.String())
		if p == nil {
			return id, remote.FormatError(
				remote.ErrorCannotInferVersion,
				id.Name,
				id.Version,
			)
		}
		id.Version = lucytypes.RawVersion(p.Release.LatestVersion)
		return id, nil
	}
	return id, remote.FormatError(
		remote.ErrorCannotInferVersion,
		id.Name,
		id.Version,
	)
}
