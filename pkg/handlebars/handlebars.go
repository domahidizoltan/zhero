// Package handlebars is for Handlebars helpers and variables
package handlebars

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/pkg/logging"
	"github.com/rs/zerolog/log"
	"github.com/russross/blackfriday"
)

var (
	replacer              = strings.NewReplacer("\\n", "<br/>", "<a ", "<a target='_blank' ")
	schemaLinkPlaceholder = regexp.MustCompile(`\[\[([^]]*)\]\]`)
	schemaLinkReplacer    = strings.NewReplacer("(/docs/", "(https://schema.org/docs/")

	helpers = map[string]any{
		"beautify":      beautify,
		"use":           use,
		"compareAndUse": compareAndUse,
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

func use(use string, enabled, condition bool) string {
	if enabled && condition {
		return use
	}
	return ""
}

func compareAndUse(use string, enabled bool, optionVal, comparisonVal any) string {
	if enabled && reflect.DeepEqual(optionVal, comparisonVal) {
		return use
	}
	return ""
}

var once sync.Once

func MustParse(filePath string) *raymond.Template {
	once.Do(func() {
		logging.ConfigureLogging(nil)
	})

	tpl, err := raymond.ParseFile(filePath)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", filePath).
			Msg("failed to parse template")
	}
	return tpl
}
