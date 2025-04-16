package mcdr

import (
	"fmt"
	"time"

	"lucy/lucytypes"
	"lucy/tools"

	"github.com/sahilm/fuzzy"
)

// match is a helper function to for self.Search, as the MCDR API just gives
// the whole catalogue in a single file, we need to filter the results by
// query.
//
// This has a side effect to getEverything()
func match(
	query string,
) (err error) {
	everything, err := getEverything()
	if err != nil {
		return err
	}
	var ids = make([]string, 0, len(everything.Plugins))
	for id := range everything.Plugins {
		ids = append(ids, id)
	}

	matches := fuzzy.Find(query, ids)

	var matchedPlugins = make(map[string]plugin, len(matches))
	for _, match := range matches {
		matchedPlugins[match.Str] = everything.Plugins[match.Str]
	}
	everything.Plugins = matchedPlugins

	return nil
}

// sortBy is a helper function to sort the plugins by the given index.
//
// Side effects on getEverything() will disturb this function.
func sortBy(
	index lucytypes.SearchIndex,
) (res []lucytypes.ProjectName, err error) {
	everything, err := getEverything()
	if err != nil {
		return nil, err
	}
	switch index {
	case lucytypes.ByRelevance:
		arr := make(
			[]tools.KeyValue[lucytypes.ProjectName, float64],
			0,
			len(everything.Plugins),
		)
		for _, plugin := range everything.Plugins {
			id := lucytypes.ToProjectName(plugin.Meta.Id)
			relevance := tools.NormalizedLevenshteinDistance(id.String(), "")
			arr = append(
				arr,
				tools.KeyValue[lucytypes.ProjectName, float64]{
					Item: id, Index: relevance,
				},
			)
		}
		return tools.SortAndExtract(
			arr,
			func(a, b tools.KeyValue[lucytypes.ProjectName, float64]) int {
				if a.Index == b.Index {
					return 0
				}
				if a.Index > b.Index {
					return -1
				}
				return 1
			},
		), nil
	case lucytypes.ByDownloads:
		arr := make(
			[]tools.KeyValue[lucytypes.ProjectName, int],
			0,
			len(everything.Plugins),
		)
		for _, plugin := range everything.Plugins {
			id := lucytypes.ToProjectName(plugin.Meta.Id)
			download := 0
			for _, release := range plugin.Release.Releases {
				download += release.Asset.DownloadCount
			}
			arr = append(
				arr,
				tools.KeyValue[lucytypes.ProjectName, int]{
					Item: id, Index: download,
				},
			)
		}
		return tools.SortAndExtract(
			arr,
			func(a, b tools.KeyValue[lucytypes.ProjectName, int]) int {
				return b.Index - a.Index
			},
		), nil
	case lucytypes.ByNewest:
		arr := make(
			[]tools.KeyValue[lucytypes.ProjectName, time.Time],
			0,
			len(everything.Plugins),
		)
		for _, plugin := range everything.Plugins {
			id := lucytypes.ToProjectName(plugin.Meta.Id)
			timestamp := plugin.Release.Releases[0].CreatedAt
			arr = append(
				arr,
				tools.KeyValue[lucytypes.ProjectName, time.Time]{
					Item: id, Index: timestamp,
				},
			)
		}
		return tools.SortAndExtract(
			arr,
			func(a, b tools.KeyValue[lucytypes.ProjectName, time.Time]) int {
				if a.Index.Equal(b.Index) {
					return 0
				}
				if a.Index.After(b.Index) {
					return -1
				}
				return 1
			},
		), nil
	}

	return nil, fmt.Errorf("unknown index: %s", index)
}
