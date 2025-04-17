package mcdr

import (
	"fmt"
	"strings"
	"time"

	"lucy/lucytypes"
	"lucy/tools"

	"github.com/sahilm/fuzzy"
)

func search(
	obj *queriedEverything,
) ([]lucytypes.ProjectName, error) {
	matches, err := match(&obj.Everything, obj.Query)
	if err != nil {
		return nil, err
	}
	res := make([]lucytypes.ProjectName, 0, len(matches))
	for _, match := range matches {
		res = append(res, lucytypes.ProjectName(match.Str))
	}
	err = sortBy(res, obj.IndexBy)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// match is a helper function to for self.Search, as the MCDR API just gives
// the whole catalogue in a single file, we need to filter the results by
// query.
func match(
	everything *everything,
	query string,
) (matches fuzzy.Matches, err error) {
	var ids = make([]string, 0, len(everything.Plugins))
	for id := range everything.Plugins {
		ids = append(ids, id)
	}
	matches = fuzzy.Find(query, ids)
	return matches, nil
}

// sortBy is a helper function to sort the plugins by the given index.
//
// This is in-place sorting, so the original slice is modified.
func sortBy(
	res []lucytypes.ProjectName,
	index lucytypes.SearchIndex,
) (err error) {
	switch index {
	case lucytypes.ByRelevance:
		// Do nothing to res. Since if res is processed by match(), it is
		// already in relevance order.
		return nil
	case lucytypes.ByDownloads:
		iarr := make(
			[]tools.KeyValue[lucytypes.ProjectName, int], 0, len(res),
		)
		for _, name := range res {
			download := 0
			plugin := getPlugin(projectNameToMcdrId(name))
			if plugin == nil {
				continue
			}
			for _, release := range plugin.Release.Releases {
				download += release.Asset.DownloadCount
			}
			iarr = append(
				iarr,
				tools.KeyValue[lucytypes.ProjectName, int]{
					Item: name, Index: download,
				},
			)
		}
		res = tools.SortAndExtract(
			iarr,
			func(a, b tools.KeyValue[lucytypes.ProjectName, int]) int {
				return b.Index - a.Index
			},
		)
		return nil
	case lucytypes.ByNewest:
		iarr := make(
			[]tools.KeyValue[lucytypes.ProjectName, time.Time], 0, len(res),
		)
		for _, name := range res {
			plugin := getPlugin(projectNameToMcdrId(name))
			if plugin == nil {
				continue
			}
			timestamp := time.Unix(0, 0)
			if len(plugin.Release.Releases) > 0 {
				timestamp = plugin.Release.Releases[0].CreatedAt
			}
			iarr = append(
				iarr,
				tools.KeyValue[lucytypes.ProjectName, time.Time]{
					Item: name, Index: timestamp,
				},
			)
		}
		res = tools.SortAndExtract(
			iarr,
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
		return nil
	}

	return fmt.Errorf("unknown index: %s", index)
}

func projectNameToMcdrId(
	name lucytypes.ProjectName,
) (id string) {
	return strings.Replace(name.String(), "-", "_", -1)
}
