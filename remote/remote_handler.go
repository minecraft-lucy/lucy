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

package remote

import "lucy/lucytypes"

type SourceHandler interface {
	Search(query string, options lucytypes.SearchOptions) (
		res RawSearchResults,
		err error,
	)
	Fetch(id lucytypes.PackageId) (
		remote RawPackageRemote,
		err error,
	)
	Information(name lucytypes.ProjectName) (
		info RawProjectInformation,
		err error,
	)
	Dependencies(id lucytypes.PackageId) (
		deps RawPackageDependencies,
		err error,
	)
	Support(name lucytypes.ProjectName) (
		supports RawProjectSupport,
		err error,
	)
	ParseAmbiguousVersion(id lucytypes.PackageId) (
		parsed lucytypes.PackageId,
		err error,
	)
	Name() lucytypes.Source
}

type (
	RawProjectSupport interface {
		ToProjectSupport() lucytypes.ProjectSupport
	}
	RawProjectInformation interface {
		ToProjectInformation() lucytypes.ProjectInformation
	}
	RawPackageRemote interface {
		ToPackageRemote() lucytypes.PackageRemote
	}
	RawPackageDependencies interface {
		ToPackageDependencies() lucytypes.PackageDependencies
	}

	// TODO: Consider make SortBy a method on the RawSearchResults interface

	RawSearchResults interface {
		ToSearchResults() lucytypes.SearchResults
	}
)
