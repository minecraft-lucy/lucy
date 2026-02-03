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

package types

import (
	"os/exec"

	"lucy/exttype"
)

// ServerInfo components that do not exist, use an empty string. Note Executable
// must exist, otherwise the program will exit; therefore, it is not a pointer.
type ServerInfo struct {
	WorkPath     string
	SavePath     string
	ModPath      string
	Mods         []Package
	Executable   *ExecutableInfo
	Activity     *Activity
	Environments EnvironmentInfo
}

type ExecutableInfo struct {
	Path           string
	GameVersion    RawVersion
	LoaderPlatform Platform
	LoaderVersion  RawVersion
	BootCommand    *exec.Cmd
}

type Activity struct {
	Active bool
	Pid    int
}

type EnvironmentInfo struct {
	Lucy *LucyEnv
	Mcdr *McdrEnv
}

type McdrEnv struct {
	ConfigPath string
	Config     *exttype.FileMcdrConfig
	PluginList []Package
}

type LucyEnv struct {
	ConfigPath string
}
