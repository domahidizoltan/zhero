// Package template is for collecting template files
package template

import (
	"embed"

	"github.com/aymerick/raymond"
)

var (
	//go:embed *.css *.js *.hbs
	//go:embed page/* schemaorg/*
	templates embed.FS

	Index                        = mustParse("index.hbs")
	AdminIndex                   = mustParse("admin_index.hbs")
	PageMain                     = mustParse("page/main.hbs")
	PageList                     = mustParse("page/list.hbs")
	PageEdit                     = mustParse("page/edit.hbs")
	SchemaorgSearch              = mustParse("schemaorg/search.hbs")
	SchemaorgEdit                = mustParse("schemaorg/edit.hbs")
	SchemaorgEditPropertyPartial = mustParse("schemaorg/edit-property.partial.hbs")

	Assets = map[string][]byte{
		"/index.css":              mustLoad("index.css"),
		"/admin_index.js":         mustLoad("admin_index.js"),
		"/page/page.js":           mustLoad("page/page.js"),
		"/schemaorg/schemaorg.js": mustLoad("schemaorg/schemaorg.js"),
	}
)

func mustLoad(filename string) []byte {
	data, err := templates.ReadFile(filename)
	if err != nil {
		panic("failed to read template file: " + err.Error())
	}
	return data
}

func mustParse(filename string) *raymond.Template {
	data := mustLoad(filename)
	return raymond.MustParse(string(data))
}

func RegisterPartials() {
	SchemaorgEdit.RegisterPartialTemplate("editProperty", SchemaorgEditPropertyPartial)
}
