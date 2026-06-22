# Helpers and Shared utilities

## Scope

This directory (`pkg/`) holds the shared, reusable library packages for the **zhero** application
(module `github.com/domahidizoltan/zhero`). There is no `go.mod` here — the module root is the
parent directory. These packages are consumed by the app's sibling packages (`config`, `template`,
`domain/page`) which live outside `pkg/`.

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
