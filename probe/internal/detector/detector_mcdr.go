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

package detector

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"path"

	"lucy/externtype"
	"lucy/syntax"
	"lucy/tools"
	"lucy/types"

	"gopkg.in/yaml.v3"

	"lucy/logger"
)

const mcdrConfigFileName = "config.yml"

// McdrDetector detects MCDR (MCDReforged) installations
type McdrDetector struct{}

func (d *McdrDetector) Name() string {
	return "mcdr"
}

func (d *McdrDetector) Detect(workDir string) any {
	configPath := path.Join(workDir, mcdrConfigFileName)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	// File exists, try to read it
	configFile, err := os.Open(configPath)
	if err != nil {
		logger.Warn(err)
		return nil
	}
	defer func(configFile io.ReadCloser) {
		err := configFile.Close()
		if err != nil {
			logger.Warn(err)
		}
	}(configFile)

	configData, err := io.ReadAll(configFile)
	if err != nil {
		logger.Warn(err)
		return nil
	}

	config := &externtype.McdrConfig{}
	if err := yaml.Unmarshal(configData, config); err != nil {
		logger.Warn(err)
		return nil
	}

	return config
}

func init() {
	RegisterEnvironmentDetector(&McdrDetector{})
}

// GetMcdrConfig uses the new environment detector to check for MCDR installation
var GetMcdrConfig = tools.Memoize(
	func() (config *externtype.McdrConfig) {
		environment := Environment(".")

		if environment.Mcdr != nil {
			return environment.Mcdr.Config
		}

		return nil
	},
)

var GetMcdrPlugins = tools.Memoize(
	func() (plugins []types.Package) {
		plugins = make([]types.Package, 0)
		// Remember that MCDR can have multiple plugin directories
		PluginDirectories := GetMcdrConfig().PluginDirectories
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
				defer tools.CloseReader(pluginFile, logger.Warn)
				if err != nil {
					logger.Warn(err)
					continue
				}
				plugin, err := analyzeMcdrPlugin(pluginFile)
				if err != nil {
					logger.Warn(err)
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
	plugin *types.Package,
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
			pluginInfo := &externtype.McdrPluginIdentifierFile{}
			err = json.Unmarshal(data, pluginInfo)
			if err != nil {
				return nil, err
			}
			return &types.Package{
				Id: types.PackageId{
					Platform: types.Mcdr,
					Name:     syntax.ToProjectName(pluginInfo.Id),
					Version:  types.RawVersion(pluginInfo.Version),
				},
				Local: &types.PackageInstallation{
					Path: file.Name(),
				},
				Dependencies: nil, // TODO: This is not yet implemented, mcdr includes external (python packages) dependencies
			}, nil
		}
	}

	return
}
