// Package handlebars is for Handlebars helpers and variables
package handlebars

import (
	"regexp"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/russross/blackfriday"
)

var (
	replacer              = strings.NewReplacer("\\n", "<br/>", "<a ", "<a target='_blank' ")
	schemaLinkPlaceholder = regexp.MustCompile(`\[\[([^]]*)\]\]`)
	schemaLinkReplacer    = strings.NewReplacer("(/docs/", "(https://schema.org/docs/")

	helpers = map[string]any{
		"beautify": beautify,
	}
)

func InitHelpers() {
	raymond.RegisterHelpers(helpers)
}

func beautify(text string) string {
	text = schemaLinkReplacer.Replace(text)
	text = string(blackfriday.MarkdownBasic([]byte(text)))
	text = schemaLinkPlaceholder.ReplaceAllString(text, "<a href='https://schema.org/$1'>$1</a>")
	text = replacer.Replace(text)
	return text
}
