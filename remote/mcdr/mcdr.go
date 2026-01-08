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
	"lucy/lucytype"
	"lucy/remote"
)

type self struct{}

func (s self) Name() lucytype.Source {
	return lucytype.McdrCatalogue
}

var Self self

// Just a trivial type to implement the SearchResults interface
type mcdrSearchResult []string

func (m mcdrSearchResult) ToSearchResults() lucytype.SearchResults {
	var res lucytype.SearchResults
	for _, id := range m {
		res.Results = append(res.Results, lucytype.ProjectName(id))
	}
	res.Source = lucytype.McdrCatalogue
	return res
}

// TODO: handle search options

func (s self) Search(
	query string,
	options lucytype.SearchOptions,
) (res remote.RawSearchResults, err error) {
	res, err = searchPlugin(query)
	return
}

func (s self) Fetch(id lucytype.PackageId) (
	rem remote.RawPackageRemote,
	err error,
) {
	// TODO implement me
	panic("implement me")
}

func (s self) Information(name lucytype.ProjectName) (
	info remote.RawProjectInformation,
	err error,
) {
	// TODO implement me
	panic("implement me")
}

func (s self) Dependencies(id lucytype.PackageId) (
	remote.RawPackageDependencies,
	error,
) {
	// TODO implement me
	panic("implement me")
}

func (s self) Support(name lucytype.ProjectName) (
	supports remote.RawProjectSupport,
	err error,
) {
	// TODO implement me
	panic("implement me")
}

func (s self) ParseAmbiguousVersion(id lucytype.PackageId) (
	parsed lucytype.PackageId,
	err error,
) {
	// TODO implement me
	panic("implement me")
}
