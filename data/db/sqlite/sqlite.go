// Package sqlite is for managing SQLite DB migration scripts.
package sqlite

import (
	_ "embed"
)

//go:embed 0000_init_schemas.sql
var schemametadataDdl string

var Scripts = []string{
	schemametadataDdl,
}
