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

// Package modrinth provides functions to interact with Modrinth API
//
// We here use Modrinth terms in private functions:
//   - project: A project is a mod, plugin, or resource pack.
//   - Version: A version is a release, beta, or alpha version of a project.
//
// Generally, a project in Modrinth is equivalent to a project in Lucy. And
// a version in Modrinth is equivalent to a package in Lucy.
//
// Here, while referring to a project in lucy, we would try to the term "slug"
// to refer to the project (or it's name).
package modrinth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"lucy/lucytype"
	"lucy/remote"

	"lucy/logger"
	"lucy/tools"
)

type self struct{}

func (s self) Name() lucytype.Source {
	return lucytype.Modrinth
}

var Self self

// Search
//
// For Modrinth search API, see:
// https://docs.modrinth.com/api/operations/searchprojects/
func (s self) Search(
	query string,
	options lucytype.SearchOptions,
) (res remote.RawSearchResults, err error) {
	var facets []facetItems
	switch options.Platform {
	case lucytype.Forge:
		facets = append(facets, facetForge)
	case lucytype.Fabric:
		facets = append(facets, facetFabric)
	case lucytype.AllPlatform:
		fallthrough
	default:
		facets = append(facets, facetForge, facetAllLoaders)
	}

	if options.ShowClientPackage {
		facets = append(facets, facetServerSupported, facetClientSupported)
	} else {
		facets = append(facets, facetServerSupported)
	}

	internalOptions := searchOptions{
		index:  options.IndexBy.ToModrinth(),
		facets: facets,
	}
	searchUrl := searchUrl(lucytype.ProjectName(query), internalOptions)

	// Make the call to Modrinth API
	logger.Debug("searching via modrinth api: " + searchUrl)
	httpRes, err := http.Get(searchUrl)
	if err != nil {
		return nil, ErrInvalidAPIResponse
	}
	data, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	defer tools.CloseReader(httpRes.Body, logger.Warn)
	res = &searchResultResponse{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s self) Fetch(id lucytype.PackageId) (
	remote remote.RawPackageRemote,
	err error,
) {
	id, err = s.ParseAmbiguousVersion(id)
	version, err := getVersion(id)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (s self) Information(name lucytype.ProjectName) (
	info remote.RawProjectInformation,
	err error,
) {
	project, err := getProjectByName(name)
	if err != nil {
		return nil, err
	}
	return project, nil
}

// Support from Modrinth API is extremely unreliable. A local check (if any
// files were downloaded) is recommended.
func (s self) Support(name lucytype.ProjectName) (
	supports remote.RawProjectSupport,
	err error,
) {
	project, err := getProjectByName(name)
	if err != nil {
		return nil, err
	}
	return project, nil
}

var ErrInvalidAPIResponse = errors.New("invalid data from modrinth api")

func (s self) Dependencies(id lucytype.PackageId) (
	deps remote.RawPackageDependencies,
	err error,
) {
	// TODO implement me
	panic("implement me")
}

func (s self) ParseAmbiguousVersion(p lucytype.PackageId) (
	parsed lucytype.PackageId,
	err error,
) {
	parsed.Platform = p.Platform
	parsed.Name = p.Name
	var v *versionResponse

	switch p.Version {
	case lucytype.LatestCompatibleVersion:
		v, err = LatestCompatibleVersion(p.Name)
	case lucytype.AllVersion, lucytype.NoVersion, lucytype.LatestVersion:
		v, err = latestVersion(p.Name)
	default:
		return p, nil
	}
	if err != nil {
		return p, err
	}
	parsed.Version = lucytype.RawVersion(v.VersionNumber)

	return parsed, nil
}
