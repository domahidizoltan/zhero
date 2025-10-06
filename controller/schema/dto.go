package schema

import (
	"maps"

	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/domain/schemaorg"
	"github.com/domahidizoltan/zhero/pkg/collection"
)

type (
	schemaDto struct {
		IsLoaded            bool
		Name                string
		Description         string
		CanonicalURL        string
		Properties          []schemaPropDto
		Identifier          string `form:"identifier"` // binding:"required"`
		SecondaryIdentifier string `form:"secondary-identifier"`
	}
	schemaPropDto struct {
		NotUsed       bool
		Name          string
		CanonicalURL  string
		PossibleTypes []string
		Description   string
		Mandatory     bool
		Searchable    bool
		SelectedType  string
		Component     string
		Order         int
	}
)

func schemaDtoFrom(orgCls schemaorg.SchemaClass, domain *schema.SchemaMeta) schemaDto {
	props := make([]schemaPropDto, 0, len(orgCls.Properties))
	domainPropsByName := map[string]schema.Property{}

	if domain != nil {
		mapFn := func(p schema.Property) (string, schema.Property) {
			return p.Name, p
		}
		domainPropsByName = maps.Collect(collection.MapBy(domain.Properties, mapFn))
	}
	for _, prop := range orgCls.Properties {
		var domainProp *schema.Property
		if p, found := domainPropsByName[prop.Name]; found {
			domainProp = &p
		}
		props = append(props, schemaPropDtoFrom(prop, domainProp))
	}

	dto := schemaDto{
		Name:         orgCls.Name,
		Description:  orgCls.Description,
		CanonicalURL: orgCls.CanonicalURL,
		Properties:   props,
	}
	if domain != nil {
		dto.IsLoaded = true
		dto.Identifier = domain.Identifier
		dto.SecondaryIdentifier = domain.SecondaryIdentifier
	}
	return dto
}

func schemaPropDtoFrom(orgProp schemaorg.ClassProperty, domain *schema.Property) schemaPropDto {
	dto := schemaPropDto{
		Name:          orgProp.Name,
		CanonicalURL:  orgProp.CanonicalURL,
		PossibleTypes: orgProp.PossibleTypes,
		Description:   orgProp.Description,
		NotUsed:       true,
	}

	if domain != nil {
		dto.NotUsed = false
		dto.Mandatory = domain.Mandatory
		dto.Searchable = domain.Searchable
		dto.SelectedType = domain.Type
		dto.Component = domain.Component
		dto.Order = domain.Order
	}

	return dto
}
