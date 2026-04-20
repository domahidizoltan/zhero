// Package page manages the schema pages.
package page

import "github.com/domahidizoltan/zhero/pkg/paging"

const MaxSearchVals = 5

type (
	ReferenceMatch struct {
		Identifier          string
		SecondaryIdentifier string
	}

	Page struct {
		Route               string
		SchemaName          string
		Identifier          string
		SecondaryIdentifier string
		ListableData        map[string]any
		Data                map[string]any
		Meta                PageMeta
		References          []string
		IsEnabled           bool
		SearchVals          [MaxSearchVals]any
	}

	PageMeta struct {
		Title         string   `json:"title,omitempty"`
		Description   string   `json:"description,omitempty"`
		OGTitle       string   `json:"ogTitle,omitempty"`
		OGDescription string   `json:"ogDescription,omitempty"`
		Rating        string   `json:"rating,omitempty"`
		Robots        []string `json:"robots,omitempty"`
	}

	ListOptions struct {
		paging.PageOpts
		SecondaryIdentifierLike string
	}
)

func (pm PageMeta) ToMap() map[string]any {
	return map[string]any{
		"title":         pm.Title,
		"description":   pm.Description,
		"rating":        pm.Rating,
		"robots":        pm.Robots,
		"ogTitle":       pm.OGTitle,
		"ogDescription": pm.OGDescription,
	}
}
