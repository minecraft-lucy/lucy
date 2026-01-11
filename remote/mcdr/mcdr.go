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
	"lucy/types"
)

type self struct{}

func (s self) Name() types.Source {
	return types.McdrCatalogue
}

var Self self

// Just a trivial type to implement the SearchResults interface
type mcdrSearchResult []string

func (m mcdrSearchResult) ToSearchResults() types.SearchResults {
	var res types.SearchResults
	for _, id := range m {
		res.Results = append(res.Results, types.ProjectName(id))
	}
	res.Source = types.McdrCatalogue
	return res
}

// TODO: handle search options

func (s self) Search(
query string,
options types.SearchOptions,
) (res remote.RawSearchResults, err error) {
	res, err = searchPlugin(query)
	return
}

func (s self) Fetch(id types.PackageId) (
rem remote.RawPackageRemote,
err error,
) {
	return getRelease(id.Name.Pep8String(), id.Version)
}

func (s self) Information(name types.ProjectName) (
info remote.RawProjectInformation,
err error,
) {
	plugin, err := getPluginInfo(name.Pep8String())
	if err != nil {
		return nil, err
	}
	meta, err := getMeta(name.Pep8String())
	if err != nil {
		return nil, err
	}
	repo, err := getRepositoryInfo(name.Pep8String())
	if err != nil {
		return nil, err
	}

	info = rawProjectInformation{
		Info:       plugin,
		Meta:       meta,
		Repository: repo,
	}

	return info, nil
}

func (s self) Dependencies(id types.PackageId) (
remote.RawPackageDependencies,
error,
) {
	// TODO implement me
	panic("implement me")
}

func (s self) Support(name types.ProjectName) (
supports remote.RawProjectSupport,
err error,
) {
	// TODO implement me
	panic("implement me")
}

func (s self) ParseAmbiguousVersion(id types.PackageId) (
parsed types.PackageId,
err error,
) {
	var rel *release
	switch id.Version {
	case types.LatestVersion:
		rel, err = getRelease(id.Name.Pep8String(), id.Version)
		if err != nil {
			return id, err
		}
	}
	parsed = types.PackageId{
		Platform: types.Mcdr,
		Name:     id.Name,
		Version:  types.RawVersion(rel.Meta.Version),
	}
	return parsed, nil
}
