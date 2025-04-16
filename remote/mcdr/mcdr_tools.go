package mcdr

import (
	"lucy/lucytypes"
	"lucy/tools"
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

func sortBy(
	everything *everything,
	index lucytypes.SearchIndex,
) (res []lucytypes.ProjectName) {

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
				tools.KeyValue[lucytypes.ProjectName, float64]{id, relevance},
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
		)
	case lucytypes.ByDownloads:
		arr := make(
			[]tools.KeyValue[lucytypes.ProjectName, int],
			0,
			len(everything.Plugins),
		)
		for _, plugin := range everything.Plugins {
			download := 0
			for _, release := range plugin.Release.Releases {
				download += release.Asset.DownloadCount
			}
			arr = append(
				arr,
				tools.KeyValue[lucytypes.ProjectName, int]{
					lucytypes.ToProjectName(plugin.Meta.Id), download,
				},
			)
		}
		return tools.SortAndExtract(
			arr,
			func(a, b tools.KeyValue[lucytypes.ProjectName, int]) int {
				return b.Index - a.Index
			},
		)
	case lucytypes.ByNewest:
		arr := make(
			[]tools.KeyValue[lucytypes.ProjectName, time.Time],
			0,
			len(everything.Plugins),
		)
		for _, plugin := range everything.Plugins {
			timestamp := plugin.Release.Releases[0].CreatedAt
			arr = append(
				arr,
				tools.KeyValue[lucytypes.ProjectName, time.Time]{
					lucytypes.ToProjectName(plugin.Meta.Id), timestamp,
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
		)
	}

	return nil
}
