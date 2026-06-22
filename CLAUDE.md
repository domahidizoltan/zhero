# CLAUDE.md

## Project Overview

Zhero is a Schema.org-driven, SEO-first headless CMS written in Go. It is designed for non-technical users to create, manage, and deliver structured content effortlessly. The system leverages Schema.org's vocabulary to ensure all content is machine-readable and optimized for search engines, generating JSON-LD for rich results. The application is a single self-contained executable, making it easy to deploy on various platforms, including low-resource hardware (e.g. Raspberry Pi Zero, or Android phones).

## Commands

Run Go commands from the module root (where `go.mod` lives).

- **Build:** `make build` (= `go build -o zhero main.go`); all packages: `go build ./...`
- **Run:** `./zhero` or `go run main.go` — starts BOTH servers: Admin on :7080, Public on :8080
- **Test (all):** `go test ./...`
- **Test (single package):** `go test ./domain/schemaorg/`
- **Test (single test):** `go test ./domain/schemaorg/ -run TestSchemaorg`
- **Test (subtest):** `go test ./domain/schemaorg/ -run 'TestSchemaorg/gets_schema_class'`
- **Vet / format:** `go vet ./...` / `go fmt ./...`
- **Lint:** `golangci-lint run ./...` and/or `staticcheck ./...` (assumed globally installed; no in-repo config)
- **Cross-compile RPi Zero:** `make build-rpi-zero`
- **Build Android AAR:** `make build-android-lib` (gomobile; required before the Android Gradle build can resolve the Go server)
- **Dev DB UI:** `cd localdev && docker compose up` (Adminer at http://localhost:9000)

Note: tests are intentionally minimal (only `domain/schemaorg` and `pkg/url`); project guidance currently says not to add tests.

## Architecture

Single Go module, layered (clean architecture). Dependencies flow inward:

`controller (Gin HTTP) → domain (Service interfaces + models) → repository (SQLite)`, with `pkg/*`, `config`, and `template` as shared support.

- **`main.go`** — process entry; starts/stops the server, waits on SIGINT/SIGTERM.
- **`server/server.go`** — composition root (DI wiring); inits SQLite, runs migrations, wires repos→services→router, starts the two Gin servers.
- **`controller/`** — HTTP layer (Gin) handlers and DTOs split into `router`, `adminschema`, `adminpage`, `dynamicpage`, `pagerenderer`, `preview`, `file`, `template`.
- **`domain/`** — entities and business logic in `schemaorg`, `schema`, `page`, `route`; ports-and-adapters (each Service declares private collaborator interfaces, never imports concrete DB/HTTP types).
- **`repository/`** — hand-written `database/sql` queries (no ORM) for `page`, `route`, `schema`.
- **`pkg/`** — shared libs (database/tx, handlebars, paging, rdf, jsonld, session, url, …).
- **`data/db/sqlite/`** — embedded SQL migrations + ordered `Scripts` slice.
- **`template/`** — embedded `.hbs`/`.css`/`.js` (Handlebars via `aymerick/raymond`, NOT html/template); frontend uses HTMX + Tailwind/DaisyUI (admin) via CDN.
- **`zhero-android-app/`** — Kotlin wrapper hosting the Go server via the gomobile AAR.

## Tech stack

- **Backend:** [Go](https://go.dev/), [Gin Web Framework](https://gin-gonic.com/)
- **Frontend:** [HTMX](https://htmx.org/), JavaScript, [Tailwind CSS](https://tailwindcss.com/), [DaisyUI CSS Framework](https://daisyui.com/)
- **Templating:** [Handlebars](https://github.com/aymerick/raymond) (via `aymerick/raymond`)
- **Database:** [SQLite](https://sqlite.org/index.html) (via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite))
- **Data Formats:** [JSON-LD](https://json-ld.org/) for structured data, RDF for [Schema.org](https://schema.org/) graph processing.
- **Logging:** [Zerolog](https://github.com/rs/zerolog)
- **Configuration:** [Viper](https://github.com/spf13/viper)

## Critical Cross-Cutting Concerns

1. **Context-propagated transactions** (`pkg/database`): `InTx(ctx, fn)` stores a `*sql.Tx` on the context; repository WRITE methods require a tx and return `ErrTransactionNotFound` if absent (they never open their own). Wrap multi-step domain writes (e.g. `page.Create` inserts page + assigns route) in a single `InTx`. Reads use `db` directly.
2. **4-stage reference system** (Go + JS must stay in sync): authoring token `#ZHERO#<target>#{json}#` → compact `#refN` on persist → expanded on load → resolved to `<a href>` on render. Regexes/tokens live in `controller/controller.go` and `template/admin/page/page.js`.
3. **Four derived projections must stay consistent:** `page.data` (JSON-LD) is the source of truth; `listable_data`, `meta`, `references` columns and the `page_search` FTS5 rows are denormalized from it. Every write updates all applicable projections.
4. **Component system:** a property's Schema.org type maps to an HTML component. Adding one requires updating `adminpage.determineComponent`, `pagerenderer`, `adminschema/controller.go`, and `template/admin/schemaorg/schemaorg.js`.
5. **`onlyEnabled bool`** threads through page reads: `true` = public (published only), `false` = admin/editing.
6. **Build-tag platform split:** default uses `modernc.org/sqlite` (pure Go, driver `"sqlite"`); `-tags android` uses `mattn/go-sqlite3` (cgo, driver `"sqlite3"`). Keep `*_android.go` and default files in sync.
7. **Two disjoint route sets** (`controller/router/router.go`): public (`/`, `/:class`, previews, NoRoute→page) and admin (`/admin/...`).

## Database

- SQLite, no ORM.
- See database migration specific details in `data/CLAUDE.md`
- See data persistence and repository layer specific details in `repository/CLAUDE.md`

## Per-Layer Docs

Each layer has its own detailed `CLAUDE.md` — consult them for specifics: `controller/`, `domain/`, `repository/`, `pkg/`, `data/`, `template/`, `zhero-android-app/`.

## Notes

- No CI/CD and no app Dockerfile (only a dev Adminer compose file).
- Use Context7 to get up-to-date technical documentations.
