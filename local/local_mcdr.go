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

package local

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"os"
	"path"

	"lucy/syntax"

	"gopkg.in/yaml.v3"

	"lucy/datatypes"
	"lucy/lucytypes"

	"lucy/lnout"
	"lucy/tools"
)

const mcdrConfigFileName = "config.yml"

// For this part of code, refer to the original MCDR project
// MCDR detects its installation under cwd by check whether the config.yml file exists
// No validation is performed, for empty fields the default value will be filled
// Therefore to align with it, we only detect for the existence of the config.yml file
var getMcdrConfig = tools.Memoize(
	func() (config *McdrConfigDotYml) {
		if _, err := os.Stat(mcdrConfigFileName); os.IsNotExist(err) {
			return nil
		}
		config = &McdrConfigDotYml{}

		configFile, err := os.Open(mcdrConfigFileName)
		if err != nil {
			lnout.Warn(err)
		}

		configData, err := io.ReadAll(configFile)
		defer func(configFile io.ReadCloser) {
			err := configFile.Close()
			if err != nil {
				lnout.Warn(err)
			}
		}(configFile)
		if err != nil {
			lnout.Warn(err)
		}

		if err := yaml.Unmarshal(configData, config); err != nil {
			log.Fatal(err)
		}

		return
	},
)

var getMcdrPlugins = tools.Memoize(
	func() (plugins []lucytypes.Package) {
		plugins = make([]lucytypes.Package, 0)
		// Remember that MCDR can have multiple plugin directories
		PluginDirectories := getMcdrConfig().PluginDirectories
		if PluginDirectories == nil {
			return plugins
		}
		for _, pluginDirectory := range PluginDirectories {
			pluginEntry, _ := os.ReadDir(pluginDirectory)
			for _, pluginPath := range pluginEntry {
				if path.Ext(pluginPath.Name()) != ".mcdr" {
					continue
				}
				pluginFile, err := os.Open(
					path.Join(
						pluginDirectory,
						pluginPath.Name(),
					),
				)
				defer tools.CloseReader(pluginFile, lnout.Warn)
				if err != nil {
					lnout.Warn(err)
					continue
				}
				plugin, err := analyzeMcdrPlugin(pluginFile)
				if err != nil {
					lnout.Warn(err)
					continue
				}
				plugins = append(plugins, *plugin)
			}
		}
		return plugins
	},
)

const mcdrPluginIdentifierFile = "mcdreforged.plugin.json"

func analyzeMcdrPlugin(file *os.File) (
	plugin *lucytypes.Package,
	err error,
) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	r, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil, err
	}

	for _, f := range r.File {
		if f.Name == mcdrPluginIdentifierFile {
			rr, err := f.Open()
			data, err := io.ReadAll(rr)
			if err != nil {
				return nil, err
			}
			pluginInfo := &datatypes.McdrPluginIdentifierFile{}
			err = json.Unmarshal(data, pluginInfo)
			if err != nil {
				return nil, err
			}
			return &lucytypes.Package{
				Id: lucytypes.PackageId{
					Platform: lucytypes.Mcdr,
					Name:     syntax.PackageName(pluginInfo.Id),
					Version:  lucytypes.RawVersion(pluginInfo.Version),
				},
				Local: &lucytypes.PackageInstallation{
					Path: file.Name(),
				},
				Dependencies: nil, // TODO: This is not yet implemented, mcdr includes external (python packages) dependencies
			}, nil
		}
	}

	return
}
