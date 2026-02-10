package source

import (
	"lucy/remote"
	"lucy/remote/mcdr"
	"lucy/remote/modrinth"
	"lucy/types"
)

// All is currently hardcoded, but in the future, this could be made customizable
var All = []remote.SourceHandler{
	modrinth.Self,
	mcdr.Self,
}

var (
	Modrinth = modrinth.Self
	Mcdr     = mcdr.Self
)

var Map = map[types.Source]remote.SourceHandler{
	types.Modrinth:      modrinth.Self,
	types.McdrCatalogue: mcdr.Self,
}
