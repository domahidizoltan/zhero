// Package page manages the schema pages.
package page

import "github.com/domahidizoltan/zhero/pkg/paging"

const MaxSearchVals = 5

type (
	Page struct {
		Route               string
		SchemaName          string
		Identifier          string
		SecondaryIdentifier string
		Data                map[string]any
		IsEnabled           bool
		SearchVals          [MaxSearchVals]any
	}

	ListOptions struct {
		paging.PageOpts
		SecondaryIdentifierLike string
	}
)
