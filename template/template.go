// Package template is for collecting template files
package template

import (
	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/pkg/handlebars"
)

const (
	tpl          = "template/"
	schemaorgTpl = tpl + "schemaorg/"
	pagesTpl     = tpl + "page/"
)

var Index, AdminIndex, PageMain, PageList, PageEdit, SchemaorgSearch, SchemaorgEdit, SchemaorgEditPropertyPartial *raymond.Template

func InitTemplates() {
	Index = handlebars.MustParse(tpl + "index.hbs")
	AdminIndex = handlebars.MustParse(tpl + "admin_index.hbs")

	PageMain = handlebars.MustParse(pagesTpl + "main.hbs")
	PageList = handlebars.MustParse(pagesTpl + "list.hbs")
	PageEdit = handlebars.MustParse(pagesTpl + "edit.hbs")
	SchemaorgSearch = handlebars.MustParse(schemaorgTpl + "search.hbs")
	SchemaorgEdit = handlebars.MustParse(schemaorgTpl + "edit.hbs")
	SchemaorgEditPropertyPartial = handlebars.MustParse(schemaorgTpl + "edit-property.partial.hbs")

	SchemaorgEdit.RegisterPartialTemplate("editProperty", SchemaorgEditPropertyPartial)
	handlebars.InitHelpers()
}
