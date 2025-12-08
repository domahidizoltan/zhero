// Package jsonld provides helper functions to generate JSON-LD from Page model
package jsonld

import (
	"encoding/json"
	"fmt"

	"github.com/domahidizoltan/zhero/domain/page"
)

func FromPage(page page.Page) ([]byte, error) {
	jsonLD := make(map[string]any)
	jsonLD["@context"] = "https://schema.org/"
	jsonLD["@type"] = page.SchemaName

	for _, field := range page.Fields {
		if field.Name == page.Identifier {
			jsonLD["@id"] = field.Value
			continue
		}
		jsonLD[field.Name] = field.Value
	}

	data, err := json.MarshalIndent(jsonLD, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON-LD: %w", err)
	}

	return data, nil
}
