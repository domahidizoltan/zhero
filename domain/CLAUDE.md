# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Scope

This directory is the **domain layer** of the larger `github.com/domahidizoltan/zhero` Go module. There is no `go.mod` here — the module root (and `go.mod`, `config`, `pkg/`, `rdf_schema.jsonld`) lives two levels up, at the `zhero/` root. Import paths like `github.com/domahidizoltan/zhero/pkg/...` resolve against that root.

## Commands

Run from the module root (`zhero/`), not from this directory:

```bash
# Test the whole domain layer
go test ./domain/...

# Test a single package
go test ./domain/schemaorg/

# Run a single test / subtest (subtests use t.Run names)
go test ./domain/schemaorg/ -run TestSchemaorg
go test ./domain/schemaorg/ -run 'TestSchemaorg/gets_schema_class'

# Vet + build
go vet ./domain/...
go build ./domain/...
```

Note: `schemaorg/service_test.go` reads the RDF schema from `../../rdf_schema.jsonld` (relative to the package), falling back to downloading from the Schema.org release URL in `config.RdfConfig.Source`. The file must exist at the module root or be fetchable for tests to pass.

## Architecture

The domain is split into four packages, each modeling one concept with a `Service` (in `service.go`) and its types (in `model.go`).

- **`schemaorg`** — Read-only access to the Schema.org RDF vocabulary. Wraps `pkg/rdf.Graph` (built on `deiu/rdf2go`) to answer questions about classes, subclass hierarchies, and properties. The module root's `rdf_schema.jsonld` is the data source.
- **`schema`** — User-defined "data blueprints" (`SchemaMeta` + `Property`). Persists blueprint metadata and consumes `schemaorg` to resolve real Schema.org classes and hierarchies.
- **`page`** — Content instances of a schema (`Page`). Holds the actual data, SEO metadata (`PageMeta`), references, and search values.
- **`route`** — Custom URL slugs mapping to pages, with versioning.

### Ports-and-adapters convention

Every service follows the same pattern, and new code must match it:

1. Define collaborator dependencies as **private interfaces** at the top of `service.go` (e.g. `pageRepo`, `routeSvc`, `schemaMetaRepo`, `schemaProvider`). These are the "ports."
2. The `Service` struct holds those interfaces, never concrete types.
3. A `NewService(...)` constructor injects the implementations. Repositories and cross-service wiring are supplied by the caller (outside this layer) — this domain layer never imports a concrete DB or HTTP type.

This is why `page.Service` depends on a `routeSvc` interface (not `route.Service` directly) and `schema.Service` depends on a `schemaProvider` interface (satisfied by `schemaorg.Service`).

### Cross-service dependencies

- `page` → `route`: creating/updating a `Page` with a non-empty `Route` calls `routeSvc.AssignRoute` inside the same transaction.
- `schema` → `schemaorg`: `schema.Service` calls `schemaProvider` to fetch classes and build the class hierarchy.

### Transactions

Write operations wrap their repo calls in `database.InTx(ctx, func(ctx) error { ... })` from `pkg/database`. The transaction is propagated through `context.Context` — repo methods pull the active tx off the context. Any multi-step write (e.g. `page.Create` inserting a page **and** assigning a route) must run inside a single `InTx` so it commits or rolls back atomically. Read methods call the repo directly without `InTx`.

### `onlyEnabled` flag

Page read methods take an `onlyEnabled bool`. Pass `true` for public-facing queries (only published content) and `false` for admin/editing contexts. Preserve this distinction when adding new read paths.

## schemaorg specifics

- **Unstable node filtering**: nodes marked `isPartOf` `attic.schema.org` or `pending.schema.org` are collected once at startup (`getUnstableNodes`) and excluded from all results via `filterUnstableValues`. When adding queries, route results through `prepareValues`/`filterUnstableValues` so unstable terms stay hidden (tests assert e.g. `PodcastSeason`, `BioChemEntity` are absent).
- **RDF terms** are constructed with `term(context, name)`, where `context` is a namespace constant (`schema`, `rdf`, `rdfs`). Most extra vocabularies are commented out in `model.go` — uncomment a `context` constant to enable it.
- **Names vs URLs**: `getTermName` strips the namespace prefix to get a short name (`LiveBlogPosting`); `RawValue()` gives the full canonical URL.
- Graph init and unstable-node scanning are guarded by `sync.Once`, so the graph loads a single time per process.
