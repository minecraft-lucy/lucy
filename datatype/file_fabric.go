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

package datatype

type FabricModIdentifier struct {
	SchemaVersion int      `json:"schemaVersion"`
	Id            string   `json:"id"`
	Version       string   `json:"version"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Authors       []string `json:"authors"`
	Contact       struct {
		Homepage string `json:"homepage"`
		Issues   string `json:"issues"`
		Sources  string `json:"sources"`
	} `json:"contact"`
	License     string `json:"license"`
	Icon        string `json:"icon"`
	Environment string `json:"environment"`
	Entrypoints struct {
		Client []string `json:"client"`
		Server []string `json:"server"`
	} `json:"entrypoints"`
	Mixins        []string `json:"mixins"`
	AccessWidener string   `json:"accessWidener"`
	Depends       struct {
		Minecraft    string `json:"minecraft"`
		Fabricloader string `json:"fabricloader"`
		Java         string `json:"java"`
	} `json:"depends"`
}
