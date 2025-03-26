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
	"lucy/lnout"
	"lucy/lucytypes"
	"lucy/tools"
)

var ErrorInvalidAPIResponse = errors.New("invalid data from modrinth api")

// Search
//
// For Modrinth search API, see:
// https://docs.modrinth.com/api/operations/searchprojects/
func Search(
	name lucytypes.ProjectName,
	showClient bool,
	indexBy string,
	platform lucytypes.Platform,
) (result *datatypes.ModrinthSearchResults, err error) {
	result = &datatypes.ModrinthSearchResults{}
	var facets []facetItems
	switch platform {
	case lucytypes.Forge:
		facets = append(facets, facetForge)
	case lucytypes.Fabric:
		facets = append(facets, facetFabric)
	case lucytypes.AllPlatform:
		fallthrough
	default:
		facets = append(facets, facetForge, facetAllLoaders)
	}

	if showClient {
		facets = append(facets, facetServerSupported, facetClientSupported)
	} else {
		facets = append(facets, facetServerSupported)
	}

	internalOptions := searchOptions{
		index:  indexBy,
		facets: facets,
	}
	searchUrl := searchUrl(name, internalOptions)

	// Make the call to Modrinth API
	lnout.Debug("searching via modrinth api: " + searchUrl)
	resp, err := http.Get(searchUrl)
	if err != nil {
		return result, ErrorInvalidAPIResponse
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer tools.CloseReader(resp.Body, lnout.Warn)
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	if result.TotalHits > 100 {
		lnout.InfoNow(strconv.Itoa(result.TotalHits) + " results found on modrinth, only showing first 100")
	}
	return result, nil
}

func Fetch(id lucytypes.PackageId) (
	remote *lucytypes.PackageRemote,
	err error,
) {
	id = inferVersion(id)
	// project, err := getProjectByName(id.Name)
	// if err != nil {
	// 	return nil, err
	// }
	version, err := getVersion(id)
	if err != nil {
		return nil, err
	}
	fileUrl, filename := getFile(version)

	remote = &lucytypes.PackageRemote{
		FileUrl:  fileUrl,
		Filename: filename,
	}

	return remote, nil
}

func Information(slug lucytypes.ProjectName) (
	information *lucytypes.ProjectInformation,
	err error,
) {
	project, err := getProjectByName(slug)
	if err != nil {
		return nil, err
	}
	information = &lucytypes.ProjectInformation{
		Title:       project.Title,
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
	members, err := getProjectMembers(project.Id)
	if err != nil {
		lnout.WarnNow(err)
	} else {
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
	}

	return information, nil
}

// Support from Modrinth API is extremely unreliable. A local check (if any
// files were downloaded) is recommended.
func Support(name lucytypes.ProjectName) (
	supports *lucytypes.ProjectSupport,
	err error,
) {
	project, _ := getProjectByName(name)
	supports = &lucytypes.ProjectSupport{
		MinecraftVersions: make([]lucytypes.RawVersion, 0),
		Platforms:         make([]lucytypes.Platform, 0),
	}

	for _, version := range project.GameVersions {
		supports.MinecraftVersions = append(
			supports.MinecraftVersions,
			lucytypes.RawVersion(version),
		)
	}

	for _, platform := range project.Loaders {
		supports.Platforms = append(
			supports.Platforms,
			lucytypes.Platform(platform),
		)
	}

	return supports, nil
}

func inferVersion(p lucytypes.PackageId) (infer lucytypes.PackageId) {
	infer.Platform = p.Platform
	infer.Name = p.Name
	var v *datatypes.ModrinthVersion
	var err error

	switch p.Version {
	case lucytypes.LatestCompatibleVersion:
		v, err = LatestCompatibleVersion(p.Name)
	case lucytypes.AllVersion, lucytypes.NoVersion, lucytypes.LatestVersion:
		v, err = latestVersion(p.Name)
	default:
		return p
	}
	if err != nil {
		return p
	}
	infer.Version = lucytypes.RawVersion(v.VersionNumber)

	return infer
}
