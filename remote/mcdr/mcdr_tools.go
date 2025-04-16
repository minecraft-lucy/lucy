package mcdr

import (
	"lucy/lucytypes"
	"lucy/tools"
	"slices"
	"time"
)

const matchThreshold = 0.266667

// match is a helper function to for self.Search, as the MCDR API just gives
// the whole catalogue in a single file, we need to filter the results by
// query.
//
// This is in-place.
func match(
everything *everything,
query string,
) {
	plugins := everything.Plugins
	for i, plugin := range plugins {
		id := lucytypes.ToProjectName(plugin.Meta.Id).String()
		normDist := tools.NormalizedLevenshteinDistance(id, query)
		if normDist > matchThreshold {
			delete(plugins, i)
		}
	}
}

// sortBy has a VERY BAD implementation.
// TODO: Refactor
func sortBy(
everything *everything,
index lucytypes.SearchIndex,
) (res []lucytypes.ProjectName) {
	n := len(everything.Plugins)
	switch index {
	case lucytypes.ByRelevance:
		type keyValueArray []struct {
			item  lucytypes.ProjectName
			index float64
		}
		var arr keyValueArray
		return cmpByRelevance(arr)
	case lucytypes.ByDownloads:
		type keyValueArray []struct {
			item  lucytypes.ProjectName
			index int
		}
		var arr = make(keyValueArray, 0, n)
		for _, plugin := range everything.Plugins {
			download := 0
			for _, release := range plugin.Release.Releases {
				download += release.Asset.DownloadCount
			}
			arr = append(
				arr,
				struct {
					item  lucytypes.ProjectName
					index int
				}{lucytypes.ToProjectName(plugin.Meta.Id), download},
			)
		}
		return cmpByDownloads(arr)
	case lucytypes.ByNewest:
		type keyValueArray []struct {
			item  lucytypes.ProjectName
			index time.Time
		}
		var arr = make(keyValueArray, 0, n)
		for _, plugin := range everything.Plugins {
			timestamp := plugin.Release.Releases[0].CreatedAt
			arr = append(
				arr, struct {
					item  lucytypes.ProjectName
					index time.Time
				}{lucytypes.ToProjectName(plugin.Meta.Id), timestamp},
			)
		}
		return cmpByNewest(arr)
	}

	return nil
}

func cmpByRelevance(
arr []struct {
	item  lucytypes.ProjectName
	index float64
},
) (res []lucytypes.ProjectName) {
	slices.SortFunc(
		arr, func(
		a, b struct {
			item  lucytypes.ProjectName
			index float64
		},
		) int {
			if a.index == b.index {
				return 0
			}
			if a.index > b.index {
				return -1
			}
			return 1
		},
	)

	for _, item := range arr {
		res = append(res, item.item)
	}
	return res
}

func cmpByDownloads(
arr []struct {
	item  lucytypes.ProjectName
	index int
},
) (res []lucytypes.ProjectName) {
	slices.SortFunc(
		arr, func(
		a, b struct {
			item  lucytypes.ProjectName
			index int
		},
		) int {
			return b.index - a.index
		},
	)

	for _, item := range arr {
		res = append(res, item.item)
	}
	return res
}

func cmpByNewest(
arr []struct {
	item  lucytypes.ProjectName
	index time.Time
},
) (res []lucytypes.ProjectName) {
	slices.SortFunc(
		arr, func(
		a, b struct {
			item  lucytypes.ProjectName
			index time.Time
		},
		) int {
			if a.index.Equal(b.index) {
				return 0
			}
			if a.index.After(b.index) {
				return -1
			}
			return 1
		},
	)

	for _, item := range arr {
		res = append(res, item.item)
	}
	return res
}
