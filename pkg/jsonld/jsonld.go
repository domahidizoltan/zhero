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

	// Use the Data map directly
	for key, value := range page.Data {
		if key == page.Identifier {
			jsonLD["@id"] = value
			continue
		}
		jsonLD[key] = value
	}

	data, err := json.MarshalIndent(jsonLD, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON-LD: %w", err)
	}

	return data, nil
}
