# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Scope

This directory (`pkg/`) holds the shared, reusable library packages for the **zhero** application
(module `github.com/domahidizoltan/zhero`). There is no `go.mod` here — the module root is the
parent directory. These packages are consumed by the app's sibling packages (`config`, `template`,
`domain/page`) which live outside `pkg/`.

zhero is a schema.org-driven web CMS: pages are typed by schema.org schemas, structured data is
emitted as JSON-LD, and schema definitions are loaded from an RDF graph.

## Commands

```bash
go test ./...                      # test all pkg packages (only ./url has tests today)
go test ./url -run TestSlugify     # run a single test
go vet ./...                       # vet all pkg packages
go build ./...                     # build all pkg packages (default/non-android)
go build -tags android ./...       # build the android variant (requires cgo)
```

Wildcard commands (`./...`) are run from this `pkg/` directory and only cover `pkg/*` packages.
Use the parent module root to exercise the whole application.

## Build tags & platform

- **SQLite driver is swapped by build tag**, and the driver *name string* differs:
  - default (`//go:build !android`) → `modernc.org/sqlite`, driver name `"sqlite"` (pure Go, no cgo)
  - `//go:build android` → `github.com/mattn/go-sqlite3`, driver name `"sqlite3"` (cgo)
  See `database/database.go` vs `database/database_android.go`.
- Android targets build via gomobile and **require cgo** (the `android/386 requires external (cgo)
  linking` diagnostic is expected when cgo is off — it is not a code bug).
- `logging.ConfigureLogging` force-disables colored output when `cfg.Env.Platform == "android"`.

## Gotchas

- **`pkg/_err` is invisible to `./...`.** Go tooling ignores directories whose names start with `_`,
  so `_err` is excluded from `go test/vet/build ./...` even though it is imported explicitly as
  `github.com/domahidizoltan/zhero/pkg/_err` (e.g. by `session`). Changes there are *not* covered by
  wildcard commands — verify packages that import it.

## Architecture

Each package is a focused, mostly-stateless helper layer. Cross-cutting themes:

- **Web stack:** Gin (HTTP), `aymerick/raymond` (Handlebars templating), HTMX, and `blackfriday`
  markdown. `handlebars` registers template helpers/partials and produces HTMX-aware HTML
  (e.g. `htmxSortButton` emits `hx-get`/`hx-target` anchors); `paging` builds HTMX pagination DTOs
  and parses sort/search query params from Gin requests.
- **Semantic web / schema.org:** `rdf` loads a schema.org JSON-LD file into an `rdf2go` graph
  (downloaded once via `sync.Once` + `file.DownloadToPath`) and exposes triple queries. `jsonld`
  serializes a `domain/page.Page` into schema.org JSON-LD (`@context`/`@type`/`@id` from the page's
  `Data` map and `Identifier`).
- **Persistence:** `database` owns a package-level `*sql.DB` singleton plus **context-based
  transactions** — `InTx` stores a `*sql.Tx` under a context key, nested `InTx` calls reuse the
  existing tx, and `GetTx` retrieves it. Migrations are applied as a slice of raw SQL strings.
- **Session/flash:** `session` wraps `gin-contrib/sessions` (in-memory `memstore`) for flash messages.

## Conventions

- **Errors:** declare package-level sentinel errors (`ErrXxx`) and double-wrap with
  `fmt.Errorf("%w: %w", ErrSentinel, err)` to preserve both the category and the cause. `_err.WrapNotNil`
  wraps only when the result error is non-nil.
- **Iterators:** `collection` returns Go 1.23+ range-over-func iterators (`iter.Seq`/`iter.Seq2`)
  rather than materialized slices.
- **Logging:** use `rs/zerolog` (`log.Debug()/Err()` etc.), configured centrally by
  `logging.ConfigureLogging`.
- **Package-level state:** several packages keep deliberate global state (`database.db`,
  `paging.jump` via `SetJump`, `rdf.once`). Initialize it before use.
