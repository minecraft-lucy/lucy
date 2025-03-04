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
	"fmt"
	"lucy/logger"
	"lucy/lucytypes"
)

func FetchSource(
	source lucytypes.Source,
	id lucytypes.PackageId,
) (remote *lucytypes.PackageRemote) {
	if source == lucytypes.Auto {
		source = SelectSource(id.Platform)
	}

	switch source {
	case lucytypes.Modrinth:
		// return modrinth.Fetch(id)
	default:
		logger.Fatal(fmt.Errorf("source fetch not supported yet:" + source.String()))
	}

	return nil // unreachable
}

func GetDependencies(
	source lucytypes.Source,
	id lucytypes.PackageId,
) *lucytypes.PackageDependencies {
	return nil
}

func GetInformation(
	source lucytypes.Source,
	id lucytypes.PackageId,
) *lucytypes.PackageInformation {
	return nil
}

func SearchForProject(
	source lucytypes.Source,
	query string,
) []lucytypes.PackageName {
	return nil
}

// InferVersion replaces inferable version constants with their inferred versions
// through sources. You should call this function before parsing the version to
// SemanticVersion.
//
// TODO: Implement InferVersion for all RawVersion constants.
func InferVersion(id lucytypes.PackageId) (infer lucytypes.PackageId) {
	switch id.Version {
	case lucytypes.AllVersion, lucytypes.LatestVersion:
		// API call
	case lucytypes.LatestCompatibleVersion:
		// API call
	case lucytypes.NoVersion, lucytypes.UnknownVersion:
		// Do nothing
	default:
		// Do nothing
	}
	return id
}
