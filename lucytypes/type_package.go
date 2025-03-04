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

package lucytypes

type PackageUrlType uint8

const (
	FileUrl PackageUrlType = iota
	HomepageUrl
	SourceUrl
	WikiUrl
	OthersUrl
)

func (p PackageUrlType) String() string {
	switch p {
	case FileUrl:
		return "File"
	case HomepageUrl:
		return "Homepage"
	case SourceUrl:
		return "Source"
	case WikiUrl:
		return "Wiki"
	case OthersUrl:
		return "URL"
	default:
		return "Unknown"
	}
}

type PackageUrl struct {
	Name string
	Type PackageUrlType
	Url  string
}

// Package is a package identifier with its related information. In principle,
// only packages remote and local can provide a Package.
//
// This is an adapter type that uses composition method to provide a unified
// interface for both local and remote packages. It is used to represent a
// package in the system, and can be used to store information about the package
// such as its dependencies, installation path, and remote source.
type Package struct {
	// Id is the basic package identifier
	Id PackageId

	// Optional attributions
	Information  *PackageInformation
	Dependencies *PackageDependencies
	Local        *PackageInstallation
	Remote       *PackageRemote
}

// PackageDependencies is one of the optional attributions that can be added to
// a Package struct. It is usually used in any command that requires operating
// local packages, such as `lucy install` or `lucy remove`.
type PackageDependencies struct {
	SupportedVersions  []RawVersion
	SupportedPlatforms []Platform
	Required           []PackageId
	Optional           []PackageId
	Incompatible       []PackageId
}

// PackageInformation is a struct that contains informational data about the
// package. It is typically used in `lucy info`.
type PackageInformation struct {
	Name        string
	Brief       string
	Description string
	Author      []PackageMember
	Urls        []PackageUrl
	License     string
}

type PackageMember struct {
	Name  string
	Role  string
	Url   string
	Email string
}

// PackageInstallation is an optional attribution to lucytypes.Package. It is
// used for packages that are known to be installed in the local filesystem.
type PackageInstallation struct {
	Path string
}

// PackageRemote is an optional attribution to lucytypes.Package. It is used to
// represent package's presence in a remote source.
type PackageRemote struct {
	Source Source
	// Whatever the remote source uses to identify this package.
	RemoteId string
	// The URL to download the package's specified version When package.Id.Version
	FileUrl  string
	Filename string
}

// PackageUpdate is a struct to represent the update status of a package. It must
// be used with PackageRemote.
type PackageUpdate struct {
	// Just refer this to the PackageRemote under the same Package struct.
	Current *PackageRemote
}
