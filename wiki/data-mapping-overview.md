# Data mapping overview

This document is the cross-layer companion to the per-feature docs. It shows how data
moves through the layers and how the four derived projections of a page stay in sync.
Use it to orient before planning a feature that touches persistence or rendering.

## The layer chain

Data crosses five representations. Each hop is an explicit conversion, never a passthrough
of domain types into templates.

```
HTML form  ⇄  DTO            ⇄  domain model   ⇄  SQL row
(template)    (controller)      (domain)          (repository)
```

| Layer | Location | Page example | Schema example |
|-|-|-|-|
| HTML | `template/admin/**` (Handlebars) | `field-<name>`, `meta-*` inputs | `property-*` inputs |
| DTO | `controller/<domain>/dto.go` (private) | `adminpage.pageDto` / `fieldDto` | `adminschema.schemaDto` / `schemaPropDto` |
| Domain | `domain/<domain>/model.go` (public) | `page.Page` / `page.PageMeta` | `schema.SchemaMeta` / `schema.Property` |
| SQL | `repository/<domain>/repository.go` | `page`, `page_search` tables | `schema_meta`, `schema_meta_properties` |

Conversions are owned by the DTOs (`…From`, `ToModel`, `enhanceFromModel`, `ToMap`).
The domain layer never imports DB/HTTP types; repositories map domain↔SQL by hand
(`database/sql`, no ORM). Schema and page identifiers map **1:1** to their SQL columns
(no column swap — see `repository/schema/repository.go`).

## The source of truth and its four projections

`page.data` (JSON-LD) is the source of truth for a page. Three sibling columns and one FTS
table are **denormalized from it** and must all be rewritten on every page write
(`repository/page` `Insert`/`Update`):

| Projection | SQL | Derived from | Purpose |
|-|-|-|-|
| Canonical data | `page.data` | the field map + `@id`/`@type`/id-field injected on write | JSON-LD read model, SEO/API output |
| Listable data | `page.listable_data` | fields flagged `listable` | fast listing tiles without loading full data |
| Meta | `page.meta` | `PageMeta` (SEO/OpenGraph, `fieldMeta`) | header rendering, avoids SQL joins |
| Search | `page_search.col0..col4` | first 5 `searchable` field values (`MaxSearchVals`) | FTS5 full-text search |

Notes:
- Search has a **fixed width of 5 slots**. Searchable values are lifted into `col0..col4`;
  the secondary identifier is searchable via `page.secondary_identifier` but does **not**
  consume a slot (`dto.ToModel`, `f.Name != dto.SecondaryIdentifier`).
- If you add a projection or change the slot count, update every `page` write path plus the
  `page_search` insert/update/delete so all projections stay consistent.

## Reference lifecycle (4 stages)

Cross-page links round-trip through four representations; changing one regex/token format
breaks the others. Regexes live in `controller/controller.go`
(`RefPattern`, `ShortRefPattern`) and mirror `template/admin/page/page.js`.

1. **Author** — field embeds `#ZHERO#<schema_name>/<identifier>#{<json meta>}#`
   (`meta`: `linkText`, `altText`).
2. **Persist** — `adminpage/dto.go extractReferences` swaps each long ref for compact `#refN`
   and stores the targets in `Page.References` (index = N) → `page.references` column.
3. **Load** — `replaceShortRefs` expands `#refN` back to the long form for the editor.
4. **Resolve** — `dynamicpage.transformReferences` looks up targets via
   `pageSvc.ListEnabledRoutesByRef` and rewrites to `<a href>` (plain text when the target is
   disabled/missing or `disableAllLinks` is set, e.g. previews).

## Write atomicity

Multi-step domain writes run in one context-propagated transaction (`database.InTx`).
Repository write methods pull the active `*sql.Tx` off the context and return
`ErrTransactionNotFound` if absent — they never open their own. Example: `page.Create`
inserts the page **and** assigns its route inside a single `InTx`.

## The `onlyEnabled` flag

Page read paths take `onlyEnabled bool`: `true` for public-facing queries (published pages
only), `false` for admin/editing. Preserve this distinction when adding read paths.
