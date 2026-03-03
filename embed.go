package embed

import "embed"

//go:embed static
var StaticFiles embed.FS

//go:embed template
var Templates embed.FS
