package templates

import "embed"

//go:embed pages/* services/* agents/*
var FS embed.FS
