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
	Packages     []Package
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
}

type LucyEnv struct {
	ConfigPath string
}
