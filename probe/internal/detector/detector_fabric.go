package detector

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"

	externaltype "lucy/exttype"
	"lucy/logger"
	"lucy/syntax"
	"lucy/tools"
	"lucy/types"
)

// fabricServerSingleFileDetector detects Fabric single-file servers
// This is one of the two methods of Fabric installation. One larger .jar file
// it placed at the root of the server directory. It handles the initialization
// and the downloading of the required libraries and minecraft version.
type fabricServerSingleFileDetector struct{}

func (d *fabricServerSingleFileDetector) Name() string {
	return "fabric server"
}

func (d *fabricServerSingleFileDetector) Detect(
	filePath string,
	zipReader *zip.Reader,
	fileHandle *os.File,
) (exec *types.ExecutableInfo, err error) {
	loaderVersion := types.UnknownVersion
	gameVersion := types.UnknownVersion
	for _, f := range zipReader.File {
		if f.Name == "install.properties" {
			r, err := f.Open()
			if err != nil {
				continue
			}
			defer tools.CloseReader(r, logger.Warn)

			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				line := scanner.Text()
				if after, found := strings.CutPrefix(line, "fabric-loader-version="); found {
					loaderVersion = types.RawVersion(after)
				} else if after, found := strings.CutPrefix(line, "game-version="); found {
					gameVersion = types.RawVersion(after)
				}
			}
			if loaderVersion == types.UnknownVersion || gameVersion == types.UnknownVersion {
				continue
			}
			break
		}
	}

	if loaderVersion == types.UnknownVersion || gameVersion == types.UnknownVersion {
		return nil, nil
	}

	exec = &types.ExecutableInfo{
		Path:           filePath,
		GameVersion:    gameVersion,
		LoaderPlatform: types.Fabric,
		LoaderVersion:  loaderVersion,
		BootCommand:    nil,
	}

	return exec, nil
}

// fabricServerLauncherDetector detects Fabric server launchers
// This is one of the two methods of Fabric installation. A lightweight
// launcher .jar file is placed at the root of the server directory. It only
// record the paths to the required libraries.
//
// The detection porcess is rather complicated.
type fabricServerLauncherDetector struct{}

func (d *fabricServerLauncherDetector) Name() string {
	return "fabric server"
}

func (d *fabricServerLauncherDetector) Detect(
	filePath string,
	zipReader *zip.Reader,
	fileHandle *os.File,
) (exec *types.ExecutableInfo, err error) {
	var valid bool
	for _, f := range zipReader.File {
		if f.Name == "fabric-server-launch.properties" {
			r, err := f.Open()
			if err != nil {
				continue
			}
			defer tools.CloseReader(r, logger.Warn)

			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "launch.mainClass=net.fabricmc.loader.impl.launch.knot.KnotServer" {
					valid = true
					break
				}
			}
		}
	}

	if !valid {
		return nil, nil
	}

	loaderVersion := types.UnknownVersion
	gameVersion := types.UnknownVersion
	for _, f := range zipReader.File {
		if f.Name == "META-INF/MANIFEST.MF" {
			r, err := f.Open()
			if err != nil {
				continue
			}
			defer tools.CloseReader(r, logger.Warn)

			err = tools.MoveReaderToLineWithPrefix(r, "Main-Class: ")
			if err != nil {
				continue
			}

			line := bufio.NewScanner(r)
			line.Scan()
			mainClassPaths := strings.Split(
				strings.ReplaceAll(
					line.Text()[len("Main-Class: "):],
					tools.CRLF+" ", ""),
				" ")

			// Here we just parse the paths to find the versions.
			//
			// Although been seemingly unreliable, this is a justified method.
			// The lightweight launcher .jar's idea is to not restrictively
			// specify anything but only the paths to the libraries(classes).
			// Besides, it is the user's responsibility to ensure the presence
			// of the required libraries.
			for _, path := range mainClassPaths {
				if after, found := strings.CutPrefix(path, "libraries/net/fabricmc/fabric-loader/"); found {
					loaderVersion = types.RawVersion(
						strings.Split(after, "/")[0])
				} else if after, found := strings.CutPrefix(path, "libraries/net/fabricmc/intermediary/"); found {
					gameVersion = types.RawVersion(
						strings.Split(after, "/")[0])
				}
			}

			if loaderVersion == types.UnknownVersion || gameVersion == types.UnknownVersion {
				continue
			}

			exec = &types.ExecutableInfo{
				Path:           filePath,
				GameVersion:    gameVersion,
				LoaderPlatform: types.Fabric,
				LoaderVersion:  loaderVersion,
				BootCommand:    nil,
			}

			return exec, nil
		}
	}

	return nil, nil
}

// FabricModDetector detects Fabric mods in JAR files
type FabricModDetector struct{}

func (d *FabricModDetector) Name() string {
	return "fabric mod"
}

func (d *FabricModDetector) Detect(
	zipReader *zip.Reader,
	fileHandle *os.File,
) (packages []types.Package, err error) {
	for _, f := range zipReader.File {
		if f.Name == "fabric.mod.json" {
			r, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer tools.CloseReader(r, logger.Warn)

			data, err := io.ReadAll(r)
			if err != nil {
				return nil, err
			}

			modInfo := &externaltype.FileFabricModIdentifier{}
			err = json.Unmarshal(data, modInfo)
			if err != nil {
				return nil, err
			}

			pkg := types.Package{
				Id: types.PackageId{
					Platform: types.Fabric,
					Name:     syntax.ToProjectName(modInfo.Id),
					Version:  types.RawVersion(modInfo.Version),
				},
				Local: &types.PackageInstallation{
					Path: fileHandle.Name(),
				},
				Dependencies: &types.PackageDependencies{},
				Information:  &types.ProjectInformation{},
			}

			// Parse dependencies
			d.buildDependency(&pkg, modInfo.Depends, true, false)
			d.buildDependency(&pkg, modInfo.Recommends, false, false)
			d.buildDependency(&pkg, modInfo.Suggests, false, false)
			d.buildDependency(&pkg, modInfo.Breaks, true, true)
			d.buildDependency(&pkg, modInfo.Conflicts, false, true)

			// Parse info
			pkg.Information = &types.ProjectInformation{
				Title:       modInfo.Name,
				Description: modInfo.Description,
				License:     modInfo.License,
				Authors: func() []types.Person {
					authors := make([]types.Person, len(modInfo.Authors))
					for i, author := range modInfo.Authors {
						authors[i] = types.Person{Name: author}
					}
					return authors
				}(),
			}

			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

func (d *FabricModDetector) buildDependency(
	pkg *types.Package,
	deps map[string]string,
	mandatory bool,
	inverse bool,
) {
	for k, v := range deps {
		dep := types.Dependency{
			Id: types.PackageId{
				Platform: types.Fabric,
				Name:     syntax.ToProjectName(k),
			},
			Constraint: parseFabricVersionRange(v),
			Mandatory:  mandatory,
		}
		if inverse {
			dep.Constraint.Inverse()
		}
		pkg.Dependencies.Value = append(pkg.Dependencies.Value, dep)
	}
}

func init() {
	RegisterExecutableDetector(&fabricServerSingleFileDetector{})
	RegisterExecutableDetector(&fabricServerLauncherDetector{})
	RegisterModDetector(&FabricModDetector{})
}
