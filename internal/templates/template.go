package templates

import "embed"

//go:embed pages/* services/* agents/* tests/* all:uca
var FS embed.FS
