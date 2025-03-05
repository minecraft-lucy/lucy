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
// sources. They are eventually unified into a single interface lucytypes.Package.
//
// lucytypes.Package itself utilizes a composite pattern, where its most fields,
// except the id, are optional and will be filled in as needed.
package remote

import (
	"errors"
	"fmt"

	"lucy/dependency"
	"lucy/lucytypes"
	"lucy/remote/modrinth"
)

var (
	ESourceNotSupported  = errors.New("source not supported")
	ECannotInferPlatform = errors.New("cannot infer platform")
	ECannotInferSource   = errors.New("cannot infer source") // This will be reintroduced when the auto source detection is implemented.
)

func Fetch(
	source lucytypes.Source,
	id lucytypes.PackageId,
) (remote *lucytypes.PackageRemote, err error) {
	switch source {
	case lucytypes.Modrinth:
		fetch, err := modrinth.Fetch(id)
		if err != nil {
			return nil, err
		}
		return fetch, nil
	case lucytypes.CurseForge:
		fallthrough
	case lucytypes.McdrWebsite:
		fallthrough
	default:
		return nil, fmt.Errorf("%w: %s", ESourceNotSupported, source)
	}
}

func Dependencies(
	source lucytypes.Source,
	id lucytypes.PackageId,
) (deps *lucytypes.PackageDependencies, err error) {
	switch source {

	case lucytypes.Modrinth:
		fallthrough
	case lucytypes.CurseForge:
		fallthrough
	case lucytypes.McdrWebsite:
		fallthrough
	default:
		return nil, fmt.Errorf("%w: %s", ESourceNotSupported, source)
	}
}

func Information(
	source lucytypes.Source,
	name lucytypes.ProjectName,
) (info *lucytypes.ProjectInformation, err error) {
	switch source {
	case lucytypes.Modrinth:
		info, err = modrinth.Information(name)
		if err != nil {
			return nil, err
		}
		return info, nil
	case lucytypes.CurseForge:
		fallthrough
	case lucytypes.McdrWebsite:
		fallthrough
	default:
		return nil, fmt.Errorf("%w: %s", ESourceNotSupported, source)
	}
}

func Search(
	source lucytypes.Source,
	name lucytypes.ProjectName,
	option lucytypes.SearchOptions,
) (res *lucytypes.SearchResults, err error) {
	switch source {
	case lucytypes.Modrinth:
		return modrinth.Search(name, option)
	default:
		return nil, fmt.Errorf("%w: %s", ESourceNotSupported, source)
	}
}

// InferVersion replaces inferable version constants with their inferred versions
// through sources. You should call this function before parsing the version to
// SemanticVersion.
//
// TODO: Implement InferVersion for all RawVersion constants.
func InferVersion(
	source lucytypes.Source,
	id lucytypes.PackageId,
) (infer lucytypes.PackageId) {
	switch id.Version {
	case dependency.AllVersion, dependency.LatestVersion:
		// API call
	case dependency.LatestCompatibleVersion:
		// API call
	case dependency.NoVersion, dependency.UnknownVersion:
		// Do nothing
	default:
		// Do nothing
	}
	return id
}
