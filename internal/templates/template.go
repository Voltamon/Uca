package templates

import "embed"

//go:embed pages/* services/* agents/* uca/*
var FS embed.FS
