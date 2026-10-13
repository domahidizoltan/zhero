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

`page.data` (JSON-LD) is the source of truth for a page. Four sibling columns/tables are
**denormalized from it** and must all be rewritten on every page write
(`repository/page` `Insert`/`Update`):

| Projection | SQL | Derived from | Purpose |
|-|-|-|-|
| Canonical data | `page.data` | the field map + `@id`/`@type`/id-field injected on write | JSON-LD read model, SEO/API output |
| Listable data | `page.listable_data` | fields flagged `listable` | fast listing tiles without loading full data |
| Meta | `page.meta` | `PageMeta` (SEO/OpenGraph, `fieldMeta`) | header rendering, avoids SQL joins |
| References | `page.references` | long-form targets from `extractReferences` | cross-page links, index = `#refN` |
| Search | `page_search.col0..col4` | first 5 `searchable` field values (`MaxSearchVals`) | FTS5 full-text search |

Notes:
- Search has a **fixed width of 5 slots**. Searchable values are lifted into `col0..col4`;
  the secondary identifier is searchable via `page.secondary_identifier` but does **not**
  consume a slot (`dto.ToModel`, `f.Name != dto.SecondaryIdentifier`).
- If you add a projection or change the slot count, update every `page` write path plus the
  `page_search` insert/update/delete so all projections stay consistent.

## Reference lifecycle (4 stages)

Cross-page links round-trip through four representations: authoring token → `#refN` on persist
→ expand on load → resolve to `<a href>`. Full detail incl. functions: `controller/CLAUDE.md`.
Token formats: `wiki/dynamic-data-conventions.md` §4.

## Write atomicity

Multi-step domain writes run in one context-propagated `database.InTx`; repo writes return
`ErrTransactionNotFound` without one. Mechanics: `pkg/CLAUDE.md`; repo-side contract:
`repository/CLAUDE.md`.

## The `onlyEnabled` flag

Page read paths take `onlyEnabled bool`: `true` = public/published only, `false` = admin.
Details: `domain/CLAUDE.md`.
