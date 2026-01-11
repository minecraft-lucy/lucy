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

package mcdr

import (
	"lucy/types"
	"time"
)

// GitHub API file ref: https://api.github.com/repos/MCDReforged/PluginCatalogue/contents/plugins/{plugin_name}/plugin_info.json
// The purpose of this file is quite unclear to me.
// For this project, meta.json under the meta branch is more handy.
type pluginInfo struct {
	Id           string   `json:"id"`
	Authors      []author `json:"authors"`
	Repository   string   `json:"repository"`
	Branch       string   `json:"branch"`
	RelatedPath  string   `json:"related_path"`
	Labels       []string `json:"labels"`
	Introduction struct {
		EnUs string `json:"en_us"`
		ZhCn string `json:"zh_cn"`
	} `json:"introduction"`
}

type author struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

// GitHub API file ref: https://api.github.com/repos/MCDReforged/PluginCatalogue/{plugin_name}/release.json?ref=meta
type pluginRelease struct {
	SchemaVersion      int       `json:"schema_version"`
	Id                 string    `json:"id"`
	LatestVersion      string    `json:"latest_version"`
	LatestVersionIndex int       `json:"latest_version_index"`
	Releases           []release `json:"releases"`
}
type release struct {
	Url         string     `json:"url"`
	Name        string     `json:"name"`
	TagName     string     `json:"tag_name"`
	CreatedAt   time.Time  `json:"created_at"`
	Description string     `json:"description"`
	Prerelease  bool       `json:"prerelease"`
	Asset       asset      `json:"asset"`
	Meta        pluginMeta `json:"meta"`
}

func (r release) ToPackageRemote() types.PackageRemote {
	remote := types.PackageRemote{
		Source:   types.McdrCatalogue,
		FileUrl:  r.Asset.BrowserDownloadUrl,
		Filename: r.Asset.Name,
	}
	return remote
}

type asset struct {
	Id                 int       `json:"id"`
	Name               string    `json:"name"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	BrowserDownloadUrl string    `json:"browser_download_url"`
	HashMd5            string    `json:"hash_md5"`
	HashSha256         string    `json:"hash_sha256"`
}

// GitHub API file ref: https://api.github.com/repos/MCDReforged/PluginCatalogue/{plugin_name}/meta.json?ref=meta
type pluginMeta struct {
	SchemaVersion int               `json:"schema_version"`
	Id            string            `json:"id"`
	Name          string            `json:"name"`
	Version       string            `json:"version"`
	Link          string            `json:"link"`
	Authors       []string          `json:"authors"`
	Dependencies  map[string]string `json:"dependencies"`
	Requirements  []string          `json:"requirements"`
	Description   struct {
		EnUs string `json:"en_us"`
		ZhCn string `json:"zh_cn"`
	} `json:"description"`
}
