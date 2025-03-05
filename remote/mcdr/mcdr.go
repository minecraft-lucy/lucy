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

package mcdr

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/google/go-github/v50/github"
	"lucy/datatypes"
	"lucy/lucytypes"
	"lucy/syntax"
)

func SearchMcdrPluginCatalogue(search lucytypes.ProjectName) (
	pluginInfo *datatypes.McdrPluginInfo,
	err error,
) {
	plugins := getMcdrPluginCatalogue()

	for _, plugin := range plugins {
		p := syntax.Parse(*plugin.Name)
		if p.Name == search {
			return getMcdrPluginInfo(*plugin.Path), nil
		}
	}

	return nil, fmt.Errorf("plugin not found")
}

func getMcdrPluginCatalogue() []*github.RepositoryContent {
	ctx := context.Background()
	client := github.NewClient(nil)

	_, directoryContent, _, _ := client.Repositories.GetContents(
		ctx,
		"MCDReforged",
		"PluginCatalogue",
		"plugins",
		nil,
	)

	return directoryContent
}

func getMcdrPluginInfo(pluginPath string) (pluginInfo *datatypes.McdrPluginInfo) {
	ctx := context.Background()
	client := github.NewClient(nil)

	fileContent, _, _, _ := client.Repositories.GetContents(
		ctx,
		"MCDReforged",
		"PluginCatalogue",
		path.Join(pluginPath, "plugin_info.json"),
		nil,
	)
	pluginInfo = &datatypes.McdrPluginInfo{}
	content, _ := fileContent.GetContent()
	_ = json.Unmarshal([]byte(content), pluginInfo)

	return
}
