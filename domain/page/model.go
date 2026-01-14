// Package page manages the schema pages.
package page

var (
	SortDirAsc  SortDir = "asc"
	SortDirDesc SortDir = "desc"
)

const MaxSearchVals = 5

type (
	Page struct {
		SchemaName          string
		Identifier          string
		SecondaryIdentifier string
		Data                map[string]any
		IsEnabled           bool
		SearchVals          [MaxSearchVals]any
	}

	PagingMeta struct {
		TotalItems  uint
		PageSize    uint
		TotalPages  uint
		CurrentPage uint
	}

	SortDir string

	ListOptions struct {
		SecondaryIdentifierLike string
		SortBy                  string
		SortDir                 SortDir
		Page                    uint
		PageSize                uint
	}
)
