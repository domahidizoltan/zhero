# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Android wrapper (`zhero.app`) that hosts an embedded Go HTTP server inside a foreground
service. The Go server is consumed as a **gomobile/gobind AAR binding** (Kotlin imports
`server.Server` / `server.Server_`). This module is the thin native shell; all server logic
lives in the parent Go project and is compiled into the AAR.

## Build & test commands

Run from this directory (`zhero-android-app/`):

```bash
./gradlew assembleDebug                 # build debug APK
./gradlew installDebug                  # build + install on connected device/emulator
./gradlew build                         # full build incl. lint + tests
./gradlew clean
./gradlew lint                          # Android lint -> app/build/reports/lint-results-*.html

./gradlew test                          # JVM unit tests (app/src/test)
./gradlew testDebugUnitTest --tests "zhero.app.ExampleUnitTest"          # single class
./gradlew testDebugUnitTest --tests "zhero.app.ExampleUnitTest.addition_isCorrect"  # single method
./gradlew connectedAndroidTest          # instrumented tests (needs device, app/src/androidTest)
```

## Critical prerequisite: the server AAR

`app/build.gradle.kts` pulls native binaries via `fileTree("app/libs", include = ["*.aar", "*.jar"])`.
**`app/libs/` is not checked in and does not exist by default.** Without the gomobile-generated
AAR placed there, `import server.Server` will not resolve and the build fails. The AAR is produced
by the parent Go project's gomobile bind step — generate it there and copy the output into
`app/libs/` before building.

`kotlinx.coroutines` is used in `ServerService` but is **not declared** in `libs.versions.toml`;
it resolves transitively (via the AndroidX navigation/appcompat deps). Don't assume it's a direct
dependency when reasoning about versions.

## Architecture

Two Kotlin files in `app/src/main/java/zhero/app/`:

- **`MainActivity.kt`** — single-screen UI (`activity_main.xml`, ViewBinding). One toggle button
  starts/stops the server and a scrolling log view. Server control is **fire-and-forget via Intent
  actions** (`START_SERVER` / `STOP_SERVER`) sent to `ServerService` — there is no service binding.
  Also drives runtime storage-permission requests.

- **`ServerService.kt`** — foreground service (`foregroundServiceType="dataSync"`) that owns the Go
  server lifecycle on a `Dispatchers.IO` coroutine scope: `Server.new_()` →
  `setAbsolutePath(<external storage root>/)` → `start()` / `stop()`. The server's base path is the
  device's external storage root, so it reads/writes user files there.

### Cross-component conventions (read before changing)

- **Log bridge:** `MainActivity` displays Go output by shelling out to `Runtime.exec("logcat -s GoLog")`
  on a background thread and appending lines to the TextView. `ServerService`, however, logs under the
  tag `ServerService` (and Go stderr is re-logged under that tag). The filter tag `GoLog` and the
  service's actual tags must agree for logs to appear — keep them in sync when touching either side.
- **Permissions:** API 30+ requests `MANAGE_EXTERNAL_STORAGE` (All Files Access) via Settings intent;
  below 30 requests `WRITE_EXTERNAL_STORAGE`. Both paths share `PERMISSION_REQUEST_CODE = 100`.
  The service assumes this access is already granted when it sets the server's absolute path.
- `isServerRunning` is tracked as local UI state in `MainActivity`; it is not queried back from the
  service, so it can drift if the service stops on its own.

## Config facts

- Gradle 8.13 (wrapper), AGP 8.13.2, Kotlin 2.0.21, version catalog at `gradle/libs.versions.toml`.
- `minSdk 21`, `compileSdk`/`targetSdk 36`, Java/JVM target 11, ViewBinding on, namespace `zhero.app`.
- Release build has `isMinifyEnabled = false` (ProGuard rules are present but effectively unused).
