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

// Package remote is an adapter to its nested packages, which are responsible for
// fetching, searching, and providing information about packages from different
// sources. They are eventually unified into a single interface lucytype.Package.
//
// lucytype.Package itself utilizes a composite pattern, where its most fields,
// except the id, are optional and will be filled in as needed.
package remote

import (
	"fmt"

	"lucy/lucytype"
)

// IoC via dependency injection

func Fetch(
	source SourceHandler,
	id lucytype.PackageId,
) (remote lucytype.PackageRemote, err error) {
	raw, err := source.Fetch(id)
	if err != nil {
		return lucytype.PackageRemote{}, err
	}
	remote = raw.ToPackageRemote()
	return remote, nil
}

func Dependencies(
	source SourceHandler,
	id lucytype.PackageId,
) (deps *lucytype.PackageDependencies, err error) {
	// TODO: Implement
	return nil, fmt.Errorf("%w: %s", ErrorSourceNotSupported, source)
}

func Support(source lucytype.Source, name lucytype.ProjectName) (
	supports *lucytype.ProjectSupport,
	err error,
) {
	// TODO: Implement
	panic("not implemented")
}

func Information(
	source SourceHandler,
	name lucytype.ProjectName,
) (info lucytype.ProjectInformation, err error) {
	raw, err := source.Information(name)
	if err != nil {
		return lucytype.ProjectInformation{}, err
	}
	info = raw.ToProjectInformation()
	return info, nil
}

func Search(
	source SourceHandler,
	query lucytype.ProjectName,
	option lucytype.SearchOptions,
) (res lucytype.SearchResults, err error) {
	raw, err := source.Search(string(query), option)
	if err != nil {
		return res, err
	}
	res = raw.ToSearchResults()
	if len(res.Results) == 0 {
		return res, ErrorNoResults
	}
	return res, nil
}

// InferVersion replaces inferable version constants with their inferred versions
// through sources. You should call this function before parsing the version to
// ComparableVersion.
//
// TODO: Remove, infer version should not be exposed. All inference will be done in the SourceHandlers
func InferVersion(
	source SourceHandler,
	id lucytype.PackageId,
) (infer lucytype.PackageId) {
	return id
}
