// Package schemametadata manages the data blueprint.
package schemametadata

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
