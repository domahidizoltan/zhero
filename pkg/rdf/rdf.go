// Package rdf loads and queries RDF graph
package rdf

import (
	"os"
	"sync"

	"github.com/deiu/rdf2go"
	"github.com/domahidizoltan/zhero/pkg/file"
)

type Graph struct {
	graph *rdf2go.Graph
}

var once sync.Once

func Init(filePath, downloadURL string, overwrite bool) (*Graph, error) {
	// log
	var err error
	once.Do(func() {
		if err = file.DownloadToPath(filePath, downloadURL, overwrite); err != nil {
			return
		}
	})

	if err != nil {
		return nil, err
	}

	r, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	g := rdf2go.NewGraph("")
	if err := g.Parse(r, "application/ld+json"); err != nil {
		return nil, err
	}

	return &Graph{
		graph: g,
	}, nil
}

func (s *Graph) One(subject, predicate, object rdf2go.Term) *rdf2go.Triple {
	return s.graph.One(subject, predicate, object)
}

func (s *Graph) All(subject, predicate, object rdf2go.Term) []*rdf2go.Triple {
	return s.graph.All(subject, predicate, object)
}

func (s *Graph) Remove(triple *rdf2go.Triple) {
	s.graph.Remove(triple)
}
