# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Scope

This directory is the **controller (HTTP) layer** of Zhero, a schema-driven CMS. It is one package tree inside the larger Go module `github.com/domahidizoltan/zhero`, whose root (with `go.mod`, `main`, `domain/`, `template/`, `pkg/`, `config/`) lives in the parent directory. Code here imports those siblings via the full module path (e.g. `github.com/domahidizoltan/zhero/domain/page`).

## Commands

Go commands must run from the **module root** (parent dir), since `go.mod` lives there. Target this layer with the `./controller/...` pattern:

```bash
go build ./controller/...                         # build this layer
go test ./controller/...                          # test this layer
go test ./controller/dynamicpage/ -run TestName   # single package / single test
go vet ./controller/...
golangci-lint run ./controller/...                # linter (staticcheck-style QF/fmtappendf hints are enforced)
```

There are currently no `_test.go` files in this layer.

## Architecture

Zhero is a **schema-driven CMS**. The data model has three layers, all defined in sibling `domain/` packages and injected into controllers as service interfaces:

- **schema.org classes** (`domain/schemaorg`) — the catalog of available types and their possible properties.
- **SchemaMeta** (`domain/schema`) — a user-selected/configured subset of a schema.org class. Defines `Identifier`, `SecondaryIdentifier`, and `Properties` (each with `Type`, `Component`, `Mandatory/Searchable/Listable` flags, `Order`).
- **Page** (`domain/page`) — a concrete instance of a SchemaMeta, holding `Data`, `ListableData`, SEO `Meta`, and `References`.

Controllers never construct services; `router.Services` is populated at startup and passed in. Each sub-package exposes a `NewController(...)` taking the service interfaces it needs.

### Request flow

Two disjoint route sets are wired in `router/router.go`:

- **`SetPublicRoutes`** — serves rendered pages to end users. `GET /:class` lists pages; `NoRoute` falls through to `dynamicpage.LoadPage` for individual pages; `/preview/*` renders unsaved/saved previews.
- **`SetAdminRoutes`** — everything under `/admin` for managing schemas, pages, and files. Responses are HTML fragments driven by **HTMX** (note `HX-Trigger` headers and `hx-*` attributes in generated markup).

Public requests pass through **`CustomRouteMiddleware`** (`router/middleware.go`), which is the non-obvious heart of routing:

- It resolves the request path against `route.Service` (custom slugs with **versioning**), issuing `301` redirects when a page's route is outdated or has a canonical custom slug.
- When no custom route matches, it parses the path into `class` + `identifier` and injects them as Gin params via `setParams`/`mergeParams`.
- `skipPrefixes` (`static`, `asset`, `preview`, `favicon.ico`) bypass this and set `skipLoadPage`.

### Rendering pipeline

`dynamicpage.Controller` orchestrates rendering but delegates HTML generation through the `controller.UserFacingPageRenderer` / `UserFacingPageListRenderer` interfaces (defined in `controller.go`). The concrete implementation is `pagerenderer.DynamicPageRenderer`, which walks `SchemaMeta.Properties` and emits `<p>/<h1>/<img>/<a>` based on each property's `Component` (`File`, `Image`, default). Output is then wrapped in the site layout by `controller/template.WithLayout` (public) or `template.AdminIndex` (admin), which call into the sibling `template` package's Handlebars templates (via `aymerick/raymond`).

`preview.Controller` reuses `dynamicpage.Controller`: `LoadPage` previews saved pages; `InFlightPage` renders directly from posted form fields without persisting (used for live admin preview).

### The reference system (cross-page links)

This is the most subtle mechanism and spans several files. References let a page field link to another page:

1. **Authoring**: a field value embeds a long reference matching `RefPattern` (`#ZHERO#<target>#{<json meta>}#`), defined in `controller.go`.
2. **Persisting** (`adminpage/dto.go` `extractReferences`): long refs are replaced with compact `#refN` tokens (`ShortRefPattern`), and the actual targets are stored in `Page.References` (index = N).
3. **Loading** (`dto.go` `replaceShortRefs`, `dynamicpage` `transformShortRef`): `#refN` tokens are expanded back to long refs using the stored slice.
4. **Resolving** (`dynamicpage` `transformReferences`): long refs are looked up via `pageSvc.ListEnabledRoutesByRef` and rewritten to `<a href>` anchors (or plain text when the target is disabled/missing, or when `disableAllLinks` is set, e.g. in previews). Resolved links land in `data["references"]`, which the renderer substitutes into field values.

When touching references, keep all four stages in sync — changing one regex or token format breaks round-tripping.

### Package map

| Package | Responsibility |
|---|---|
| `controller` (root) | Shared error responses (`BadRequest`, `InternalServerError`, `TemplateRenderError`), renderer interfaces, and the `RefPattern`/`ShortRefPattern` regexes. |
| `router` | Public + admin route tables, `CustomRouteMiddleware`, static/asset serving, public Handlebars helper registration. |
| `adminschema` | Configure a SchemaMeta from a schema.org class (component list, mandatory/searchable/listable flags, property order). Binds forms via `gin` + `validator`. |
| `adminpage` | Page CRUD: create/edit/list/save, enable/disable/delete actions, slug validation, reference search. `dto.go` mediates form ⇄ `page.Page` and owns reference extraction. |
| `dynamicpage` | Public page rendering: `List`, `LoadPage`, `Render`, and all reference transformation. |
| `pagerenderer` | Pure HTML generation from `SchemaMeta` + data; component-aware. |
| `preview` | Saved + in-flight (unsaved form) previews, wrapping `dynamicpage`. |
| `file` | Upload/delete handlers with type validation and thumbnail generation. Uploads live under `UploadsPath` (`/zhero-content/uploads`). |
| `controller/template` | Layout wrappers (`WithLayout`, `AdminIndex`, `PageNotFoundLayout`) and flash-message handling via `pkg/session`. |

## Conventions

- **Logging** is structured via `zerolog` (`log.Error().Err(err).Str(...).Msg(...)`). Error helpers in `controller.go` already log; don't double-log.
- **Components** are derived by name/type heuristics in `adminpage.determineComponent` and consumed by `pagerenderer`. Adding a new component means updating both, plus the component list in `adminschema/controller.go`.
- DTO structs (`adminpage/dto.go`, `adminschema/dto.go`) are the boundary between domain models and templates; convert through `…From`/`ToModel`/`ToMap` rather than passing domain types to templates directly.
