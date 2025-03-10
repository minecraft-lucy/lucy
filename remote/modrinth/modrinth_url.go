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
	"net/url"
	"strings"
	"text/template"

	"lucy/lnout"
	"lucy/lucytypes"
)

const projectUrlPrefix = "https://api.modrinth.com/v2/project/"

func versionsUrl(slug lucytypes.ProjectName) (urlString string) {
	urlString, _ = url.JoinPath(
		projectUrlPrefix,
		string(slug),
		"version",
	)
	return
}

const versionUrlPrefix = `https://api.modrinth.com/v2/version/`

func versionUrl(id string) (urlString string) {
	return versionUrlPrefix + id
}

// projectUrl returns the URL for a project with the given Modrinth project id
// or slug (package name).
func projectUrl(suffix string) (urlString string) {
	return projectUrlPrefix + string(suffix)
}

func projectMemberUrl(suffix string) (urlString string) {
	return projectUrl(suffix) + "/members"
}

func projectDependencyUrl(suffix string) (urlString string) {
	return projectUrl(suffix) + "/dependencies"
}

const searchUrlTemplate = `https://api.modrinth.com/v2/search?query={{.query}}&limit=100&index={{.index}}&facets={{.facets}}`

func searchUrl(
	query lucytypes.ProjectName,
	option searchOptions,
) (urlString string) {
	urlTemplate, _ := template.New("modrinth_search_url").Parse(searchUrlTemplate)
	urlBuilder := strings.Builder{}
	err := urlTemplate.Execute(
		&urlBuilder,
		map[string]any{
			"query":  query,
			"index":  option.index,
			"facets": url.QueryEscape(serializeFacet(option.facets...)),
		},
	)
	if err != nil {
		lnout.Error(err)
	}

	urlString = urlBuilder.String()
	return urlString
}

const userHomepageUrlPrefix = `https://modrinth.com/user/`

// userHomepageUrl's suffix is the user's username or id.
func userHomepageUrl(suffix string) (urlString string) {
	return userHomepageUrlPrefix + suffix
}
