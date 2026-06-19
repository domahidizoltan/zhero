# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Scope

This `data/` directory is the persistence layer of the larger `zhero` Go project. It has no `go.mod` of its own — it is a sub-package of the parent module, so `go` commands (build/vet/test) run from the parent module root, not here.

## Architecture

The package implements a **schema-driven CMS** persisted in SQLite. Content types are defined by data (schemas), not by Go structs, so adding new content types is a data operation rather than a code change.

### Migration mechanism (`db/sqlite/sqlite.go`)

SQL files are compiled into the binary with `//go:embed` and exposed as an ordered `Scripts []string` slice. The slice order is the execution order — the consumer (in the parent module) runs each script in sequence. There is no rollback/down-migration concept; every script is written to be idempotent (`CREATE TABLE IF NOT EXISTS`, `CREATE INDEX IF NOT EXISTS`, etc.).

### Adding a migration

1. Create `db/sqlite/NNNN_description.sql` using the next zero-padded ordinal (existing: `0000_init_schemas.sql`).
2. Add a matching `//go:embed NNNN_description.sql` directive and `var` in `sqlite.go`.
3. Append the new var to the `Scripts` slice **in order** — position determines when it runs.
4. Keep statements idempotent; existing databases re-run the full script list.

### Data model (`db/sqlite/0000_init_schemas.sql`)

The tables form a content-modeling system rather than a fixed domain schema:

- **`schema_meta`** — content type definitions, keyed by unique `name`, with an `identifier` / `secondary_identifier` pair used to address instances.
- **`schema_meta_properties`** — per-schema field definitions. Flags drive behavior elsewhere: `mandatory` (validation), `searchable` (indexed into `page_search`), `listable` (denormalized into `page.listable_data`). `type` + `component` describe how a field is edited/rendered (HTML components); `order` controls field ordering.
- **`page`** — content instances. `data` holds the instance as **JSON-LD**; `listable_data`, `meta`, and `references` are denormalized JSON projections derived from `data`. `enabled` gates public visibility. Indexed by `(schema_name, identifier)`.
- **`page_search`** — FTS5 virtual table for full-text search. Columns are generic (`col0`–`col4`), so a schema's `searchable` properties are mapped positionally into these five slots rather than by name.
- **`route`** — custom URL routing. Maps a `route` string to a `page` with a `version`, enabling redirect middleware and versioned addressing.

When changing how content is stored, keep the four projections consistent: a write to `page.data` (JSON-LD) must be reflected into `listable_data`, the `page_search` columns, and `route` as applicable — these are derived views of the same content, not independent sources of truth.
