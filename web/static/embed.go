package static

import "embed"

//go:embed css/* fonts/* *
var StaticFS embed.FS
