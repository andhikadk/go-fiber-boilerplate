package assets

import "embed"

//go:embed migrations/*
var MigrationsFS embed.FS
