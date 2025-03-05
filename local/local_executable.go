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
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"

	"lucy/dependency"

	"github.com/pelletier/go-toml"

	"lucy/datatypes"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/output"
	"lucy/tools"
)

// TODO: Improve probe logic, plain executable unpacking do not work well
// TODO: Research on forge installation

var getExecutableInfo = tools.Memoize(
	func() *lucytypes.ExecutableInfo {
		var valid []*lucytypes.ExecutableInfo
		workPath := getServerWorkPath()
		jars, err := findJar(workPath)
		if err != nil {
			logger.Warning(err)
			logger.Info("cannot read the current directory, most features will be disabled")
		}
		for _, jar := range jars {
			exec := analyzeExecutable(jar)
			if exec == nil {
				continue
			}
			valid = append(valid, exec)
		}

		if len(valid) == 0 {
			logger.Info("no server jar found, trying to find under libraries")
			jarPaths := findJarRecursive(path.Join(workPath, "libraries"))
			if len(jarPaths) == 0 {
				// if still no jars found in libraries, search the whole directory
				logger.Info("still no server jar found, attempting even more aggressive search")
				logger.Info("note that this may take a long time, and the accuracy is not guaranteed")
				jarPaths = findJarRecursive(workPath)
			}
			mu := sync.Mutex{}
			wg := sync.WaitGroup{}
			for _, jarPath := range jarPaths {
				wg.Add(1)
				go func(jarPath string) {
					exec := analyzeExecutable(jarPath)
					if exec == nil {
						wg.Done()
						return
					}
					mu.Lock()
					valid = append(valid, exec)
					mu.Unlock()
					wg.Done()
				}(jarPath)
			}
			wg.Wait()
		}

		if len(valid) == 0 {
			logger.Info("no server under current directory")
			return UnknownExecutable
		} else if len(valid) == 1 {
			return valid[0]
		}
		index := output.PromptSelectExecutable(valid)
		return valid[index]
	},
)

func findJar(dir string) (jarFiles []string, err error) {
	jarFiles = []string{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if path.Ext(entry.Name()) == ".jar" {
			jarFiles = append(jarFiles, path.Join(dir, entry.Name()))
		}
	}

	return jarFiles, nil
}

const fileCountThreshold = 50000

func findJarRecursive(dir string) (jarFiles []string) {
	jarFiles = []string{}
	entries, _ := os.ReadDir(dir)
	var wg sync.WaitGroup
	var fileCount int32
	var mu sync.Mutex

	// TODO: Use semaphore to limit the number of goroutines
	for _, entry := range entries {
		if atomic.LoadInt32(&fileCount) >= fileCountThreshold {
			logger.Info("file count threshold reached, stopping search")
			break
		}
		if entry.IsDir() {
			wg.Add(1)
			go func(subDir string) {
				defer wg.Done()
				subJarFiles := findJarRecursive(subDir)
				mu.Lock()
				jarFiles = append(jarFiles, subJarFiles...)
				mu.Unlock()
			}(path.Join(dir, entry.Name()))
		} else {
			atomic.AddInt32(&fileCount, 1)
			if path.Ext(entry.Name()) == ".jar" {
				mu.Lock()
				jarFiles = append(jarFiles, path.Join(dir, entry.Name()))
				mu.Unlock()
			}
		}
	}

	wg.Wait()
	return
}

var UnknownExecutable = &lucytypes.ExecutableInfo{
	Path:        "",
	GameVersion: "unknown",
	BootCommand: nil,
	Platform:    lucytypes.UnknownPlatform,
}

const (
	fabricSingleIdentifierFile   = "install.properties"
	vanillaIdentifierFile        = "version.json"
	fabricLauncherIdentifierFile = "fabric-server-launch.properties"
	javaManifest                 = "META-INF/MANIFEST.MF"
	forgeModIdentifierFile       = "META-INF/mods.toml"
)

// analyzeExecutable gives nil if the jar file is invalid. The constant UnknownExecutable
// is not yet used in the codebase, however still reserved for future use.
func analyzeExecutable(filePath string) (exec *lucytypes.ExecutableInfo) {
	// exec is a nil before an analysis function is called
	// Anything other than exec.Path is set in the analysis function
	file, _ := os.Open(filePath)
	if file == nil {
		return nil
	}
	stat, err := file.Stat()
	if err != nil {
		return nil
	}
	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil
	}

	for _, f := range reader.File {
		switch f.Name {
		case fabricSingleIdentifierFile:
			if exec != nil {
				return nil
			}
			exec = analyzeFabricSingle(f)
		case fabricLauncherIdentifierFile:
			if exec != nil {
				return nil
			}
			for _, ff := range reader.File {
				if ff.Name == javaManifest {
					exec = analyzeFabricLauncher(ff)
				}
			}
		case vanillaIdentifierFile:
			if exec != nil {
				return nil
			}
			exec = analyzeVanilla(f)
		case forgeModIdentifierFile:
			if exec != nil {
				return nil
			}
			exec = analyzeForge(f)
		}
	}

	if exec == nil {
		return
	}
	// Set the path to the file at the end
	exec.Path = file.Name()
	return
}

func analyzeVanilla(versionJson *zip.File) (exec *lucytypes.ExecutableInfo) {
	exec = &lucytypes.ExecutableInfo{}
	exec.Platform = lucytypes.Minecraft
	reader, _ := versionJson.Open()
	defer tools.CloseReader(reader, logger.Warning)
	data, _ := io.ReadAll(reader)
	obj := VersionDotJson{}
	_ = json.Unmarshal(data, &obj)
	exec.GameVersion = dependency.RawVersion(obj.Id)
	return
}

// install.properties looks like this:
// fabric-loader-version=0.16.9
// game-version=1.21.4

func analyzeFabricSingle(installProperties *zip.File) (exec *lucytypes.ExecutableInfo) {
	exec = &lucytypes.ExecutableInfo{}
	exec.Platform = lucytypes.Fabric
	r, _ := installProperties.Open()
	defer tools.CloseReader(r, logger.Warning)
	data, _ := io.ReadAll(r)
	s := string(data)

	// Read second line, split by "=" and get the second part
	exec.GameVersion = dependency.RawVersion(
		strings.Split(strings.Split(s, "\n")[1], "=")[1],
	)

	// Read first line, split by "=" and get the second part
	exec.LoaderVersion = dependency.RawVersion(
		strings.Split(strings.Split(s, "\n")[0], "=")[1],
	)

	return
}

// META-INF/MANIFEST.MF looks like this:
// Manifest-Version: 1.0
// Main-Class: net.fabricmc.loader.impl.launch.server.FabricServerLauncher
// Class-Path: libraries/org/ow2/asm/asm/9.7.1/asm-9.7.1.jar libraries/org/
// ow2/asm/asm-analysis/9.7.1/asm-analysis-9.7.1.jar libraries/org/ow2/asm
// /asm-commons/9.7.1/asm-commons-9.7.1.jar libraries/org/ow2/asm/asm-tree
// /9.7.1/asm-tree-9.7.1.jar libraries/org/ow2/asm/asm-util/9.7.1/asm-util
// -9.7.1.jar libraries/net/fabricmc/sponge-mixin/0.15.4+mixin.0.8.7/spong
// e-mixin-0.15.4+mixin.0.8.7.jar libraries/net/fabricmc/intermediary/1.21
// .4/intermediary-1.21.4.jar libraries/net/fabricmc/fabric-loader/0.16.9/
// fabric-loader-0.16.9.jar
// Note that line breaks are "\r\n " and the last line ends with "\r\n"

func analyzeFabricLauncher(
	manifest *zip.File,
) (exec *lucytypes.ExecutableInfo) {
	exec = &lucytypes.ExecutableInfo{}
	exec.Platform = lucytypes.Fabric
	r, _ := manifest.Open()
	defer tools.CloseReader(r, logger.Warning)
	data, _ := io.ReadAll(r)
	s := string(data)
	if !strings.Contains(s, "Class-Path: ") {
		return nil
	}
	s = strings.Split(s, "Class-Path: ")[1] // Start reading from Class-Path
	s = strings.ReplaceAll(s, "\r\n ", "")  // Remove line breaks
	s = strings.ReplaceAll(s, "\r\n", "")   // Remove last line breaks
	classPaths := strings.Split(s, " ")
	for _, classPath := range classPaths {
		if strings.Contains(classPath, "libraries/net/fabricmc/intermediary") {
			exec.GameVersion = dependency.RawVersion(
				strings.Split(classPath, "/")[4],
			)
		}
		if strings.Contains(classPath, "libraries/net/fabricmc/fabric-loader") {
			exec.LoaderVersion = dependency.RawVersion(
				strings.Split(classPath, "/")[4],
			)
		}
	}
	return
}

func analyzeForge(file *zip.File) (exec *lucytypes.ExecutableInfo) {
	r, _ := file.Open()
	defer tools.CloseReader(r, logger.Warning)
	data, _ := io.ReadAll(r)
	p := &datatypes.ForgeModIdentifierNew{}
	err := toml.Unmarshal(data, p)
	if err != nil {
		return nil
	}
	for _, mod := range p.Mods {
		if mod.ModID == "forge" {
			return &lucytypes.ExecutableInfo{
				GameVersion:   dependency.UnknownVersion,
				Platform:      lucytypes.Forge,
				LoaderVersion: mod.Version,
				BootCommand:   nil,
			}
		}
	}

	return nil
}
