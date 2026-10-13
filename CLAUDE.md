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
- **Verify (fmt + vet + lint + test):** `make check`
- **Lint:** `golangci-lint run ./...` (= `make lint`; config: `.golangci.yml`)
- **Cross-compile RPi Zero:** `make build-rpi-zero`
- **Build Android AAR:** `make build-android-lib` (gomobile; required before the Android Gradle build can resolve the Go server)
- **Dev DB UI:** `cd localdev && docker compose up` (Adminer at http://localhost:9000)

Note: tests are deliberately minimal — subtle-logic tests only (reference round-trip, projection derivation, component mapping); see `.claude/rules/testing.md`. `domain/schemaorg` requires `rdf_schema.jsonld` at the repo root (gitignored; auto-downloaded on first run).

## Architecture

Single Go module, layered (clean architecture). Dependencies flow inward:

`controller (Gin HTTP) → domain (Service interfaces + models) → repository (SQLite)`, with `pkg/*`, `config`, and `template` as shared support.

- **`main.go`** — process entry; starts/stops the server, waits on SIGINT/SIGTERM.
- **`server/server.go`** — composition root (DI wiring) and gomobile binding target (`make build-android-lib` → `gomobile bind ./server`); `Start()` sequence: LoadConfig → helpers → logging → SetJump → InitSqliteDB → Migrate(Scripts) → repos→services→router → both Gin servers; `Stop()` graceful 5s shutdown.
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

1. **Context-propagated transactions** (`pkg/database.InTx`): repo WRITE methods require an ambient tx (`ErrTransactionNotFound`) and never open their own; wrap multi-step domain writes in one `InTx`. Details: `pkg/CLAUDE.md`.
2. **4-stage reference system** (Go + JS in sync): `#ZHERO#<target>#{json}#` → `#refN` → expanded → `<a href>`. Details: `controller/CLAUDE.md` (Go), `template/CLAUDE.md` (JS).
3. **Four derived projections** must stay consistent: `page.data` is the source of truth; `listable_data`, `meta`, `references`, `page_search` are denormalized from it. Details: `wiki/data-mapping-overview.md`.
4. **Component system:** Schema.org type/name → HTML component; adding one touches `adminpage.determineComponent`, `pagerenderer`, `adminschema/controller.go`, `template/admin/schemaorg/schemaorg.js`. Details: `controller/CLAUDE.md`; mapping table in `wiki/dynamic-data-conventions.md` §5.
5. **`onlyEnabled bool`** on page reads: `true` = public/published only, `false` = admin. Details: `domain/CLAUDE.md`.
6. **Build-tag platform split:** default uses `modernc.org/sqlite` (pure Go, driver `"sqlite"`); `-tags android` uses `mattn/go-sqlite3` (cgo, driver `"sqlite3"`). Keep `*_android.go` and default files in sync. Details: `repository/CLAUDE.md`.
7. **Two disjoint route sets** (`controller/router/router.go`): public (`/`, `/:class`, previews, NoRoute→page) and admin (`/admin/...`). Details: `controller/CLAUDE.md`.

## Database

- SQLite, no ORM.
- See database migration specific details in `data/CLAUDE.md`
- See data persistence and repository layer specific details in `repository/CLAUDE.md`

## Per-Layer Docs

Each layer has its own detailed `CLAUDE.md` — consult them for specifics: `controller/`, `domain/`, `repository/`, `pkg/`, `data/`, `template/`, `config/`, `zhero-android-app/`.

## Notes

- No CI/CD and no app Dockerfile (only a dev Adminer compose file).
- Use Context7 to get up-to-date technical documentations.
- See the project knowledgebase table of contents in `wiki/index.md`
- `rich_results.yaml` — reference data mapping Google rich-result features → Schema.org; not yet referenced by code; top-level key misspelled `richResuls` — fix before wiring up.