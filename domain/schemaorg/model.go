// Package schemaorg is a package to help process the RDF graph from Schema.org
package schemaorg

import (
	"github.com/deiu/rdf2go"
)

var (
	Class          = term(rdfs, "Class")
	Comment        = term(rdfs, "comment")
	IsPartOf       = term(schema, "isPartOf")
	DomainIncludes = term(schema, "domainIncludes")
	RangeIncludes  = term(schema, "rangeIncludes")
	SubClassOf     = term(rdfs, "subClassOf")
	Type           = term(rdf, "type")

	attic   = rdf2go.NewResource("https://attic.schema.org")
	pending = rdf2go.NewResource("https://pending.schema.org")

	RootClass = term(schema, "Thing")
)

type context string

const (
	// brick    context = "https://brickschema.org/schema/Brick#"
	// csvw     context = "http://www.w3.org/ns/csvw#"
	// dc       context = "http://purl.org/dc/elements/1.1/"
	// dcam     context = "http://purl.org/dc/dcam/"
	// dcat     context = "http://www.w3.org/ns/dcat#"
	// dcmitype context = "http://purl.org/dc/dcmitype/"
	// dcterms  context = "http://purl.org/dc/terms/"
	// doap     context = "http://usefulinc.com/ns/doap#"
	// foaf     context = "http://xmlns.com/foaf/0.1/"
	// odrl     context = "http://www.w3.org/ns/odrl/2/"
	// org      context = "http://www.w3.org/ns/org#"
	// owl      context = "http://www.w3.org/2002/07/owl#"
	// prof     context = "http://www.w3.org/ns/dx/prof/"
	// prov     context = "http://www.w3.org/ns/prov#"
	// qb       context = "http://purl.org/linked-data/cube#"
	rdf    context = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	rdfs   context = "http://www.w3.org/2000/01/rdf-schema#"
	schema context = "https://schema.org/"
	// sh       context = "http://www.w3.org/ns/shacl#"
	// skos     context = "http://www.w3.org/2004/02/skos/core#"
	// sosa     context = "http://www.w3.org/ns/sosa/"
	// ssn      context = "http://www.w3.org/ns/ssn/"
	// time     context = "http://www.w3.org/2006/time#"
	// vann     context = "http://purl.org/vocab/vann/"
	// void     context = "http://rdfs.org/ns/void#"
	// xsd      context = "http://www.w3.org/2001/XMLSchema#"
)

func term(ctx context, value string) rdf2go.Term {
	return rdf2go.NewResource(string(ctx) + value)
}

type (
	SchemaClass struct {
		Name         string
		Description  string
		CanonicalURL string
		Properties   []ClassProperty
	}

	ClassProperty struct {
		Property      string
		CanonicalURL  string
		ExpectedTypes []string
		Description   string
	}
)
