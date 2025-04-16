package mcdr

import (
	"time"

	"lucy/lucytypes"
	"lucy/tools"
)

type everything struct {
	Timestamp int               `json:"timestamp"`
	Authors   map[string]author `json:"authors"`
	Plugins   map[string]plugin `json:"plugins"`
}

func (e everything) ToSearchResults() lucytypes.SearchResults {
	res := lucytypes.SearchResults{}
	res.Source = lucytypes.McdrCatalogue
	for _, plugin := range e.Plugins {
		res.Results = append(
			res.Results,
			lucytypes.ToProjectName(plugin.Meta.Id),
		)
	}
	return res
}

type author struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type plugin struct {
	Meta struct {
		SchemaVersion int           `json:"schema_version"`
		Id            string        `json:"id"`
		Name          string        `json:"name"`
		Version       string        `json:"version"`
		Link          string        `json:"link"`
		Authors       []string      `json:"authors"`
		Dependencies  struct{}      `json:"dependencies"`
		Requirements  []interface{} `json:"requirements"`
		Description   struct {
			EnUs string `json:"en_us"`
			ZhCn string `json:"zh_cn"`
		} `json:"description"`
	} `json:"meta"`
	Plugin struct {
		SchemaVersion int      `json:"schema_version"`
		Id            string   `json:"id"`
		Authors       []string `json:"authors"`
		Repository    string   `json:"repository"`
		Branch        string   `json:"branch"`
		RelatedPath   string   `json:"related_path"`
		Labels        []string `json:"labels"`
		Introduction  struct {
			EnUs string `json:"en_us"`
			ZhCn string `json:"zh_cn"`
		} `json:"introduction"`
		IntroductionUrls struct {
			EnUs string `json:"en_us"`
			ZhCn string `json:"zh_cn"`
		} `json:"introduction_urls"`
	} `json:"plugin"`
	Release struct {
		SchemaVersion      int    `json:"schema_version"`
		Id                 string `json:"id"`
		LatestVersion      string `json:"latest_version"`
		LatestVersionIndex int    `json:"latest_version_index"`
		Releases           []struct {
			Url         string    `json:"url"`
			Name        string    `json:"name"`
			TagName     string    `json:"tag_name"`
			CreatedAt   time.Time `json:"created_at"`
			Description string    `json:"description"`
			Prerelease  bool      `json:"prerelease"`
			Asset       struct {
				Id                 int       `json:"id"`
				Name               string    `json:"name"`
				Size               int       `json:"size"`
				DownloadCount      int       `json:"download_count"`
				CreatedAt          time.Time `json:"created_at"`
				BrowserDownloadUrl string    `json:"browser_download_url"`
				HashMd5            string    `json:"hash_md5"`
				HashSha256         string    `json:"hash_sha256"`
			} `json:"asset"`
			Meta struct {
				SchemaVersion int           `json:"schema_version"`
				Id            string        `json:"id"`
				Name          string        `json:"name"`
				Version       string        `json:"version"`
				Link          string        `json:"link"`
				Authors       []string      `json:"authors"`
				Dependencies  struct{}      `json:"dependencies"`
				Requirements  []interface{} `json:"requirements"`
				Description   struct {
					EnUs string `json:"en_us"`
					ZhCn string `json:"zh_cn"`
				} `json:"description"`
			} `json:"meta"`
		} `json:"releases"`
	} `json:"release"`
	Repository struct {
		Url             string `json:"url"`
		Name            string `json:"name"`
		FullName        string `json:"full_name"`
		HtmlUrl         string `json:"html_url"`
		Description     string `json:"description"`
		Archived        bool   `json:"archived"`
		StargazersCount int    `json:"stargazers_count"`
		WatchersCount   int    `json:"watchers_count"`
		ForksCount      int    `json:"forks_count"`
		Readme          string `json:"readme"`
		ReadmeUrl       string `json:"readme_url"`
		License         struct {
			Key    string `json:"key"`
			Name   string `json:"name"`
			SpdxId string `json:"spdx_id"`
			Url    string `json:"url"`
		} `json:"license"`
	} `json:"repository"`
}

func (p plugin) ToProjectInformation() lucytypes.ProjectInformation {
	info := lucytypes.ProjectInformation{
		Title:       p.Meta.Name,
		Brief:       p.Plugin.Introduction.EnUs,
		Description: p.Meta.Description.EnUs,
		Author:      make([]lucytypes.PackageMember, 0, len(p.Plugin.Authors)),
		Urls:        make([]lucytypes.PackageUrl, 0),
		License:     p.Repository.License.Name,
	}

	// authors
	for _, authorName := range p.Meta.Authors {
		name := tools.Ternary(
			getAuthor(authorName) == nil,
			authorName,
			getAuthor(authorName).Name,
		)
		url := tools.Ternary(
			getAuthor(authorName) == nil,
			"",
			getAuthor(authorName).Link,
		)
		info.Author = append(
			info.Author, lucytypes.PackageMember{
				Name:  name,
				Role:  "Author",
				Url:   url,
				Email: "",
			},
		)
	}

	// urls
	info.Urls = append(
		info.Urls, lucytypes.PackageUrl{
			Name: "GitHub",
			Type: lucytypes.SourceUrl,
			Url:  p.Repository.Url,
		},
	)

	info.Urls = append(
		info.Urls, lucytypes.PackageUrl{
			Name: "Latest Release",
			Type: lucytypes.FileUrl,
			Url:  p.Release.Releases[0].Asset.BrowserDownloadUrl,
		},
	)

	return info
}
