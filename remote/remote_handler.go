package remote

import "lucy/lucytypes"

type SourceHandler interface {
	Fetch(id lucytypes.PackageId) (
		remote RawPackageRemote,
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

	Information(name lucytypes.ProjectName) (
		info RawProjectInformation,
		err error,
	)

	Search(query string, option lucytypes.SearchOptions) (
		results []lucytypes.ProjectName,
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
	RawSearchResults interface {
		ToSearchResults() lucytypes.SearchResults
	}
)
