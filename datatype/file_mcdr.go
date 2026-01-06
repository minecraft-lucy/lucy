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

import "lucy/tools"

type McdrPluginInfo struct {
	Id      string `json:"id"`
	Authors []struct {
		Name string `json:"name"`
		Link string `json:"link"`
	} `json:"authors"`
	Repository   string   `json:"repository"`
	Branch       string   `json:"branch"`
	RelatedPath  string   `json:"related_path"`
	Labels       []string `json:"labels"`
	Introduction struct {
		EnUs string `json:"en_us"`
		ZhCn string `json:"zh_cn"`
	} `json:"introduction"`
}

type McdrPluginIdentifierFile struct {
	Id          string `json:"id"`
	Version     string `json:"version"`
	Name        string `json:"name"`
	Description struct {
		EnUs string `json:"en_us"`
		ZhCn string `json:"zh_cn"`
	} `json:"description"`
	Author       tools.StringOrStringSlice `json:"author"`
	Link         string                    `json:"link"`
	Dependencies struct {
		Mcdreforged string `json:"mcdreforged"`
	} `json:"dependencies"`
	Resources []string `json:"resources"`
}
