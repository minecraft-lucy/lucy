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

import "lucy/types"

type SourceHandler interface {
	Search(query string, options types.SearchOptions) (
		res RawSearchResults,
		err error,
	)
	Fetch(id types.PackageId) (
		remote RawPackageRemote,
		err error,
	)
	Information(name types.ProjectName) (
		info RawProjectInformation,
		err error,
	)
	Dependencies(id types.PackageId) (
		deps RawPackageDependencies,
		err error,
	)
	Support(name types.ProjectName) (
		supports RawProjectSupport,
		err error,
	)
	ParseAmbiguousVersion(id types.PackageId) (
		parsed types.PackageId,
		err error,
	)
	Name() types.Source
}

// Raw interfaces are designed for lazy evaluation and conversion to typed
// structures only when necessary. More functionality can be added to these
// interfaces as needed.

type (
	RawProjectSupport interface {
		ToProjectSupport() types.PlatformSupport
	}
	RawProjectInformation interface {
		ToProjectInformation() types.ProjectInformation
	}
	RawPackageRemote interface {
		ToPackageRemote() types.PackageRemote
	}
	RawPackageDependencies interface {
		ToPackageDependencies() types.PackageDependencies
	}

	// TODO: Consider make SortBy a method on the RawSearchResults interface

	RawSearchResults interface {
		ToSearchResults() types.SearchResults
	}
)
