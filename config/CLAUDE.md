# Config

## Scope

Viper-based configuration for the whole app. `LoadConfig(absolutePath)` reads
`<absolutePath>config.yaml` (file must sit next to the binary — committed, no secrets),
unmarshals via `mapstructure` tags, injects `Env.AbsolutePath`.

## Key ↔ struct mapping

| YAML key | Go struct | Notes |
|---|---|---|
| `log.*` | `LogConfig` | format, level, color |
| `database.sqlite.file` | `DBConfig` | SQLite file path |
| `app.pagination.jump` | `AppConfig` | feeds `paging.SetJump` in server |
| `app.pagination.defaultPageSize` | `AppConfig` | default page size |
| `admin.server.port` | `ServerConfig` | default 7080 |
| `admin.rdf.file` | `RdfConfig` | local schemaorg file (gitignored) |
| `admin.rdf.source` | `RdfConfig` | remote schemaorg URL |
| `admin.uploads.thumbnailSize` | `UploadsConfig` | parsed `WxH` by `GetThumbnailSize`, default 360×180 |
| `admin.uploads.supportedFileTypes` | `UploadsConfig` | allowed upload extensions |
| `admin.uploads.supportedImageTypes` | `UploadsConfig` | allowed image extensions |
| `public.server.port` | `ServerConfig` | default 8080 |
| `env.platform` | `EnvConfig` | optional, e.g. "android" |

`Env.AbsolutePath` is injected post-load (not from yaml). Config.yaml is committed and
secret-free. Adding a key = edit `config.yaml` **and** the struct tag together (per
security-guidelines).