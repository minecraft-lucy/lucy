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

// Package util is a general package for network and file system operations.
package util

import (
	"fmt"
	"os"
	"path"

	"lucy/logger"
)

func init() {
	if err := os.MkdirAll(CacheDir, os.ModePerm); err != nil {
		DoCache = false
		logger.Warn(
			fmt.Errorf(
				"failed to create cache directory, disabling caching: %w",
				err,
			),
		)
	}
	DoCache = true
}

const (
	ProgramName  = "lucy"
	ProgramPath  = ".lucy"
	ConfigFile   = ProgramPath + "/config.json"
	DownloadPath = ProgramPath + "/downloads"
)

var CacheDir = path.Join(os.TempDir(), ProgramName)

var DoCache bool

// TODO: implement download task

// DownloadTask return a function that downloads a file from url to path.
func DownloadTask(from string, to string) func() error {
	return func() error {
		return nil
	}
}
