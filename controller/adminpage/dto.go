package adminpage

import (
	"regexp"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	page_domain "github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/gin-gonic/gin"
)

type (
	pageDto struct {
		Route                    string
		SchemaName               string
		Fields                   []fieldDto
		Identifier               string
		SecondaryIdentifier      string
		SecondaryIdentifierValue any
		ListableData             map[string]any
		References               []string
		CreatedBy                string
		CreatedAt                time.Time
		UpdatedBy                string
		UpdatedAt                time.Time
		IsEnabled                bool
		Meta                     pageMeta
	}

	fieldDto struct {
		Name         string
		Order        uint
		IsMandatory  bool
		IsSearchable bool
		IsListable   bool
		Type         string
		Component    string
		InputType    bool
		Value        any
		References   []string
	}

	pageMeta struct {
		Title         string
		Description   string
		OGTitle       string
		OGDescription string
		Rating        string
		Robots        []string
	}
)

func PageDtoFrom(meta *schema.SchemaMeta) pageDto {
	if meta == nil {
		return pageDto{}
	}

	dto := pageDto{
		SchemaName:          meta.Name,
		Identifier:          meta.Identifier,
		SecondaryIdentifier: meta.SecondaryIdentifier,
		IsEnabled:           false,
	}

	dto.Fields = make([]fieldDto, 0, len(meta.Properties))
	for _, p := range meta.Properties {
		component := p.Component
		// Auto-determine component if empty or legacy "TODO"
		if component == "" || component == "TODO" {
			component = determineComponent(p.Type, p.Name)
		}
		dto.Fields = append(dto.Fields, fieldDto{
			Name:         p.Name,
			Order:        p.Order,
			IsMandatory:  p.Mandatory,
			IsSearchable: p.Searchable,
			IsListable:   p.Listable,
			Type:         p.Type,
			Component:    component,
			InputType:    slices.Contains([]string{"Color", "Email", "File", "Tel", "URL", "Number", "Date", "DateTime", "Time"}, component),
		})
	}

	return dto
}

func (dto *pageDto) EnhanceFromForm(c *gin.Context) {
	for i, f := range dto.Fields {
		dto.Fields[i].Value = c.PostForm("field-" + f.Name)
	}
	dto.IsEnabled = c.PostForm("is-enabled") == "on"
	dto.Route = c.PostForm("route")

	dto.Meta = pageMeta{
		Title:         c.PostForm("meta-title"),
		Description:   c.PostForm("meta-description"),
		OGTitle:       c.PostForm("meta-og-title"),
		OGDescription: c.PostForm("meta-og-description"),
	}

	if c.PostForm("meta-robots-noindex") == "on" {
		dto.Meta.Robots = append(dto.Meta.Robots, "noindex")
	}
	if c.PostForm("meta-robots-nofollow") == "on" {
		dto.Meta.Robots = append(dto.Meta.Robots, "nofollow")
	}

	if c.PostForm("meta-rating-adult") == "on" {
		dto.Meta.Rating = "adult"
	}
}

func (dto *pageDto) enhanceFromModel(p *page_domain.Page) {
	if p == nil {
		return
	}

	dto.IsEnabled = p.IsEnabled
	dto.Route = p.Route
	dto.Meta.FromModel(p.Meta)
	dto.SecondaryIdentifierValue = p.SecondaryIdentifier
	dto.ListableData = p.ListableData
	for i, f := range dto.Fields {
		if val, ok := p.Data[f.Name]; ok {
			dto.Fields[i].Value = val
		}
	}
}

func (dto *pageDto) ToModel() page_domain.Page {
	searchVals := [page_domain.MaxSearchVals]any{}
	data := make(map[string]any, len(dto.Fields))
	scIdx := 0
	listableData := make(map[string]any)
	references := make([]string, 0)

	for _, f := range dto.Fields {
		val := f.Value
		data[f.Name] = val

		if f.IsListable {
			listableData[f.Name] = val
		}

		if f.IsSearchable && f.Name != dto.SecondaryIdentifier && scIdx < 5 {
			searchVals[scIdx] = f.Value
			scIdx++
		}

		// Collect references from fields (will be extracted via extractReferences)
		references = append(references, f.References...)
	}

	// Deduplicate references
	refSet := make(map[string]struct{})
	for _, ref := range references {
		refSet[ref] = struct{}{}
	}
	uniqueRefs := make([]string, 0, len(refSet))
	for ref := range refSet {
		uniqueRefs = append(uniqueRefs, ref)
	}
	slices.Sort(uniqueRefs)

	return page_domain.Page{
		Route:               dto.Route,
		SchemaName:          dto.SchemaName,
		Identifier:          data[dto.Identifier].(string),
		SecondaryIdentifier: data[dto.SecondaryIdentifier].(string),
		Data:                data,
		IsEnabled:           dto.IsEnabled,
		SearchVals:          searchVals,
		Meta:                dto.Meta.ToModel(),
		ListableData:        listableData,
		References:          uniqueRefs,
	}
}

// TODO: extractReferences scans text fields for #ZHERO#... reference patterns
func (dto *pageDto) extractReferences() {
	refPattern := regexp.MustCompile(`#ZHERO#([^#]+)#\{([^}]*)\}#`)
	refSet := make(map[string]struct{})

	for i, f := range dto.Fields {
		if f.Value == nil {
			continue
		}
		strVal, ok := f.Value.(string)
		if !ok {
			continue
		}
		// Only extract from Text and TextArea fields
		if (f.Type == "Text" || f.Type == "TextArea") && strings.Contains(strVal, "#") {
			matches := refPattern.FindAllStringSubmatch(strVal, -1)
			fieldRefs := make([]string, 0, len(matches))
			for _, match := range matches {
				ref := match[1] // "Thing/123"
				fieldRefs = append(fieldRefs, ref)
				refSet[ref] = struct{}{}
			}
			dto.Fields[i].References = fieldRefs
		}
	}

	// Build global deduplicated references list
	dto.References = make([]string, 0, len(refSet))
	for ref := range refSet {
		dto.References = append(dto.References, ref)
	}
	slices.Sort(dto.References)
}

func (dto *pageDto) ToMap() map[string]any {
	if dto == nil {
		return nil
	}

	fields := make(map[string]any, len(dto.Fields))
	for _, f := range dto.Fields {
		fields[f.Name] = f.Value
	}
	return map[string]any{
		"route":                    dto.Route,
		"schemaName":               dto.SchemaName,
		"fields":                   fields,
		"identifier":               dto.Identifier,
		"secondaryIdentifier":      dto.SecondaryIdentifier,
		"secondaryIdentifierValue": dto.SecondaryIdentifierValue,
		"listableData":             dto.ListableData,
		"references":               dto.References,
		"isEnabled":                dto.IsEnabled,
		"meta":                     dto.Meta.ToMap(),
	}
}

func (dm *pageMeta) FromModel(pm page_domain.PageMeta) {
	*dm = pageMeta{
		Title:         pm.Title,
		Description:   pm.Description,
		OGTitle:       pm.OGTitle,
		OGDescription: pm.OGDescription,
		Rating:        pm.Rating,
		Robots:        pm.Robots,
	}
}

func (dm *pageMeta) ToModel() page_domain.PageMeta {
	if dm == nil {
		return page_domain.PageMeta{}
	}

	return page_domain.PageMeta{
		Title:         dm.Title,
		Description:   dm.Description,
		OGTitle:       dm.OGTitle,
		OGDescription: dm.OGDescription,
		Rating:        dm.Rating,
		Robots:        dm.Robots,
	}
}

func (dm *pageMeta) ToMap() map[string]any {
	if dm == nil {
		return nil
	}

	return map[string]any{
		"title":         dm.Title,
		"description":   dm.Description,
		"rating":        dm.Rating,
		"robots":        dm.Robots,
		"ogTitle":       dm.OGTitle,
		"ogDescription": dm.OGDescription,
	}
}
