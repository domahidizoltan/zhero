// Package schemaorg is a package to help process the RDF graph from Schema.org
package schemaorg

import (
	"slices"
	"strings"
	"sync"

	"github.com/deiu/rdf2go"
	"github.com/domahidizoltan/zhero/config"
	"github.com/domahidizoltan/zhero/pkg/collection"
	rdfPkg "github.com/domahidizoltan/zhero/pkg/rdf"
)

type schemaorg struct {
	graph         *rdfPkg.Graph
	unstableNodes map[string]struct{}
}

var once sync.Once

func New(cfg config.RdfConfig) (*schemaorg, error) {
	var (
		g             *rdfPkg.Graph
		err           error
		unstableNodes = make(map[string]struct{}, 1000)
	)

	once.Do(func() {
		g, err = rdfPkg.Init(cfg.File, cfg.Source, false)
		if err != nil {
			return
		}

		unstableNodes = getUnstableNodes(g)
	})

	if err != nil {
		return nil, err
	}
	return &schemaorg{
		graph:         g,
		unstableNodes: unstableNodes,
	}, nil
}

func (s *schemaorg) GetAllClasses() []string {
	triples := s.graph.All(nil, Type, Class)
	return s.prepareValues(triples, tripleSubject)
}

func (s *schemaorg) GetSubClassesOf(cls rdf2go.Term) []string {
	triples := s.graph.All(nil, SubClassOf, cls)
	return s.prepareValues(triples, tripleSubject)
}

func (s *schemaorg) GetSchemaClass(cls rdf2go.Term) *SchemaClass {
	desc := s.getDescription(cls)
	classTerms := s.getClassHierarchy(cls, nil)

	classes := s.filterUnstableValues(mapTerms(classTerms))
	props := s.getPropertiesOf(classes)
	allProps := []ClassProperty{}
	for _, p := range props {
		allProps = append(allProps, p...)
	}

	return &SchemaClass{
		Name:         getTermName(cls, schema),
		Description:  desc,
		CanonicalURL: cls.RawValue(),
		Properties:   allProps,
	}
}

func (s *schemaorg) getDescription(cls rdf2go.Term) string {
	if t := s.graph.One(cls, Comment, nil); t != nil {

		lit := t.Object.(*rdf2go.Literal)
		return lit.RawValue()
	}
	return ""
}

func (s *schemaorg) getExpectedType(cls rdf2go.Term) []string {
	types := s.graph.All(cls, RangeIncludes, nil)
	return s.prepareValues(types, tripleObject)
}

func (s *schemaorg) getClassHierarchy(cls rdf2go.Term, chainItems []rdf2go.Term) []rdf2go.Term {
	if t := s.graph.One(cls, SubClassOf, nil); t != nil {
		var obj rdf2go.Term = t.Object.(*rdf2go.Resource)
		chainItems = append(chainItems, s.getClassHierarchy(obj, chainItems)...)
	}
	return append(chainItems, cls)
}

func (s *schemaorg) getPropertiesOf(values []string) map[string][]ClassProperty {
	properties := make(map[string][]ClassProperty, len(values))
	for _, v := range values {
		sc := term(schema, v)
		props := s.graph.All(nil, DomainIncludes, sc)

		for _, p := range props {
			if _, found := s.unstableNodes[getTermName(p.Subject, schema)]; found {
				continue
			}
			properties[v] = append(properties[v], ClassProperty{
				Property:      getTermName(p.Subject, schema),
				CanonicalURL:  p.Subject.RawValue(),
				Description:   s.getDescription(p.Subject),
				ExpectedTypes: s.getExpectedType(p.Subject),
			})
		}
	}
	return properties
}

func getUnstableNodes(g *rdfPkg.Graph) map[string]struct{} {
	triples := g.All(nil, IsPartOf, attic)
	triples = append(triples, g.All(nil, IsPartOf, pending)...)

	values := make(map[string]struct{}, len(triples))
	for _, v := range mapNames(triples, tripleSubject) {
		values[v] = struct{}{}
	}

	return values
}

type tripleFn = func(*rdf2go.Triple) rdf2go.Term

var (
	tripleSubject = func(t *rdf2go.Triple) rdf2go.Term {
		return t.Subject
	}
	tripleObject = func(t *rdf2go.Triple) rdf2go.Term {
		return t.Object
	}
)

func (s *schemaorg) prepareValues(triples []*rdf2go.Triple, fn func(*rdf2go.Triple) rdf2go.Term) []string {
	res := s.filterUnstableValues(mapNames(triples, fn))
	slices.Sort(res)
	return res
}

func mapNames(triples []*rdf2go.Triple, fn tripleFn) []string {
	return slices.Collect(collection.MapValues(triples, func(t *rdf2go.Triple) string {
		return getTermName(fn(t), schema)
	}))
}

func mapTerms(terms []rdf2go.Term) []string {
	return slices.Collect(collection.MapValues(terms, func(t rdf2go.Term) string {
		return getTermName(t, schema)
	}))
}

func getTermName(term rdf2go.Term, ctx context) string {
	return strings.TrimPrefix(term.RawValue(), string(ctx))
}

func (s *schemaorg) filterUnstableValues(values []string) []string {
	return slices.Collect(collection.FilterValues(values, func(v string) bool {
		_, found := s.unstableNodes[v]
		return !found
	}))
}
