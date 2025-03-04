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
//   - Project: A project is a mod, plugin, or resource pack.
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
	"strconv"

	"lucy/datatypes"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/tools"
)

var ErrorInvalidAPIResponse = errors.New("invalid data from modrinth api")

// Search
//
// For Modrinth search API, see:
// https://docs.modrinth.com/api/operations/searchprojects/
func Search(
	packageId lucytypes.PackageId,
	options lucytypes.SearchOptions,
) (result *lucytypes.SearchResults, err error) {
	var facets []facetItems
	query := packageId.Name

	switch packageId.Platform {
	case lucytypes.Forge:
		facets = append(facets, facetForge)
	case lucytypes.Fabric:
		facets = append(facets, facetFabric)
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
	searchUrl := searchUrl(
		query,
		internalOptions,
	)

	// Make the call to Modrinth API
	logger.Debug("searching via modrinth api: " + searchUrl)
	resp, err := http.Get(searchUrl)
	if err != nil {
		return nil, ErrorInvalidAPIResponse
	}
	data, err := io.ReadAll(resp.Body)
	defer tools.CloseReader(resp.Body, logger.Warning)
	var searchResults datatypes.ModrinthSearchResults
	err = json.Unmarshal(data, &searchResults)
	if err != nil {
		return nil, err
	}
	if searchResults.Hits == nil {
		return nil, nil
	}
	if searchResults.TotalHits > 100 {
		logger.Info(strconv.Itoa(searchResults.TotalHits) + " results found on modrinth, only showing first 100")
	}

	result = &lucytypes.SearchResults{}
	result.Results = make([]string, 0, len(searchResults.Hits))
	result.Source = lucytypes.Modrinth
	for _, hit := range searchResults.Hits {
		result.Results = append(result.Results, hit.Slug)
	}
	return result, nil
}

func Fetch(id lucytypes.PackageId) (
	remote *lucytypes.PackageRemote,
	err error,
) {
	id = inferVersion(id)
	project := getProjectByName(id.Name)
	version, err := getVersion(id)
	if err != nil {
		logger.Fatal(err)
	}
	fileUrl, filename := getFile(version)

	remote = &lucytypes.PackageRemote{
		Source:   lucytypes.Modrinth,
		RemoteId: project.Id,
		FileUrl:  fileUrl,
		Filename: filename,
	}

	return remote, nil
}

func Information(slug lucytypes.PackageName) (
	information *lucytypes.PackageInformation,
	err error,
) {
	project := getProjectByName(slug)
	information = &lucytypes.PackageInformation{
		Name:        project.Title,
		Brief:       project.Description,
		Description: tools.MarkdownToPlainText(project.Body),
		Author:      []lucytypes.PackageMember{},
		Urls:        []lucytypes.PackageUrl{},
		License:     project.License.Name,
	}

	// Fill in URLs
	if project.WikiUrl != "" {
		information.Urls = append(
			information.Urls,
			lucytypes.PackageUrl{
				Name: "Wiki",
				Type: lucytypes.WikiUrl,
				Url:  project.WikiUrl,
			},
		)
	}

	if project.SourceUrl != "" {
		information.Urls = append(
			information.Urls,
			lucytypes.PackageUrl{
				Name: "Source Code",
				Type: lucytypes.SourceUrl,
				Url:  project.SourceUrl,
			},
		)
	}

	if project.DonationUrls != nil {
		for _, donationUrl := range project.DonationUrls {
			information.Urls = append(
				information.Urls,
				lucytypes.PackageUrl{
					Name: "Donation",
					Type: lucytypes.OthersUrl,
					Url:  donationUrl.Url,
				},
			)
		}
	}

	// Fill in authors
	members := getProjectMembers(project.Id)
	for _, member := range members {
		information.Author = append(
			information.Author,
			lucytypes.PackageMember{
				Name:  member.User.Username,
				Role:  member.Role,
				Url:   userHomepageUrl(member.User.Id),
				Email: member.User.Email,
			},
		)
	}

	return information, nil
}

// Dependencies from Modrinth API is extremely unreliable. A local check (if any
// files were downloaded) is recommended.
func Dependencies(id lucytypes.PackageId) (dependencies *lucytypes.PackageDependencies) {
	id = inferVersion(id)
	project := getProjectByName(id.Name)
	version, _ := getVersion(id)
	dependencies = &lucytypes.PackageDependencies{
		SupportedVersions:  []lucytypes.RawVersion{},
		SupportedPlatforms: []lucytypes.Platform{},
		Required:           []lucytypes.PackageId{},
	}

	for _, version := range project.GameVersions {
		dependencies.SupportedVersions = append(
			dependencies.SupportedVersions,
			lucytypes.RawVersion(version),
		)
	}

	for _, platform := range project.Loaders {
		dependencies.SupportedPlatforms = append(
			dependencies.SupportedPlatforms,
			lucytypes.Platform(platform),
		)
	}

	for _, dependency := range version.Dependencies {
		switch dependency.DependencyType {
		case datatypes.ModrinthVersionDependencyTypeIncompatible:
			d, err := DependencyToPackage(id, &dependency)
			if err != nil {
				logger.Warning(err)
				continue
			}
			dependencies.Incompatible = append(
				dependencies.Incompatible,
				d,
			)
		case datatypes.ModrinthVersionDependencyTypeOptional:
			d, err := DependencyToPackage(id, &dependency)
			if err != nil {
				logger.Warning(err)
				continue
			}
			dependencies.Optional = append(
				dependencies.Optional,
				d,
			)
		case datatypes.ModrinthVersionDependencyTypeRequired:
			d, err := DependencyToPackage(id, &dependency)
			if err != nil {
				logger.Warning(err)
				continue
			}
			dependencies.Required = append(
				dependencies.Required,
				d,
			)
		}
	}

	return dependencies
}

func GetProjectByName(packageName lucytypes.PackageName) (
	project *datatypes.ModrinthProject,
	err error,
) {
	res, err := http.Get(projectUrl(string(packageName)))
	if err != nil {
		return
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	project = &datatypes.ModrinthProject{}
	err = json.Unmarshal(data, project)
	return
}

func inferVersion(p lucytypes.PackageId) (infer lucytypes.PackageId) {
	infer.Platform = p.Platform
	infer.Name = p.Name

	switch p.Version {
	case lucytypes.AllVersion, lucytypes.NoVersion, lucytypes.LatestCompatibleVersion:
		version := LatestCompatibleVersion(p.Name)
		infer.Version = version.VersionNumber
	case lucytypes.LatestVersion:
		version := latestVersion(p.Name)
		infer.Version = version.VersionNumber
	default:
		return p
	}

	return infer
}
