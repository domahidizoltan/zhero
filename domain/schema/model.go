// Package schema manages the data blueprint.
package schema

type SchemaMeta struct {
	Name                string
	Identifier          string `form:"identifier"` // binding:"required"`
	SecondaryIdentifier string `form:"secondary-identifier"`
	Properties          []Property
}

type Property struct {
	Name       string
	Mandatory  bool
	Searchable bool
	Type       string
	Component  string
	Order      int
}
