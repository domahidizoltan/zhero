package adminpage

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTestContext(form map[string]string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formEncode(form)))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c, w
}

func formEncode(m map[string]string) string {
	var b strings.Builder
	for k, v := range m {
		if b.Len() > 0 {
			b.WriteString("&")
		}
		b.WriteString(k + "=" + v)
	}
	return b.String()
}

func makeTestSchemaMeta() *schema.SchemaMeta {
	return &schema.SchemaMeta{
		Name:                "Article",
		Identifier:          "id",
		SecondaryIdentifier: "headline",
		Properties: []schema.Property{
			{Name: "id", Type: "Text", Mandatory: true, Searchable: false, Listable: false, Order: 0},
			{Name: "headline", Type: "Text", Mandatory: true, Searchable: true, Listable: true, Order: 1},
			{Name: "body", Type: "Text", Mandatory: true, Searchable: true, Listable: false, Order: 2},
			{Name: "author", Type: "Text", Mandatory: false, Searchable: true, Listable: true, Order: 3},
		},
	}
}

func TestExtractReferencesRoundTrip(t *testing.T) {
	meta := makeTestSchemaMeta()
	dto := PageDtoFrom(meta)

	// Build form with two distinct long refs + one duplicate across two fields
	ref1 := `#ZHERO#Article/01H123#{ "linkText":"First","altText":"alt1" }#`
	ref2 := `#ZHERO#Article/01H456#{ "linkText":"Second","altText":"alt2" }#`

	form := map[string]string{
		"is-enabled":     "on",
		"route":          "test-route",
		"field-headline": "Article " + ref1,
		"field-body":     "Body with " + ref2 + " and " + ref1, // duplicate ref1
		"field-author":   "Author name",
	}
	c, _ := makeTestContext(form)

	dto.EnhanceFromForm(c)

	// Call extractReferences (unexported, so we call it via the test in same package)
	dto.extractReferences()

	// Assert field values now contain #ref0 / #ref1
	headlineVal := dto.Fields[1].Value.(string)
	bodyVal := dto.Fields[2].Value.(string)

	assert.Contains(t, headlineVal, "#ref0")
	assert.Contains(t, bodyVal, "#ref1")

	// Duplicate long token maps to same index
	refCount := 0
	for _, f := range dto.Fields {
		if strings.Contains(f.Value.(string), "#ref") {
			refCount++
		}
	}
	assert.Equal(t, 2, refCount, "duplicate long ref should map to single #refN")

	// dto.References[idx] holds the exact long forms
	require.Len(t, dto.References, 2)
	assert.Equal(t, ref1, dto.References[0])
	assert.Equal(t, ref2, dto.References[1])

	// Round-trip: enhanceFromModel with page containing the references
	p := &page.Page{
		Data: map[string]any{
			"id":       "01H123",
			"headline": "Article " + "#ref0",
			"body":     "Body with " + "#ref1" + " and " + "#ref0",
			"author":   "Author name",
		},
		ListableData: map[string]any{
			"headline": "Article " + "#ref0",
			"author":   "Author name",
		},
		References: dto.References,
	}
	dto.enhanceFromModel(p)

	// replaceShortRefs should restore the originals exactly
	assert.Equal(t, "Article "+ref1, dto.Fields[1].Value)
	assert.Equal(t, "Body with "+ref2+" and "+ref1, dto.Fields[2].Value)
}

func TestToModelProjections(t *testing.T) {
	meta := &schema.SchemaMeta{
		Name:                "Article",
		Identifier:          "id",
		SecondaryIdentifier: "headline", // this is searchable but should be SKIPPED in SearchVals
		Properties: []schema.Property{
			{Name: "id", Type: "Text", Mandatory: true, Searchable: false, Listable: false, Order: 0},
			{Name: "headline", Type: "Text", Mandatory: true, Searchable: true, Listable: true, Order: 1}, // sec-id
			{Name: "body", Type: "Text", Mandatory: true, Searchable: true, Listable: false, Order: 2},
			{Name: "author", Type: "Text", Mandatory: false, Searchable: true, Listable: true, Order: 3},
			{Name: "datePublished", Type: "Date", Mandatory: false, Searchable: true, Listable: false, Order: 4},
			{Name: "category", Type: "Text", Mandatory: false, Searchable: true, Listable: false, Order: 5},
			{Name: "tags", Type: "Text", Mandatory: false, Searchable: true, Listable: false, Order: 6}, // 6th searchable -> should be dropped (cap 5)
		},
	}

	dto := PageDtoFrom(meta)
	form := map[string]string{
		"is-enabled":          "on",
		"route":               "test-route",
		"field-id":            "01H123",
		"field-headline":      "Test Article",
		"field-body":          "Body content",
		"field-author":        "John Doe",
		"field-datePublished": "2024-01-01",
		"field-category":      "Tech",
		"field-tags":          "go,test",
	}
	c, _ := makeTestContext(form)
	dto.EnhanceFromForm(c)

	m := dto.ToModel()

	// Data contains every field value
	assert.Equal(t, "01H123", m.Data["id"])
	assert.Equal(t, "Test Article", m.Data["headline"])
	assert.Equal(t, "Body content", m.Data["body"])
	assert.Equal(t, "John Doe", m.Data["author"])
	assert.Equal(t, "2024-01-01", m.Data["datePublished"])
	assert.Equal(t, "Tech", m.Data["category"])
	assert.Equal(t, "go,test", m.Data["tags"])

	// ListableData only contains listable fields
	assert.Contains(t, m.ListableData, "headline")
	assert.Contains(t, m.ListableData, "author")
	assert.NotContains(t, m.ListableData, "body")
	assert.NotContains(t, m.ListableData, "datePublished")
	assert.NotContains(t, m.ListableData, "category")
	assert.NotContains(t, m.ListableData, "tags")

	// SearchVals fills slots 0..4 in field order, SKIPS secondary identifier (headline)
	// Searchable fields in order: headline(sec-id, SKIP), body(0), author(1), datePublished(2), category(3), tags(4, DROPPED cap 5)
	assert.Equal(t, "Body content", m.SearchVals[0])
	assert.Equal(t, "John Doe", m.SearchVals[1])
	assert.Equal(t, "2024-01-01", m.SearchVals[2])
	assert.Equal(t, "Tech", m.SearchVals[3])
	// tags is the 6th searchable non-sec-id -> dropped
	assert.Equal(t, "go,test", m.SearchVals[4])

	// Identifier read from data[id-property]
	assert.Equal(t, "01H123", m.Identifier)
	assert.Equal(t, "Test Article", m.SecondaryIdentifier)

	// Meta round-trips
	assert.Equal(t, "test-route", m.Route)
	assert.True(t, m.IsEnabled)
}

func TestDetermineComponent(t *testing.T) {
	tests := []struct {
		name     string
		propType string
		propName string
		want     string
	}{
		// Name heuristics win over type
		{"color name wins", "URL", "backgroundColor", "Color"},
		{"email name wins", "Text", "userEmail", "Email"},
		{"file name wins", "Text", "uploadFile", "File"},
		{"image name wins", "Text", "heroImage", "Image"},
		{"thumbnail name wins", "Text", "thumbnail", "Image"},
		{"phone name wins", "Text", "phoneNumber", "Tel"},
		{"tel name wins", "Text", "telNumber", "Tel"},
		// Type fallbacks
		{"Boolean type", "Boolean", "active", "Checkbox"},
		{"Date type", "Date", "createdAt", "Date"},
		{"DateTime type", "DateTime", "updatedAt", "DateTime"},
		{"Number type", "Number", "count", "Number"},
		{"Integer type", "Integer", "quantity", "Number"},
		{"Float type", "Float", "price", "Number"},
		{"Quantity type", "Quantity", "amount", "TextInput"},
		{"Text type", "Text", "description", "TextInput"},
		{"URL type (no name match)", "URL", "website", "URL"},
		{"Time type", "Time", "startTime", "Time"},
		// Default: unknown type -> ReferenceSearch
		{"unknown type defaults to ReferenceSearch", "GeoCoordinates", "location", "ReferenceSearch"},
		// Name wins over URL type
		{"color name wins over URL", "URL", "backgroundColor", "Color"},
		{"email name wins over URL", "URL", "contactEmail", "Email"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineComponent(tt.propType, tt.propName)
			assert.Equal(t, tt.want, got)
		})
	}
}
