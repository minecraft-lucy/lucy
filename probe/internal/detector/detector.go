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
	"os"

	"lucy/types"
)

// ExecutableDetector is the interface for detecting different types of
// Minecraft servers
type ExecutableDetector interface {
	Detect(
		filePath string,
		zipReader *zip.Reader,
		fileHandle *os.File,
	) (*types.ExecutableInfo, error)
	Name() string
}

// ModDetector is the interface for analyzing mods or plugins
type ModDetector interface {
	Detect(zipReader *zip.Reader, fileHandle *os.File) ([]types.Package, error)
	Name() string
}

// EnvironmentDetector is the detector that handles a directory rather than
// a single file
type EnvironmentDetector interface {
	Detect(workDir string) any
	Name() string
}
