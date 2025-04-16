package mcdr

import "strings"

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
	a, ok := everything.Authors[name]
	if !ok {
		return nil
	}
	return &a
}
