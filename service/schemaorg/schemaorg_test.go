package schemaorg

import (
	"testing"

	"github.com/domahidizoltan/zhero/pkg/config"
	"github.com/stretchr/testify/assert"
)

const jsonldURL = "https://raw.githubusercontent.com/schemaorg/schemaorg/refs/heads/main/data/releases/29.2/schemaorg-all-https.jsonld"

func TestSchemaorg(t *testing.T) {
	so, err := New(config.RdfConfig{
		Source: jsonldURL,
		File:   "../../rdf_schema.jsonld",
	})
	assert.NoError(t, err)

	t.Run("gets_all_classes", func(t *testing.T) {
		res := so.GetAllClasses()
		assert.Greater(t, len(res), 100)
		assert.Contains(t, res, "CreativeWorkSeason")
		assert.Contains(t, res, "TVSeason")
		assert.NotContains(t, res, "PodcastSeason")
		assert.NotContains(t, res, "BioChemEntity")
	})

	t.Run("gets_subclasses", func(t *testing.T) {
		res := so.GetSubClassesOf(term(schema, "Thing"))
		assert.Len(t, res, 9)
		stableElements := []string{
			"Action", "CreativeWork", "Event", "Intangible",
			"MedicalEntity", "Organization", "Person", "Place", "Product",
		}
		assert.ElementsMatch(t, res, stableElements)
		assert.NotElementsMatch(t, res, []string{"StupidType", "BioChemEntity", "Taxon"})
	})

	t.Run("gets_schema_class", func(t *testing.T) {
		res := so.GetSchemaClass(term(schema, "LiveBlogPosting"))

		description := "A [[LiveBlogPosting]] is a [[BlogPosting]] intended to provide " +
			"a rolling textual coverage of an ongoing event through continuous updates."
		assert.Equal(t, "LiveBlogPosting", res.Name)
		assert.Equal(t, description, res.Description)
		assert.Greater(t, len(res.Properties), 100)

		assert.Contains(t, res.Properties, ClassProperty{
			Property:      "liveBlogUpdate",
			CanonicalURL:  "https://schema.org/liveBlogUpdate",
			Description:   "An update to the LiveBlog.",
			ExpectedTypes: []string{"BlogPosting"},
		})
		assert.Contains(t, res.Properties, ClassProperty{
			Property:      "audio",
			CanonicalURL:  "https://schema.org/audio",
			Description:   "An embedded audio object.",
			ExpectedTypes: []string{"AudioObject", "Clip", "MusicRecording"},
		})

		var props []string
		for _, p := range res.Properties {
			props = append(props, p.Property)
		}
		assert.NotContains(t, props, "usageInfo")
		assert.NotContains(t, props, "backStory")
	})
}
