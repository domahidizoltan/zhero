// Package template is for collecting template files
package template

import (
	"github.com/domahidizoltan/zhero/pkg/handlebars"
)

const (
	tpl          = "template/"
	schemaorgTpl = tpl + "schemaorg/"
	pagesTpl     = tpl + "page/"
)

var (
	Index = handlebars.MustParse(tpl + "index.hbs")

	PageMain                     = handlebars.MustParse(pagesTpl + "main.hbs")
	PageList                     = handlebars.MustParse(pagesTpl + "list.hbs")
	SchemaorgSearch              = handlebars.MustParse(schemaorgTpl + "search.hbs")
	SchemaorgEdit                = handlebars.MustParse(schemaorgTpl + "edit.hbs")
	SchemaorgEditPropertyPartial = handlebars.MustParse(schemaorgTpl + "edit-property.partial.hbs")
)
