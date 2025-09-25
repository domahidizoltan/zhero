// Package schema manages the data blueprint.
package schema

type Schema struct {
	Name                string
	Identifier          string
	SecondaryIdentifier string
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
