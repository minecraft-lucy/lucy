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
	"bufio"
	"io"
	"os"
	"strings"

	"lucy/externtype"
	"lucy/logger"
	"lucy/syntax"
	"lucy/tools"
	"lucy/types"

	"github.com/pelletier/go-toml"
)

// forgeServerDetector detects Forge servers
type forgeServerDetector struct{}

func (d *forgeServerDetector) Name() string {
	return "forge server"
}

func (d *forgeServerDetector) Detect(
	filePath string,
	zipReader *zip.Reader,
	fileHandle *os.File,
) (*types.ExecutableInfo, error) {
	forgeVersion := types.UnknownVersion
	gameVersion := types.UnknownVersion
	for _, f := range zipReader.File {
		if f.Name == "META-INF/MANIFEST.MF" {
			r, err := f.Open()
			if err != nil {
				continue
			}
			defer tools.CloseReader(r, logger.Warn)

			err = tools.MoveReaderToLineWithPrefix(r, "Implementation-Title: net.minecraftforge")
			if err != nil {
				continue
			}

			// Read the next line for forge version
			scanner := bufio.NewScanner(r)
			if scanner.Scan() {
				line := scanner.Text()
				if after, found := strings.CutPrefix(line, "Implementation-Version: "); found {
					forgeVersion = types.RawVersion(after)
				}
			}

			// New reader to find game version
			r2, err := f.Open()
			if err != nil {
				continue
			}
			defer tools.CloseReader(r2, logger.Warn)
			err = tools.MoveReaderToLineWithPrefix(r2, "Specification-Title: Minecraft")
			if err != nil {
				continue
			}

			// Read the n+2 line for game version
			scanner2 := bufio.NewScanner(r2)
			scanner2.Scan() // Skip one line
			if scanner2.Scan() {
				line := scanner2.Text()
				if after, found := strings.CutPrefix(line, "Specification-Version: "); found {
					gameVersion = types.RawVersion(after)
				}
			}

			if forgeVersion != types.UnknownVersion && gameVersion != types.UnknownVersion {
				exec := &types.ExecutableInfo{
					Path:           filePath,
					GameVersion:    gameVersion,
					LoaderPlatform: types.Forge,
					LoaderVersion:  forgeVersion,
					BootCommand:    nil,
				}

				return exec, nil
			}
		}
	}

	return nil, nil
}

// forgeModDetector detects new Forge mods (1.13+)
type forgeModDetector struct{}

func (d *forgeModDetector) Name() string {
	return "forge mod"
}

func (d *forgeModDetector) Detect(
	zipReader *zip.Reader,
	fileHandle *os.File,
) (packages []types.Package, err error) {
	for _, f := range zipReader.File {
		if f.Name == "META-INF/mods.toml" {
			r, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer tools.CloseReader(r, logger.Warn)

			data, err := io.ReadAll(r)
			if err != nil {
				return nil, err
			}

			modIdentifier := &externtype.ForgeModIdentifierNew{}
			err = toml.Unmarshal(data, modIdentifier)
			if err != nil {
				return nil, err
			}

			for _, mod := range modIdentifier.Mods {
				// Skip the forge mod itself
				// It will be handled by the executable detector separately
				if mod.ModID == "forge" {
					continue
				}

				// Version
				version := types.RawVersion(mod.Version)
				if version == "${file.jarVersion}" {
					version = getForgeModVersion(zipReader)
				}

				// Parse as internal id
				p := types.Package{
					Id: types.PackageId{
						Platform: types.Forge,
						Name:     syntax.ToProjectName(mod.ModID),
						Version:  version,
					},
					Local: &types.PackageInstallation{
						Path: fileHandle.Name(),
					},
					Dependencies: &types.PackageDependencies{},
					Information:  &types.ProjectInformation{},
				}

				// Parse dependencies
				//
				// This provides an authentic information (rather than a remote).
				// The file is exactly what the loader checks for.
				//
				// Unexpected mod behavior is not our concern. Later we will
				// add manual dependency/conflict management features.
				deps := modIdentifier.Dependencies[mod.ModID]
				for _, dep := range deps {
					p.Dependencies.Value = append(
						p.Dependencies.Value,
						types.Dependency{
							Id: types.PackageId{
								Platform: types.Forge,
								Name:     syntax.ToProjectName(dep.ModID),
							},
							Constraint: parseMavenVersionInterval(dep.VersionRange),
						},
					)
				}

				// Parse info
				p.Information = &types.ProjectInformation{
					Title:   mod.DisplayName,
					Brief:   mod.Description,
					Authors: []types.Person{{Name: mod.Authors}},
					License: modIdentifier.License,
					Urls: []types.Url{
						{
							Name: "URL",
							Type: types.UrlHome,
							Url:  mod.DisplayURL,
						},
						{
							Name: "Issue Tracker",
							Type: types.UrlIssues,
							Url:  modIdentifier.IssueTrackerURL,
						},
					},
				}

				packages = append(packages, p)
			}
		}
	}

	return packages, nil
}

// TODO: Old forge is not yet supported. The detoctor was vibe-coded and needs
// more research.

// ForgeModDetectorOld detects old Forge mods (pre-1.13)
// type ForgeModDetectorOld struct{}

// func (d *ForgeModDetectorOld) Name() string {
// 	return "ForgeModDetectorOld"
// }

// func (d *ForgeModDetectorOld) DetectMod(
// 	zipReader *zip.Reader,
// 	fileHandle *os.File,
// ) []types.Package {
// 	for _, f := range zipReader.File {
// 		if f.Name == "META-INF/mcmod.info" {
// 			r, err := f.Open()
// 			if err != nil {
// 				return nil
// 			}
// 			defer tools.CloseReader(r, logger.Warn)

// 			data, err := io.ReadAll(r)
// 			if err != nil {
// 				return nil
// 			}

// 			modInfos := &datatype.ForgeModIdentifierOld{}
// 			err = json.Unmarshal(data, modInfos)
// 			if err != nil {
// 				return nil
// 			}

// 			var packages []types.Package
// 			for _, modInfo := range *modInfos {
// 				version := types.RawVersion(modInfo.Version)
// 				if version == "${file.jarVersion}" {
// 					version = getForgeModVersion(zipReader)
// 				}

// 				p := types.Package{
// 					Id: types.PackageId{
// 						Platform: types.Forge,
// 						Name:     syntax.ToProjectName(modInfo.ModId),
// 						Version:  version,
// 					},
// 					Local: &types.PackageInstallation{
// 						Path: fileHandle.Name(),
// 					},
// 				}
// 				packages = append(packages, p)
// 			}

// 			if len(packages) > 0 {
// 				return packages
// 			}
// 			return packages
// 		}
// 	}

// 	return nil
// }

func init() {
	RegisterExecutableDetector(&forgeServerDetector{})
	RegisterModDetector(&forgeModDetector{})
}
