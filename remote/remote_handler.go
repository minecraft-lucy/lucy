package remote

import "lucy/lucytypes"

type SourceHandler interface {
	Search(query string, options lucytypes.SearchOptions) (
		res lucytypes.SearchResults,
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
