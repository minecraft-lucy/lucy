package mcdr

import (
	"encoding/json"
	"fmt"

	"lucy/github"
	"lucy/types"

	"github.com/sahilm/fuzzy"
)

const (
	pluginCatalogueRepoEndpoint = `https://api.github.com/repos/MCDReforged/PluginCatalogue/contents/`
	branchMaster                = "?ref=master"
	branchCatalogue             = "?ref=catalogue" // I haven't figured out the difference yet
	branchMeta                  = "?ref=meta"
)

func searchPlugin(query string) (mcdrSearchResult, error) {
	ghEndpoint := pluginCatalogueRepoEndpoint + ("plugins/") + branchCatalogue
	err, msg, data := github.GetFileFromGitHub(ghEndpoint)
	if err != nil {
		return nil, err
	}
	if msg != nil && msg.Message != "" {
		return nil, fmt.Errorf("%w: %s", ErrorGhApi, msg.Message)
	}

	var ghFiles []github.GhItem
	err = json.Unmarshal(data, &ghFiles)
	if err != nil {
		return nil, err
	}

	pluginIds := make([]string, 0)
	for _, file := range ghFiles {
		pluginIds = append(pluginIds, file.Name)
	}

	matches := fuzzy.Find(query, pluginIds)
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		result = append(result, pluginIds[match.Index])
	}
	return result, nil
}

func getPluginInfo(id string) (*pluginInfo, error) {
	ghEndpoint := pluginCatalogueRepoEndpoint + ("plugins/") + id + "/plugin_info.json" + branchMaster
	var data []byte
	err, msg, data := github.GetFileFromGitHub(ghEndpoint)
	if err != nil {
		return nil, err
	}
	if msg != nil && msg.Message != "" {
		if msg.Status == "404" {
			return nil, ErrPluginNotFound
		}
		return nil, fmt.Errorf("%w: %s", ErrorGhApi, msg.Message)
	}

	var info pluginInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func getMeta(id string) (*pluginMeta, error) {
	ghEndpoint := pluginCatalogueRepoEndpoint + id + "/meta.json" + branchMeta
	var data []byte
	err, msg, data := github.GetFileFromGitHub(ghEndpoint)
	if err != nil {
		return nil, err
	}
	if msg != nil && msg.Message != "" {
		if msg.Status == "404" {
			return nil, ErrPluginNotFound
		}
		return nil, fmt.Errorf("%w: %s", ErrorGhApi, msg.Message)
	}

	var meta pluginMeta
	err = json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

func getRelease(id string, version types.RawVersion) (*release, error) {
	history, err := getReleaseHistory(id)
	if err != nil {
		return nil, err
	}

	if version == types.LatestVersion {
		return &history.Releases[history.LatestVersionIndex], nil
	}

	for _, rel := range history.Releases {
		if rel.Meta.Version == version.String() {
			return &rel, nil
		}
	}
	return nil, fmt.Errorf(
		"%w for plugin %s: %s",
		ErrVersionNotFound,
		id,
		version,
	)
}

func getReleaseHistory(id string) (*pluginRelease, error) {
	ghEndpoint := pluginCatalogueRepoEndpoint + id + "/release.json" + branchMeta
	var data []byte
	err, msg, data := github.GetFileFromGitHub(ghEndpoint)
	if err != nil {
		return nil, err
	}
	if msg != nil && msg.Message != "" {
		if msg.Status == "404" {
			return nil, ErrPluginNotFound
		}
		return nil, fmt.Errorf("%w: %s", ErrorGhApi, msg.Message)
	}

	var releaseHistory pluginRelease
	err = json.Unmarshal(data, &releaseHistory)
	if err != nil {
		return nil, err
	}
	return &releaseHistory, nil
}

func getRepositoryInfo(id string) (*pluginRepo, error) {
	ghEndpoint := pluginCatalogueRepoEndpoint + id + "/repository.json" + branchMeta
	var data []byte
	err, msg, data := github.GetFileFromGitHub(ghEndpoint)
	if err != nil {
		return nil, err
	}
	if msg != nil && msg.Message != "" {
		if msg.Status == "404" {
			return nil, ErrPluginNotFound
		}
		return nil, fmt.Errorf("%w: %s", ErrorGhApi, msg.Message)
	}

	var repo pluginRepo
	err = json.Unmarshal(data, &repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}
