package remote

import "lucy/lucytypes"

type SourceHandler interface {
	Fetch(id lucytypes.PackageId) (
	remote lucytypes.RawPackageRemote,
	err error,
	)

	Dependencies(id lucytypes.PackageId) (
	deps lucytypes.RawPackageDependencies,
	err error,
	)

	Support(name lucytypes.ProjectName) (
	supports lucytypes.RawProjectSupport,
	err error,
	)

	Information(name lucytypes.ProjectName) (
	info lucytypes.RawProjectInformation,
	err error,
	)

	Search(query string, option SearchOptions) (
	results []lucytypes.ProjectName,
	)

	ParseAmbiguousVersion(id lucytypes.PackageId) (
	parsed lucytypes.PackageId,
	err error,
	)
}
