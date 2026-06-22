# Presentation layer

## What this is

`package template` is the presentation layer of the **zhero** CMS. It bundles Handlebars
templates (`.hbs`) and static assets (`.css`, `.js`) into the Go binary via `go:embed` and
exposes them as package-level variables for the parent application to render and serve.

There is no `go.mod` here — this directory is one package inside a larger Go module rooted
above `/template`. It has no standalone build, run, or test step; it compiles when the parent
module builds. Templates are validated at parent-init time (see "How `template.go` works").

Rendering uses **`github.com/aymerick/raymond`** (a Go Handlebars engine), not Go's
`html/template`.

## How `template.go` works (read this before adding/renaming files)

`template.go` is the single source of truth wiring files into the binary. Three things must
stay in sync or the parent app **panics at init** (`mustParse`/`mustLoad` panic on missing files):

1. **`go:embed` directives** — the file must match one of the embed globs
   (`admin/*`, `*.css *.hbs *.js`, `paging/*`).
2. **Parsed-template / asset vars** — each rendered template needs a `mustParse(...)` var;
   each served asset needs a `mustLoad(...)` entry in the `Assets` or `AdminAssets` map.
3. **`RegisterPartials()`** — any Handlebars partial must be registered here before use.
   Currently: `editProperty` (scoped to `AdminSchemaorgEdit`) and `pagination` (global).

`Assets` / `AdminAssets` are `map[urlPath][]byte`; the parent serves them under `/asset/...`.

## Architecture

Two rendering surfaces, each a shell layout that injects a rendered `{{body}}`:

- **Public site** — `index.hbs` (shell) + `page_not_found.hbs`. Vanilla CSS (`index.css`),
  vanilla JS (`index.js`: search popup, 18+ age-gate driven by `meta.rating == "adult"`).
- **Admin CMS** — `admin/index.hbs` (shell) + the `admin/page/` and `admin/schemaorg/` views.
  Styled with Tailwind + DaisyUI (CDN), interactive via **HTMX**, with TomSelect (dropdowns)
  and SortableJS (drag-to-reorder).

The CMS domain model is **Schema.org-driven**:

- **Schema** (`admin/schemaorg/`) — pick a Schema.org class, choose which properties are
  visible/hidden/mandatory/searchable/listable, set an Identifier + SecondaryIdentifier, and
  order properties. Identifier config is **create-once** (fieldsets disabled after load).
  Each property maps a Schema.org *type* → an HTML *component* (`getTypeComponents` in
  `schemaorg.js` defines the type→component matrix, e.g. `Text → TextInput/TextArea/...`,
  image-ish → `URL/File/Image`).
- **Page** (`admin/page/`) — a content instance of a schema. `edit.hbs` renders one input per
  field, branching on `component` (`TextInput`, `TextArea`, `Checkbox`, `URL`, `File`, `Image`,
  `ReferenceSearch`, plus native `inputType` inputs). `list.hbs` is an HTMX-driven, sortable,
  paginated table.

## Conventions and gotchas

- **Custom Handlebars helpers are defined in ../pkg/handlebars/, not here.** Templates depend on
  them being registered on the raymond engine: `join`, `eq`, `equal`, `contains`, `concat`,
  `use`, `compareAndUse`, `mapValue`, `replace`, `beautify`, `htmxSortButton`, and the
  `eachMenuItem` block helper (exposes `@menu`). When adding a helper-like construct, register
  it in the parent — do not assume it exists.
- **`{{body}}` is raw HTML injection** — the parent renders an inner template and passes the
  result as `body` into a shell. `{{{triple-stache}}}` (e.g. `beautify`) also emits unescaped HTML.
- **Reference tokens** — cross-page references are embedded inside field values as the inline
  token `#ZHERO#<type>/<id>#{json meta}#` (meta carries `linkText`/`altText`). Parsing/insertion
  lives in `admin/page/page.js` (`refPattern`, `setReference`); the picker UI is
  `admin/page/search-reference.hbs`. Keep the token format identical across JS and any Go that
  parses it.
- **File uploads** — `File`/`Image` fields POST to `/admin/file` via HTMX with
  `X-File-Dir`, `X-File-Name-Prefix`, `X-Form-File-Name` headers. Stored files are served from
  `/zhero-content/uploads/<class>/<identifier>/<file>`; image thumbnails use the `_thumb`
  filename convention (`photo.jpg` → `photo_thumb.jpg`), applied in both `edit.hbs` and `page.js`.
- **Preview runs on a separate host/port** — `page.js` hardcodes `previewHost` to the current
  host on port **8080**, posting the edit form to `/preview/<class>[/<identifier>]`.
- **Listable cap** — `countListableProperties` in `schemaorg.js` enforces a max of ~3 extra
  listable properties (plus SecondaryIdentifier and one image/thumbnail) for public listings.
- **HX-Trigger conventions** — JS reacts to server-sent `HX-Trigger` headers: `showError`
  (popup), `fileUploaded` (`{field, fileName}` → fills hidden field + thumbnail). Match these
  names when emitting triggers server-side.
