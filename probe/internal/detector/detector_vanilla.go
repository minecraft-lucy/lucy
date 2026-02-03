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

	"lucy/exttype"
	"lucy/logger"
	"lucy/tools"
	"lucy/types"
)

// VanillaDetector detects vanilla Minecraft servers
type VanillaDetector struct{}

func (d *VanillaDetector) Name() string {
	// TODO implement me
	panic("implement me")
}

func (d *VanillaDetector) Detect(
	filePath string,
	zipReader *zip.Reader,
	fileHandle *os.File,
) (*types.ExecutableInfo, error) {
	for _, f := range zipReader.File {
		if f.Name == "version.json" {
			r, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer tools.CloseReader(r, logger.Warn)

			data, err := io.ReadAll(r)
			if err != nil {
				return nil, err
			}

			obj := exttype.FileMinecraftVersionSpec{}
			err = json.Unmarshal(data, &obj)
			if err != nil {
				return nil, err
			}

			exec := &types.ExecutableInfo{
				Path:           filePath,
				LoaderPlatform: types.Minecraft,
				GameVersion:    types.RawVersion(obj.Id),
			}

			return exec, nil
		}
	}

	return nil, nil
}

func init() {
	RegisterExecutableDetector(&VanillaDetector{})
}
