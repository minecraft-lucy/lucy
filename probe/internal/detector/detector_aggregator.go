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
	"fmt"
	"os"

	"lucy/logger"
	"lucy/tools"
	"lucy/types"
)

// Executable analyzes a JAR file using all registered detectors
// and returns the first successful match (in registration order).
// If multiple detectors match, callers should handle ambiguity separately.
func Executable(filePath string) *types.ExecutableInfo {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Debug("Failed to open file: " + err.Error())
		return nil
	}
	defer tools.CloseReader(file, logger.Warn)

	stat, err := file.Stat()
	if err != nil {
		logger.Debug("Failed to stat file: " + err.Error())
		return nil
	}

	zipReader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		logger.Debug("Failed to read JAR file: " + err.Error())
		return nil
	}

	var candidates []*types.ExecutableInfo
	detectors := GetExecutableDetectors()

	for _, detector := range detectors {
		result, err := detector.Detect(filePath, zipReader, file)
		if err != nil {
			continue
		}
		if result != nil {
			candidates = append(candidates, result)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	if len(candidates) > 1 {
		// TODO: Modify this by need to handle multiple matches better
		logger.Warn(fmt.Errorf("multiple executable detectors matched; marking as unknown"))
		return UnknownExecutable
	}

	return candidates[0]
}

// Mod analyzes a JAR file for mods using all registered mod detectors
func Mod(filePath string) (res []types.Package) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer tools.CloseReader(file, logger.Warn)

	stat, err := file.Stat()
	if err != nil {
		return nil
	}

	zipReader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil
	}

	detectors := GetModDetectors()

	for _, detector := range detectors {
		packages, err := detector.Detect(zipReader, file)
		if err != nil || packages == nil {
			continue
		}
		res = append(res, packages...)
	}

	return
}

// Environment checks for environment indicators (like MCDR)
func Environment(workDir string) (env types.EnvironmentInfo) {
	detectors := GetEnvironmentDetectors()
	for _, detector := range detectors {
		data := detector.Detect(workDir)
		if detector.Name() == "mcdr" && data != nil {
			env.Mcdr = data.(*types.McdrEnv)
		}
		if detector.Name() == "lucy" && data != nil {
			env.Lucy = data.(*types.LucyEnv)
		}
	}
	return
}
