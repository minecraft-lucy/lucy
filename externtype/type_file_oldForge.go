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

package externtype

// ForgeModIdentifierOld is for 1.12 and older forge mods. This is a json file.
type ForgeModIdentifierOld []struct {
	ModId        string        `json:"modid"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Version      string        `json:"version"`
	McVersion    string        `json:"mcversion"`
	URL          string        `json:"url"`
	UpdateURL    string        `json:"updateUrl"`
	AuthorList   []string      `json:"authorList"`
	Credits      string        `json:"credits"`
	LogoFile     string        `json:"logoFile"`
	Screenshots  []interface{} `json:"screenshots"`
	Dependencies []interface{} `json:"dependencies"`
}
