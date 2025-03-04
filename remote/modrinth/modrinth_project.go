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

package modrinth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"lucy/datatypes"
	"lucy/lucytypes"
)

func getProjectId(slug lucytypes.PackageName) (id string) {
	res, _ := http.Get(projectUrl(string(slug)))
	modrinthProject := datatypes.ModrinthProject{}
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &modrinthProject)
	id = modrinthProject.Id
	return
}

func getProjectById(id string) (project *datatypes.ModrinthProject) {
	res, _ := http.Get(projectUrl(id))
	data, _ := io.ReadAll(res.Body)
	project = &datatypes.ModrinthProject{}
	json.Unmarshal(data, project)
	return
}

func getProjectByName(slug lucytypes.PackageName) (project *datatypes.ModrinthProject) {
	res, _ := http.Get(projectUrl(string(slug)))
	data, _ := io.ReadAll(res.Body)
	project = &datatypes.ModrinthProject{}
	json.Unmarshal(data, project)
	return
}

func getProjectMembers(id string) (members []*datatypes.ModrinthMember) {
	res, _ := http.Get(projectMemberUrl(id))
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &members)
	return
}

var ErrorInvalidDependency = errors.New("invalid dependency")

func DependencyToPackage(
	depedent lucytypes.PackageId,
	dependency *datatypes.ModrinthVersionDependencies,
) (
	p lucytypes.PackageId,
	err error,
) {
	var version *datatypes.ModrinthVersion
	var project *datatypes.ModrinthProject

	// I don't see a case where a package would depend on a project on another
	// platform. So, we can safely assume that the platform of the dependent
	// package is the same as the platform of the dependency.
	p.Platform = depedent.Platform

	if dependency.VersionId != "" && dependency.ProjectId != "" {
		version = getVersionById(dependency.VersionId)
		project = getProjectById(dependency.ProjectId)
	} else if dependency.VersionId != "" {
		version = getVersionById(dependency.VersionId)
		project = getProjectById(version.ProjectId)
	} else if dependency.ProjectId != "" {
		project = getProjectById(dependency.ProjectId)
		// This is not safe, TODO: use better inference method
		version = latestVersion(lucytypes.PackageName(project.Slug))
		p.Version = lucytypes.LatestVersion
	} else {
		return p, ErrorInvalidDependency
	}

	p.Name = lucytypes.PackageName(project.Slug)
	p.Version = lucytypes.RawVersion(version.VersionNumber)

	return p, nil
}
