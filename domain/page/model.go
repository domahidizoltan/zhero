// Package page manages the schema pages.
package page

import "slices"

var (
	SortDirAsc  SortDir = "asc"
	SortDirDesc SortDir = "desc"
)

type (
	Page struct {
		SchemaName          string
		Identifier          string
		SecondaryIdentifier string
		Fields              []Field
		IsEnabled           bool
	}

	Field struct {
		Name         string
		Value        any
		Order        uint
		SearchColumn string
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

func GetFieldIdxByName(fields []Field, name string) int {
	return slices.IndexFunc(fields, func(f Field) bool {
		return f.Name == name
	})
}
