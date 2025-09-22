package schemametadata

type Schema struct {
	Name                string
	Identifier          string
	SecondaryIdentifier string
	Properties          []Property
}

// Property represents a schema property to be saved
type Property struct {
	Name       string
	Mandatory  bool
	Searchable bool
	Type       string
	Component  string
	Order      int
}
