// Package page manages the schema pages.
package page

import "slices"

type Page struct {
	SchemaName          string
	Identifier          string
	SecondaryIdentifier string
	Fields              []Field
	IsEnabled           bool
}

type Field struct {
	Name         string
	Value        any
	Order        int
	SearchColumn string
}

func GetFieldIdxByName(fields []Field, name string) int {
	return slices.IndexFunc(fields, func(f Field) bool {
		return f.Name == name
	})
}
