// Package templates is for collecting template files
package templates

import (
	"github.com/domahidizoltan/zhero/pkg/handlebars"
)

const (
	tpl          = "templates/"
	schemaorgTpl = tpl + "schemaorg/"
)

var (
	Index = handlebars.MustParse(tpl + "index.hbs")

	SchemaorgSearch              = handlebars.MustParse(schemaorgTpl + "search.hbs")
	SchemaorgEdit                = handlebars.MustParse(schemaorgTpl + "edit.hbs")
	SchemaorgEditPropertyPartial = handlebars.MustParse(schemaorgTpl + "edit-property.partial.hbs")
)
