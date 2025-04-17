package mcdr

import (
	"lucy/remote"
	"strings"
)

// getPlugin is able to tolerate the interchange of "-" and "_"
func getPlugin(id string) *plugin {
	everything, err := getEverything()
	if err != nil {
		return nil
	}
	p, ok := everything.Plugins[id]
	if !ok {
		id = strings.Replace(id, "-", "_", -1)
		p, ok = everything.Plugins[id]
		if !ok {
			return nil
		}
	}
	return &p
}

func getAuthor(name string) *author {
	everything, err := getEverything()
	if err != nil {
		return nil
	}
	a, ok := everything.Authors.Authors[name]
	if !ok {
		return nil
	}
	return &a
}

func getRelease(plugin *plugin, version string) (release *release, err error) {
	for _, release := range plugin.Release.Releases {
		if release.Meta.Version == version {
			return &release, nil
		}
	}
	return nil, remote.FormatError(
		remote.ErrorNoVersion,
		plugin.Meta.Name,
		version,
	)
}
