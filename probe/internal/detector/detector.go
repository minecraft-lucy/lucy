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
