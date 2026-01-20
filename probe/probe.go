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

// Package local provides functionality to gather and manage server information
// for a Minecraft server. It includes methods to retrieve server configuration,
// mod list, executable information, and other relevant details. The package
// utilizes memoization to avoid redundant calculations and resolve any data
// dependencies issues. Therefore, all probe functions are 100% concurrent-safe.
//
// The main exposed function is GetServerInfo, which returns a comprehensive
// ServerInfo struct containing all the gathered information. To avoid side
// effects, the ServerInfo struct is returned as a copy, rather than reference.
package probe

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"sync"

	"lucy/syntax"

	"github.com/pelletier/go-toml"

	"gopkg.in/ini.v1"

	"lucy/datatype"
	"lucy/logger"
	"lucy/tools"
	"lucy/types"
)

// GetServerInfo is the exposed function for external packages to get serverInfo.
// As we can assume that the environment does not change while the program is
// running, a sync.Once is used to prevent further calls to this function. Rather,
// the cached serverInfo is used as the return value.
var GetServerInfo = tools.Memoize(buildServerInfo)

// buildServerInfo builds the server information by performing several checks
// and gathering data from various sources. It uses goroutines to perform these
// tasks concurrently and a sync.Mutex to ensure thread-safe updates to the
// serverInfo struct.
func buildServerInfo() types.ServerInfo {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var serverInfo types.ServerInfo

	// MCDR Stage
	wg.Add(1)
	go func() {
		defer wg.Done()
		mcdrConfig := getMcdrConfig()
		if mcdrConfig != nil {
			mu.Lock()
			serverInfo.Mcdr = &types.McdrInstallation{
				PluginPaths: mcdrConfig.PluginDirectories,
			}
			mu.Unlock()
		}
	}()

	// MCDR Plugins
	if getMcdrConfig() != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			plugins := getMcdrPlugins()
			mu.Lock()
			serverInfo.Mcdr.PluginList = plugins
			mu.Unlock()
		}()
	}

	// Server Work Path
	wg.Add(1)
	go func() {
		defer wg.Done()
		workPath := getServerWorkPath()
		mu.Lock()
		serverInfo.WorkPath = workPath
		mu.Unlock()
	}()

	// Executable Stage
	wg.Add(1)
	go func() {
		defer wg.Done()
		executable := getExecutableInfo()
		mu.Lock()
		serverInfo.Executable = executable
		mu.Unlock()
	}()

	// Mod Path
	wg.Add(1)
	go func() {
		defer wg.Done()
		modPath := getServerModPath()
		mu.Lock()
		serverInfo.ModPath = modPath
		mu.Unlock()
	}()

	// Mod List
	wg.Add(1)
	go func() {
		defer wg.Done()
		modList := getMods()
		mu.Lock()
		serverInfo.Mods = modList
		mu.Unlock()
	}()

	// Save Path
	wg.Add(1)
	go func() {
		defer wg.Done()
		savePath := getSavePath()
		mu.Lock()
		serverInfo.SavePath = savePath
		mu.Unlock()
	}()

	// Check for Lucy installation
	wg.Add(1)
	go func() {
		defer wg.Done()
		hasLucy := checkHasLucy()
		mu.Lock()
		serverInfo.HasLucy = hasLucy
		mu.Unlock()
	}()

	// Check if the server is running
	wg.Add(1)
	go func() {
		defer wg.Done()
		activity := checkServerFileLock()
		mu.Lock()
		serverInfo.Activity = activity
		mu.Unlock()
	}()

	// Server Mod Path
	wg.Add(1)
	go func() {
		defer wg.Done()
		modPath := getServerModPath()
		mu.Lock()
		serverInfo.ModPath = modPath
		mu.Unlock()
	}()

	wg.Wait()
	return serverInfo
}

// Some functions that gets a single piece of information. They are not exported,
// as GetServerInfo() applies a memoization mechanism. Every time a serverInfo
// is needed, just call GetServerInfo() without the concern of redundant calculation.

var getServerModPath = tools.Memoize(
	func() string {
		if exec := getExecutableInfo(); exec != nil && (exec.LoaderPlatform == types.Fabric || exec.LoaderPlatform == types.Forge) {
			return path.Join(getServerWorkPath(), "mods")
		}
		return ""
	},
)

var getServerWorkPath = tools.Memoize(
	func() string {
		if mcdrConfig := getMcdrConfig(); mcdrConfig != nil {
			return mcdrConfig.WorkingDirectory
		}
		return "."
	},
)

var getServerDotProperties = tools.Memoize(
	func() MinecraftServerDotProperties {
		exec := getExecutableInfo()
		propertiesPath := path.Join(getServerWorkPath(), "server.properties")
		file, err := ini.Load(propertiesPath)
		if err != nil {
			if exec != UnknownExecutable {
				logger.Warn(errors.New("this server is missing a server.properties"))
			}
			return nil
		}

		properties := make(map[string]string)
		for _, section := range file.Sections() {
			for _, key := range section.Keys() {
				properties[key.Name()] = key.String()
			}
		}

		return properties
	},
)

var getSavePath = tools.Memoize(
	func() string {
		serverProperties := getServerDotProperties()
		if serverProperties == nil {
			return ""
		}
		levelName := serverProperties["level-name"]
		return path.Join(getServerWorkPath(), levelName)
	},
)

var checkHasLucy = tools.Memoize(
	func() bool {
		_, err := os.Stat(".lucy")
		return err == nil
	},
)

var getMods = tools.Memoize(
	func() (mods []types.Package) {
		path := getServerModPath()
		jarPaths, err := findJar(path)
		if err != nil {
			logger.Warn(err)
			logger.Info("this server might not have a mod folder")
			return nil
		}

		for _, jarPath := range jarPaths {
			jar, err := os.Open(jarPath)
			if err != nil {
				continue
			}
			analyzed := analyzeModJar(jar)
			if analyzed != nil {
				mods = append(mods, analyzed...)
			}
		}

		sort.Slice(
			mods,
			func(i, j int) bool { return mods[i].Id.Name < mods[j].Id.Name },
		)
		return mods
	},
)

const (
	fabricModIdentifierFile   = "fabric.mod.json"
	oldForgeModIdentifierFile = "META-INF/mcmod.info"
	newForgeModIdentifierFile = "META-INF/mods.toml"
)

// analyzeModJar is now a single large function, but it will be later split into
// smaller functions according to different mod loaders. This function will keep
// serve as an entry point to the mod analysis process.
//
// According to current information, all mod analysis can be summarized into the
// following process:
// 1. Check for the identifier file
// 2. Analyze informative files
// 3. Fill in the Package struct
func analyzeModJar(file *os.File) (packages []types.Package) {
	stat, err := file.Stat()
	if err != nil {
		return nil
	}
	zipReader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil
	}

	packages = []types.Package{}

	for _, f := range zipReader.File {
		// fabric check
		if f.Name == fabricModIdentifierFile {
			rr, err := f.Open()
			data, err := io.ReadAll(rr)
			modInfo := &datatype.FabricModIdentifier{}
			err = json.Unmarshal(data, modInfo)
			if err != nil {
				return nil
			}

			packages = append(
				packages,
				types.Package{
					Id: types.PackageId{
						Platform: types.Fabric,
						Name:     syntax.PackageName(modInfo.Id),
						Version:  types.RawVersion(modInfo.Version),
					},
					Local: &types.PackageInstallation{
						Path: file.Name(),
					},
					Dependencies: nil, // TODO: This is not yet implemented, because the deps field is an expression, we need to parse it
				},
			)

			return packages
		}

		// check for old forge identifier
		if f.Name == oldForgeModIdentifierFile {
			rr, err := f.Open()
			data, err := io.ReadAll(rr)
			if err != nil {
				return nil
			}
			modInfos := &datatype.ForgeModIdentifierOld{}
			err = json.Unmarshal(data, modInfos)
			if err != nil {
				return nil
			}

			for _, modInfo := range *modInfos {
				p := types.Package{
					Id: types.PackageId{
						Platform: types.Forge,
						Name:     syntax.PackageName(modInfo.ModId),
						Version:  types.RawVersion(modInfo.Version),
					},
					Local: &types.PackageInstallation{
						Path: file.Name(),
					},
					Dependencies: nil, // TODO: This is not yet implemented, because the deps field is an expression, we need to parse it
				}
				if p.Id.Version == "${file.jarVersion}" {
					p.Id.Version = getForgeVariableVersion(zipReader)
				}
				packages = append(packages, p)
			}

			return packages
		}

		// check for new forge identifier
		if f.Name == newForgeModIdentifierFile {
			rr, err := f.Open()
			data, err := io.ReadAll(rr)
			if err != nil {
				return nil
			}
			modInfo := &datatype.ForgeModIdentifierNew{}

			err = toml.Unmarshal(data, modInfo)
			if err != nil {
				return nil
			}

			for _, mod := range modInfo.Mods {
				p := types.Package{
					Id: types.PackageId{
						Platform: types.Forge,
						Name:     types.ProjectName(mod.ModID),
						Version:  types.RawVersion(mod.Version),
					},
					Local: &types.PackageInstallation{
						Path: file.Name(),
					},
					Dependencies: nil, // TODO: This is not yet implemented, because the deps field is an expression, we need to parse it
				}
				if p.Id.Version == "${file.jarVersion}" {
					p.Id.Version = getForgeVariableVersion(zipReader)
				}
				packages = append(packages, p)
			}

			return packages
		}

	}

	return nil
}

func getForgeVariableVersion(zip *zip.Reader) types.RawVersion {
	var r io.ReadCloser
	var err error
	for _, f := range zip.File {
		if f.Name == javaManifest {
			r, err = f.Open()
			if err != nil {
				return types.UnknownVersion
			}
		}
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return types.UnknownVersion
	}
	manifest := string(data)
	const versionField = "Implementation-Version: "
	i := strings.Index(manifest, versionField) + len(versionField)
	if i == -1 {
		return types.UnknownVersion
	}
	v := manifest[i:]
	v = strings.Split(v, "\r")[0]
	v = strings.Split(v, "\n")[0]
	return types.RawVersion(v)
}
