package static

import "embed"

//go:embed css/* fonts/* favicon.ico
var StaticFS embed.FS
