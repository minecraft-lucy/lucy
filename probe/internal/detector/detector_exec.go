package detector

import (
	"lucy/types"
)

var UnknownExecutable = &types.ExecutableInfo{
	Path:           "",
	GameVersion:    "unknown",
	BootCommand:    nil,
	LoaderPlatform: types.UnknownPlatform,
}
