package db

import (
	"embed"
)

//go:embed migrations/*.sql

// MigrationsFS contains SQL migration files embedded into the binary.
var MigrationsFS embed.FS
