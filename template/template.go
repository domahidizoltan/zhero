// Package template is for collecting template files
package template

import (
	_ "embed"

	"github.com/aymerick/raymond"
)

var (
	//go:embed index.css
	indexCSS string
	//go:embed index.hbs
	indexTpl string
	Index    = parse(indexTpl)

	//go:embed admin_index.js
	adminIndexJs string
	//go:embed admin_index.hbs
	adminIndexTpl string
	AdminIndex    = parse(adminIndexTpl)

	//go:embed page/page.js
	pagePageJs string

	//go:embed page/main.hbs
	pageMainTpl string
	PageMain    = parse(pageMainTpl)

	//go:embed page/list.hbs
	pageListTpl string
	PageList    = parse(pageListTpl)

	//go:embed page/edit.hbs
	pageEditTpl string
	PageEdit    = parse(pageEditTpl)

	//go:embed schemaorg/schemaorg.js
	schemaorgSchemaorgJs string

	//go:embed schemaorg/search.hbs
	schemaorgSearchTpl string
	SchemaorgSearch    = parse(schemaorgSearchTpl)

	//go:embed schemaorg/edit.hbs
	schemaorgEditTpl string
	SchemaorgEdit    = parse(schemaorgEditTpl)

	//go:embed schemaorg/edit-property.partial.hbs
	schemaorgEditPropertyPartialTpl string
	SchemaorgEditPropertyPartial    = parse(schemaorgEditPropertyPartialTpl)

	Assets = map[string]string{
		"/index.css":              indexCSS,
		"/admin_index.js":         adminIndexJs,
		"/page/page.js":           pagePageJs,
		"/schemaorg/schemaorg.js": schemaorgSchemaorgJs,
	}
)

func parse(templateStr string) *raymond.Template {
	tpl, err := raymond.Parse(templateStr)
	if err != nil {
		panic("failed to parse template: " + err.Error())
	}
	return tpl
}

func RegisterPartials() {
	SchemaorgEdit.RegisterPartialTemplate("editProperty", SchemaorgEditPropertyPartial)
}
