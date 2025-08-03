package dist

import (
	"embed"
)

var (
	//go:embed contents/*
	Contents embed.FS
)
