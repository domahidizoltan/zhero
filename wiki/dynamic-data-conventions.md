# Dynamic data conventions (magic strings & jump points)

Zhero is hard to read because much of its data is **dynamic**: map keys are built by
concatenating strings, field names come from Schema.org property names stored in the DB,
and rows store JSON blobs instead of typed columns. This document is a lookup table of every
convention so you can jump straight to the code instead of tracing it. When you change a
value here, grep for it — the same literal usually appears in Go **and** in a `.hbs`/`.js`
template that must stay in sync.

## 1. HTML form field names (admin edit form → controller)

Form is `template/admin/page/edit.hbs`; parsed in `adminpage/dto.go EnhanceFromForm`
and `adminpage/controller.go`. Names are **built dynamically** from schema property names.

| Form key | Meaning | Read at |
|-|-|-|
| `field-<propName>` | one page field value (propName = Schema.org property) | `dto.go:113` `c.PostForm("field-"+f.Name)` |
| `alt-text-<propName>` | SEO alt text for that field → stored as `fieldMeta["<propName>:altText"]` | `dto.go:115` |
| `meta-title` / `meta-description` | SEO title/description | `dto.go:94` |
| `meta-og-title` / `meta-og-description` | OpenGraph | `dto.go:96` |
| `meta-robots-noindex` / `meta-robots-nofollow` | `on` → appended to `Robots[]` | `dto.go:101` |
| `meta-rating-adult` | `on` → `Rating="adult"` | `dto.go:108` |
| `is-enabled` | `on` → published | `dto.go:90` |
| `route` | custom URL slug | `dto.go:91` |
| `identifier` | hidden; present = update, absent = create | `controller.go:199` |
| `item-identifier` + `item-action` | list-row actions; action ∈ `enable`/`disable`/`delete` | `controller.go:145` |
| `property-*` (define-schema form) | schema property config | `template/admin/schemaorg/edit-property.partial.hbs` |

Checkbox convention: HTML checkboxes submit the literal string `"on"` when checked; every
`c.PostForm(...) == "on"` check relies on this.

## 2. Injected / reserved keys inside `page.Data` (the JSON-LD map)

`page.Data` is `map[string]any`. Most keys are Schema.org property names, but these are
**special** and injected/consumed by code, not by the user:

| Key | Set where | Meaning |
|-|-|-|
| `@id` | `repository/page/repository.go:82,124` (write) | JSON-LD id = the page identifier |
| `@type` | `repository/page:83,125` | JSON-LD type = schema name |
| `<idField>` | `repository/page:81,123` | duplicate of identifier under the schema's id property name (`idField` arg = `SchemaMeta.Identifier`) |
| `references` | `dynamicpage/controller.go:190` | injected map `#refN`/long-ref → resolved `<a href>` or text; consumed by the renderer |
| `fieldMeta` | `dynamicpage/controller.go:129` | `map[string]string` of SEO extras, e.g. `"<prop>:altText"` |
| `identifier` | list/render context maps | used to build upload file paths |
| `listableProperties` | `dynamicpage/controller.go:63` | the listable subset for list tiles |

Renderer **skips** these props (`pagerenderer/dynamicpagerenderer.go:30-35`):
`thumbnail`, the identifier property, and any key starting with `@`.

## 3. Compound / concatenated key formats

These string formats are constructed in one place and parsed in another — keep them identical.

| Format | Example | Built at | Parsed/used at |
|-|-|-|-|
| `<schema_name>/<identifier>` (pageKey / composite id) | `Article/01H...` | `page/service.go:48,69`, `adminpage/controller.go:240` | `route.page` column; `repository/page` `ListEnabledRoutesByRef` join |
| `<propName>:altText` (fieldMeta key) | `image:altText` | `adminpage/dto.go:116` | `pagerenderer:61`, `dynamicpage` |
| CSS class `lower(<SchemaName>-<propName>)` | `article-headline` | `pagerenderer:42` | emitted into rendered HTML |
| list CSS `list-item <SchemaName>` | `list-item Article` | `pagerenderer:84` | rendered HTML |
| upload path `<UploadsPath>/<SchemaName>/<identifier>/<file>` | `/zhero-content/uploads/Article/01H.../a.jpg` | `pagerenderer:69,72,118` | served as static files |
| hierarchy line `<marker*level><ClassName>` | `>>LocalBusiness` | `schemaorg/service.go` | `schema/service.go:73` counts `marker` (`>`) for depth |

## 4. Reference tokens (the 4-stage system)

Regexes are defined once in `controller/controller.go:15-16` and mirrored in
`template/admin/page/page.js`. See `data-mapping-overview.md` for the lifecycle; formats:

| Token | Regex var | Shape |
|-|-|-|
| Long (authoring/resolved) | `RefPattern` | `#ZHERO#<schema_name>/<identifier>#{<json meta>}#` |
| Short (persisted) | `ShortRefPattern` | `#refN` (N = index into `Page.References`) |

`<json meta>` keys: `linkText` (anchor text), `altText` (SEO alt). Parsed in
`dynamicpage/controller.go:170-186` (note it wraps with `{...}` before `json.Unmarshal`).

## 5. Component resolution (Schema.org type/name → HTML widget)

`adminpage/controller.go:305 determineComponent(propType, propName)` decides the widget.
**Name heuristics win over type.** Adding a component means editing all four sites (see
CLAUDE.md cross-cutting concern #4).

- Name contains (lowercased): `color`→Color, `email`→Email, `file`→File,
  `image`/`thumbnail`→Image, `phone`/`tel`→Tel.
- Else by type: `Boolean`→Checkbox, `Date`→Date, `DateTime`→DateTime,
  `Number`/`Integer`/`Float`→Number, `Quantity`/`Text`→TextInput, `URL`→URL, `Time`→Time.
- **Default (no match): `ReferenceSearch`** — i.e. complex Schema.org types become cross-page
  reference links. This is why unknown types render as reference pickers.
- `File`/`Image` get special HTML in the renderer; everything else is a `<p>`.
  Components in `InputType` set (Color/Email/Tel/Number/Date/DateTime/Time) render as native
  `<input type=...>` (`dto.go:82`).

## 6. Search slots (fixed width = 5)

`page_search` is FTS5 with `col0..col4`. `page.SearchVals [5]any` fills them in field order,
**skipping the secondary identifier** (it's searched via `page.secondary_identifier` instead):
`adminpage/dto.go:173` `f.IsSearchable && f.Name != dto.SecondaryIdentifier && scIdx < 5`.
Widening search = change the array size **and** every `page_search` query in `repository/page`.

## 7. Identifier vs. property-name (a frequent confusion)

- In **`schema_meta`**, `identifier` / `secondary_identifier` store **property names**
  (which field acts as the id), mapped 1:1 to `SchemaMeta.Identifier`/`SecondaryIdentifier`.
- In **`page`**, `identifier` / `secondary_identifier` store **values** (the actual ULID and
  the human-readable value pulled from that field).
- `dto.ToModel()` bridges them: `Identifier: data[dto.Identifier].(string)` — it reads the
  value out of the field named by the schema's id property (`adminpage/dto.go:182-183`).
- Page identifiers are **ULIDs** minted on insert (`repository/page/repository.go:74-80`,
  `oklog/ulid`); never editable after create.

## 8. Where dynamic values originate

- **Field/property names** come from the Schema.org RDF vocabulary
  (`rdf_schema.jsonld`, read via `domain/schemaorg` → `pkg/rdf`), narrowed to a user's
  `SchemaMeta.Properties`. They are not hardcoded anywhere.
- **JSON columns** (`data`, `meta`, `listable_data`, `references`) are marshaled/unmarshaled
  in `repository/page/repository.go`; the domain sees typed structs + `map[string]any`.
- **Sort inputs** (`SortBy`/`SortDir`) are concatenated into SQL (`repository/page:242`) —
  treat as trusted/validated upstream, never raw user input.
