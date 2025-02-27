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

package datatypes

import "lucy/lucytypes"

// ForgeModIdentifierNew is for 1.13+ forge & neoforge. This is a toml file.
type ForgeModIdentifierNew struct {
	ModLoader       string `toml:"modLoader"`
	LoaderVersion   string `toml:"loaderVersion"`
	IssueTrackerURL string `toml:"issueTrackerURL"`
	LogoFile        string `toml:"logoFile"`
	License         string `toml:"license"`
	Mods            []struct {
		ModID         lucytypes.PackageName    `toml:"modId"`
		Version       lucytypes.PackageVersion `toml:"version"`
		DisplayName   string                   `toml:"displayName"`
		ItemIcon      string                   `toml:"itemIcon"`
		DisplayURL    string                   `toml:"displayURL"`
		UpdateJSONURL string                   `toml:"updateJSONURL"`
		LogoFile      string                   `toml:"logoFile"`
		Credits       string                   `toml:"credits"`
		Authors       string                   `toml:"authors"`
		Description   string                   `toml:"description"`
	} `toml:"mods"`
	Dependencies  map[string][]ForgeDependencies `toml:"dependencies"`
	ModProperties map[string]string              `toml:"-"` // ignored
}

type ForgeDependencies struct {
	ModID        lucytypes.PackageName `toml:"modId"`
	Mandatory    bool                  `toml:"mandatory"`
	VersionRange string                `toml:"versionRange"`
	Ordering     string                `toml:"ordering"`
	Side         string                `toml:"side"`
}
