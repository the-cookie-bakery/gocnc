package gocnc

import "embed"

//go:embed entities/buildings/*.png
//go:embed entities/vehicles/*.png
//go:embed entities/turrets/*.png
//go:embed entities/objects/*.png
//go:embed actions/*.ogg
var EmbeddedFS embed.FS
