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
	"lucy/lucytypes"
	"lucy/remote"
)

type self struct{}

var Mcdr self

func (s self) Search(
	query string,
	options lucytypes.SearchOptions,
) (res lucytypes.SearchResults, err error) {
	res = lucytypes.SearchResults{}
	res.Source = lucytypes.McdrCatalogue
	everything, err := getEverything()
	if err != nil {
		return res, err
	}
	match(everything, query)
	res.Results = sortBy(everything, options.IndexBy)
	return res, nil
}

func (s self) Fetch(id lucytypes.PackageId) (
	remote remote.RawPackageRemote,
	err error,
) {
	// TODO implement me
	panic("implement me")
}

func (s self) Information(name lucytypes.ProjectName) (
	info remote.RawProjectInformation,
	err error,
) {
	// TODO implement me
	panic("implement me")
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
	// TODO implement me
	panic("implement me")
}
