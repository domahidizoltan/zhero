# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Scope

This directory is the **data-persistence (repository) layer** of the `github.com/domahidizoltan/zhero`
Go module — a schema-driven CMS. The module root (`go.mod`, `main.go`) and the packages these
files depend on (`pkg/database`, `pkg/paging`, `pkg/collection`, `domain/*`) live **outside this
directory**, one or more levels up. Import paths here are `github.com/domahidizoltan/zhero/repository/{page,route,schema}`.

## Commands

Go resolves the module by walking up to the parent `go.mod`, so these run from this directory:

```sh
go build ./...                 # build the three repository packages
go vet ./...
staticcheck ./...              # linter in use (emits SA#### diagnostics)
go test ./...                  # no test files exist here yet
go test -run TestName ./page/  # run a single test in a package
```

Cross-platform driver selection is done with build tags (see below). To exercise the Android
path: `go build -tags android ./...` (requires cgo — `CGO_ENABLED=1` and a C toolchain).

## Architecture

Three independent repositories, each a `Repository` struct constructed with `NewRepo(db *sql.DB, ...)`
and backed by SQLite. They model a CMS where users define **schemas** (content types), create
**pages** (content instances) against them, and expose pages at versioned **routes**.

- **`schema/`** — the data blueprint. `schema_meta` holds a content type's `name`, `identifier`
  field, and `secondary_identifier` field; `schema_meta_properties` holds its field definitions
  (`mandatory`, `searchable`, `listable`, `type`, `component`, `order`). `Upsert` replaces all
  properties transactionally (delete-then-bulk-insert).
- **`page/`** — content instances for a schema. Largest/most central package. Stores JSON columns
  (`data`, `meta`, `listable_data`, `references`) plus an `enabled` flag, and keeps a parallel
  `page_search` row in sync for full-text search.
- **`route/`** — maps a URL `route` to a `page`, **versioned**. `Create` auto-increments version via
  `MAX(version)+1` per page; `GetByRoute` resolves a URL, `GetLatestVersion` finds newest route for a page.

### Critical convention: transactions live in the context

Write methods (`Insert`, `Update`, `Delete`, `Enable`, `Upsert`, `route.Create`) retrieve the active
transaction with `database.GetTx(ctx)` and **return `database.ErrTransactionNotFound` if it is nil** —
they never open their own transaction. Callers must wrap writes in a transaction managed by
`pkg/database`. Read methods (`Get*`, `List*`) use `r.db` directly and do not require a transaction.

### Other patterns worth knowing

- **SQLite driver via build tags.** `*.go` (`//go:build !android`) uses pure-Go `modernc.org/sqlite`
  (no cgo); `*_android.go` (`//go:build android`) uses cgo-based `github.com/mattn/go-sqlite3`. These
  blank-import files only register the driver; keep both in sync when changing drivers.
- **IDs and JSON-LD.** New pages get a ULID identifier (`oklog/ulid`). On write, `page.Data` is
  injected with the schema's id field, plus `@id` (= identifier) and `@type` (= schema name).
- **Full-text search has exactly 5 columns.** `page_search` is `col0..col4`, populated from the fixed
  `page.SearchVals[0..4]` array. Every page Insert/Update/Delete mutates `page` and `page_search`
  together; adding a searchable column means changing this fixed width across the queries.
- **Pagination.** `page.List` builds queries dynamically and returns `paging.Meta`. `SortBy`/`SortDir`
  are concatenated into SQL (defaulting to `identifier`); only bound parameters are user-safe, so
  treat sort inputs as trusted/validated upstream. `SecondaryIdentifierLike` does a `LIKE %...%` filter.
- **SQL reserved words are quoted:** `"references"`, `[type]`, `[order]`. Preserve the quoting.
- **References / route resolution.** `references` is JSON on a page. `ListEnabledRoutesByRef` takes
  `schema_name/identifier` strings, expands a single `?` placeholder to an `IN (?, ?, ...)` list, and
  returns each ref's latest route + enabled state (window function `ROW_NUMBER() ... ORDER BY version DESC`).

### Inferred SQLite schema (DDL lives in the parent module)

```
schema_meta(name PK, identifier, secondary_identifier)
schema_meta_properties(schema_name, name, mandatory, searchable, listable, type, component, order)
page(schema_name, identifier, secondary_identifier, listable_data, data, meta, references, enabled)
page_search(schema_name, identifier, col0, col1, col2, col3, col4)   -- FTS-style mirror of page
route(route, page, version)                                          -- (page, version) is the natural key
```
