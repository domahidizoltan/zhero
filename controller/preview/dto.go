package preview

type PagePreviewRequest struct {
	SchemaIdentifier    string            `json:"schemaIdentifier" form:"schemaIdentifier" binding:"required"`
	Title               string            `json:"title" form:"title" binding:"required"`
	Identifier          string            `json:"identifier" form:"identifier" binding:"required"`
	SecondaryIdentifier string            `json:"secondaryIdentifier" form:"secondaryIdentifier"`
	Slug                string            `json:"slug" form:"slug"`
	Fields              map[string]string `json:"fields" form:"fields"`
}
