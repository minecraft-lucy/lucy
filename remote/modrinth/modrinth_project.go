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

	"lucy/syntax"

	"lucy/lucytypes"
)

func getProjectId(slug lucytypes.ProjectName) (id string, err error) {
	res, _ := http.Get(projectUrl(string(slug)))
	modrinthProject := projectResponse{}
	data, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(data, &modrinthProject)
	if err != nil {
		return "", ENoProject
	}
	id = modrinthProject.Id
	return
}

func getProjectById(id string) (project *projectResponse, err error) {
	res, _ := http.Get(projectUrl(id))
	data, _ := io.ReadAll(res.Body)
	project = &projectResponse{}
	err = json.Unmarshal(data, project)
	if err != nil {
		return nil, ENoProject
	}
	return
}

func getProjectByName(slug lucytypes.ProjectName) (
	project *projectResponse,
	err error,
) {
	res, _ := http.Get(projectUrl(string(slug)))
	data, _ := io.ReadAll(res.Body)
	project = &projectResponse{}
	err = json.Unmarshal(data, project)
	if err != nil {
		return nil, ENoProject
	}
	return
}

func getProjectMembers(id string) (
	members []*memberResponse,
	err error,
) {
	res, _ := http.Get(projectMemberUrl(id))
	data, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(data, &members)
	if err != nil {
		return nil, ENoMember
	}
	return members, nil
}

var ErrorInvalidDependency = errors.New("invalid dependency")

func DependencyToPackage(
	dependent lucytypes.PackageId,
	dependency *dependenciesResponse,
) (
	p lucytypes.PackageId,
	err error,
) {
	var version *versionResponse
	var project *projectResponse

	// I don't see a case where a package would depend on a project on another
	// platform. So, we can safely assume that the platform of the dependent
	// package is the same as the platform of the dependency.
	p.Platform = dependent.Platform

	if dependency.VersionId != "" && dependency.ProjectId != "" {
		version, _ = getVersionById(dependency.VersionId)
		project, _ = getProjectById(dependency.ProjectId)
	} else if dependency.VersionId != "" {
		version, _ = getVersionById(dependency.VersionId)
		project, _ = getProjectById(version.ProjectId)
	} else if dependency.ProjectId != "" {
		project, _ = getProjectById(dependency.ProjectId)
		// This is not safe, TODO: use better inference method
		version, _ = latestVersion(lucytypes.ProjectName(project.Slug))
		p.Version = lucytypes.LatestVersion
	} else {
		return p, ErrorInvalidDependency
	}

	p.Name = syntax.PackageName(project.Slug)
	p.Version = lucytypes.RawVersion(version.VersionNumber)

	return p, nil
}
