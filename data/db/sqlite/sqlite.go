package sqlite

import (
	_ "embed"
)

//go:embed schemametadata.sql
var SchemametadataDdl string
