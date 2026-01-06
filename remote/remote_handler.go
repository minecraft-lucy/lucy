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

import "lucy/lucytype"

type SourceHandler interface {
	Search(query string, options lucytype.SearchOptions) (
		res RawSearchResults,
		err error,
	)
	Fetch(id lucytype.PackageId) (
		remote RawPackageRemote,
		err error,
	)
	Information(name lucytype.ProjectName) (
		info RawProjectInformation,
		err error,
	)
	Dependencies(id lucytype.PackageId) (
		deps RawPackageDependencies,
		err error,
	)
	Support(name lucytype.ProjectName) (
		supports RawProjectSupport,
		err error,
	)
	ParseAmbiguousVersion(id lucytype.PackageId) (
		parsed lucytype.PackageId,
		err error,
	)
	Name() lucytype.Source
}

type (
	RawProjectSupport interface {
		ToProjectSupport() lucytype.ProjectSupport
	}
	RawProjectInformation interface {
		ToProjectInformation() lucytype.ProjectInformation
	}
	RawPackageRemote interface {
		ToPackageRemote() lucytype.PackageRemote
	}
	RawPackageDependencies interface {
		ToPackageDependencies() lucytype.PackageDependencies
	}

	// TODO: Consider make SortBy a method on the RawSearchResults interface

	RawSearchResults interface {
		ToSearchResults() lucytype.SearchResults
	}
)
