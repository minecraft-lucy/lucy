package sources

import (
	"lucy/lucytypes"
	"lucy/remote"
	"lucy/remote/mcdr"
	"lucy/remote/modrinth"
)

var All = []remote.SourceHandler{
	modrinth.Self,
	mcdr.Self,
}

var Modrinth = modrinth.Self
var Mcdr = mcdr.Self

var Map = map[lucytypes.Source]remote.SourceHandler{
	lucytypes.Modrinth:      modrinth.Self,
	lucytypes.McdrCatalogue: mcdr.Self,
}
