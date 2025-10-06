package schema

import (
	"context"
	"strings"
	"sync"

	"github.com/deiu/rdf2go"
	"github.com/domahidizoltan/zhero/domain/schemaorg"
	"github.com/domahidizoltan/zhero/pkg/database"
)

type (
	schemaMetaRepo interface {
		Upsert(context.Context, SchemaMeta) error
		GetByClassName(context.Context, string) (*SchemaMeta, error)
	}

	schemaProvider interface {
		GetSchemaClassByName(cls string) *schemaorg.SchemaClass
		GetSubClassesHierarchyOf(cls rdf2go.Term, nestingLevelMarker string, currentLevel int) []string
	}
)

type Service struct {
	schemaMetaRepo schemaMetaRepo
	schemaProvider schemaProvider
	classHierarchy [][]string
}

func NewService(repo schemaMetaRepo, schemaProvider schemaProvider) Service {
	return Service{
		schemaMetaRepo: repo,
		schemaProvider: schemaProvider,
	}
}

func (s Service) SaveSchemaMeta(ctx context.Context, schema SchemaMeta) error {
	return database.InTx(ctx, func(ctx context.Context) error {
		return s.schemaMetaRepo.Upsert(ctx, schema)
	})
}

func (s Service) GetSchemaMetaByName(ctx context.Context, clsName string) (*SchemaMeta, error) {
	return s.schemaMetaRepo.GetByClassName(ctx, clsName)
}

func (s Service) GetSchemaClassByName(clsName string) *schemaorg.SchemaClass {
	cls := s.schemaProvider.GetSchemaClassByName(clsName)
	return cls
}

var clsHierarchyOnce sync.Once

func (s *Service) GetClassHierarchy() [][]string {
	clsHierarchyOnce.Do(func() {
		marker := ">"

		if len(s.classHierarchy) != 0 {
			return
		}

		lines := s.schemaProvider.GetSubClassesHierarchyOf(schemaorg.RootClass, marker, 0)
		parents := []string{lines[0]}
		s.classHierarchy = append(s.classHierarchy, []string{lines[0]})

		for _, l := range lines[1:] {
			level := strings.Count(l, marker)
			switch {
			case level == len(parents):
				parents = append(parents, l[level:])
			case level == len(parents)-1:
				parents[len(parents)-1] = l[level:]
			case level < len(parents)-1:
				parents = parents[:level]
				parents = append(parents, l[level:])
			}

			tmp := make([]string, len(parents))
			copy(tmp, parents)
			s.classHierarchy = append(s.classHierarchy, tmp)
		}
	})

	return s.classHierarchy
}
