// Package template is for collecting template files
package template

import (
	"embed"

	"github.com/aymerick/raymond"
)

const admin = "admin/"

var (
	//go:embed admin/*
	//go:embed *.css *.hbs
	//go:embed paging/*
	templates embed.FS

	AdminIndex           = mustParse(admin + "index.hbs")
	AdminPageMain        = mustParse(admin + "page/main.hbs")
	AdminPageList        = mustParse(admin + "page/list.hbs")
	AdminPageEdit        = mustParse(admin + "page/edit.hbs")
	AdminSchemaorgSearch = mustParse(admin + "schemaorg/search.hbs")
	AdminSchemaorgEdit   = mustParse(admin + "schemaorg/edit.hbs")

	AdminSchemaorgEditPropertyPartial = mustParse(admin + "schemaorg/edit-property.partial.hbs")

	AdminAssets = map[string][]byte{
		"/index.js":               mustLoad(admin + "index.js"),
		"/page/page.js":           mustLoad(admin + "page/page.js"),
		"/schemaorg/schemaorg.js": mustLoad(admin + "schemaorg/schemaorg.js"),
	}

	Index             = mustParse("index.hbs")
	PageNotFound      = mustParse("page_not_found.hbs")
	PaginationPartial = mustParse("paging/pagination.partial.hbs")

	Assets = map[string][]byte{
		"/index.css": mustLoad("index.css"),
	}
)

func mustParse(filename string) *raymond.Template {
	data := mustLoad(filename)
	return raymond.MustParse(string(data))
}

func mustLoad(filename string) []byte {
	data, err := templates.ReadFile(filename)
	if err != nil {
		panic("failed to read template file: " + err.Error())
	}
	return data
}

func RegisterPartials() {
	AdminSchemaorgEdit.RegisterPartialTemplate("editProperty", AdminSchemaorgEditPropertyPartial)
	raymond.RegisterPartialTemplate("pagination", PaginationPartial)
}
