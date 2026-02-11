package detector

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"path"

	"lucy/exttype"
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

	config := &exttype.FileMcdrConfig{}
	if err := yaml.Unmarshal(configData, config); err != nil {
		logger.Warn(err)
		return nil
	}

	return config
}

func init() {
	RegisterEnvironmentDetector(&McdrDetector{})
	RegisterModDetector(&McdrPluginDetector{})
}

// GetMcdrConfig uses the new environment detector to check for MCDR installation
var GetMcdrConfig = tools.Memoize(
	func() (config *exttype.FileMcdrConfig) {
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
		pluginDirectories := GetMcdrConfig().PluginDirectories
		if pluginDirectories == nil {
			return plugins
		}
		for _, pluginDirectory := range pluginDirectories {
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
		if f.Name == "mcdreforged.plugin.json" {
			rr, err := f.Open()
			data, err := io.ReadAll(rr)
			if err != nil {
				return nil, err
			}
			pluginInfo := &exttype.FileMcdrPluginIdentifier{}
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

type McdrPluginDetector struct{}

func (d *McdrPluginDetector) Name() string {
	return "mcdr plugin"
}

func (d *McdrPluginDetector) Detect(
	zipReader *zip.Reader,
	fileHandle *os.File,
) (packages []types.Package, err error) {
	for _, f := range zipReader.File {
		if f.Name == "mcdreforged.plugin.json" {
			r, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer tools.CloseReader(r, logger.Warn)

			data, err := io.ReadAll(r)
			if err != nil {
				return nil, err
			}
			pluginInfo := &exttype.FileMcdrPluginIdentifier{}
			if err := json.Unmarshal(data, pluginInfo); err != nil {
				return nil, err
			}

			pkg := types.Package{
				Id: types.PackageId{
					Platform: types.Mcdr,
					Name:     syntax.ToProjectName(pluginInfo.Id),
					Version:  types.RawVersion(pluginInfo.Version),
				},
				Local: &types.PackageInstallation{
					Path: fileHandle.Name(),
				},
				Dependencies: &types.PackageDependencies{},
				Information:  &types.ProjectInformation{},
			}

			// Parse dependencies
			for key, value := range pluginInfo.Dependencies {
				pkg.Dependencies.Value = append(
					pkg.Dependencies.Value,
					types.Dependency{
						Id: types.PackageId{
							Platform: types.Mcdr,
							Name:     syntax.ToProjectName(key),
						},
						Constraint: parseNpmVersionRange(value),
						Mandatory:  true,
					},
				)
			}

			// Parse info
			pkg.Information.Authors = make(
				[]types.Person,
				len(pluginInfo.Author),
			)
			for i, author := range pluginInfo.Author {
				pkg.Information.Authors[i] = types.Person{
					Name: author,
				}
			}
			pkg.Information.Title = pluginInfo.Name
			pkg.Information.Brief = pluginInfo.Description.EnUs
			pkg.Information.Urls = []types.Url{
				{
					Name: "Link",
					Type: types.UrlSource,
					Url:  pluginInfo.Link,
				},
			}
		}
	}

	return packages, nil
}
