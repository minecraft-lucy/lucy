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

package exttype

type FabricEnvironment string

const (
	FabricEnvironmentClient FabricEnvironment = "client"
	FabricEnvironmentServer FabricEnvironment = "server"
	FabricEnvironmentAny    FabricEnvironment = "*"
)

// FileFabricModIdentifier represents the structure of fabric.mod.json files found
// in Fabric mods' `.jar` files.
//
// Docs: https://fabricmc.net/wiki/documentation:fabric_mod_json_spec
type FileFabricModIdentifier struct {
	SchemaVersion int      `json:"schemaVersion"`
	Id            string   `json:"id"`
	Version       string   `json:"version"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Authors       []string `json:"authors"`

	// Fields officially supported:
	//   - "email"
	//   - "homepage"
	//   - "irc"
	//   - "issues"
	//   - "sources"
	Contact map[string]string `json:"contact"`

	// This uses the SPDX format https://spdx.org/licenses/
	// TODO: Should implement and check whether other platforms use this too.
	License string `json:"license"`

	Icon        string            `json:"icon"`
	Environment FabricEnvironment `json:"environment"`
	Jars        []struct {
		File string `json:"file"`
	} `json:"-"`
	Entrypoints      map[string][]string `json:"-"`
	Mixins           []string            `json:"-"`
	AccessWidener    string              `json:"accessWidener"`
	LanguageAdapters map[string]string   `json:"-"`

	// Depends > Recommends > Suggests
	// Breaks > Conflicts
	Depends    map[string]string `json:"depends"`
	Recommends map[string]string `json:"recommends"`
	Suggests   map[string]string `json:"suggests"`
	Breaks     map[string]string `json:"breaks"`
	Conflicts  map[string]string `json:"conflicts"`

	Custom interface{} `json:"-"`
}

type FileFabricModIdentifierOld struct {
	// TODO: See https://wiki.fabricmc.net/documentation:fabric_mod_json_spec
	// This is for very old fabric (< 0.4.0). It does not matter much right
	// now. Besides, it is poorly documented.
	//
	// When SchemaVersion is 0 or missing, it is considered old.
}
