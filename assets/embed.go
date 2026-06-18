package assets

import "embed"

//go:embed migrations/*.sql
var MigrationsFS embed.FS

//go:embed migrations/seeds/*.sql
var SeedsFS embed.FS
