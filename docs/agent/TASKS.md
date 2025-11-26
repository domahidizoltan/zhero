# Task List

This document lists the actionable tasks derived from the development plan.

## F001: Admin Interface

- [ ] `T001` Create a new Gin router group dedicated to the `/admin` path.
- [ ] `T002` Relocate all existing schema management routes (e.g., `/schema/list`, `/schema/new`) to be under the new `/admin` router group.
- [ ] `T003` Relocate all existing page management routes (e.g., `/{schemaName}/list`, `/{schemaName}/new`) to be under the new `/admin` router group.
- [ ] `T004` Update all frontend templates, links, and HTMX attributes (`hx-get`, `hx-post`, etc.) to point to the new `/admin` prefixed URLs.

## F002: Additional Actions & Refinements

- [ ] `T005` Create a new controller in a new package (`controller/preview`) for handling public page previews.
- [ ] `T006` Implement a new public-facing route, `/preview`, that accepts POST requests containing page data.
- [ ] `T007` Create a new service and repository for transforming page data into JSON-LD format.
- [ ] `T008` Update the "Preview" button of the page create/edit form in the admin UI.
- [ ] `T009` Implement client-side JavaScript to serialize the page form data and POST it to the `/preview` endpoint, displaying the result in a new browser tab.
- [ ] `T010` Create a new package (`pkg/renderer`) that transforms the JSON-LD data into a user-friendly HTML representation.
- [ ] `T011` The preview controller will use the new renderer package to process the posted data and return the final HTML.

## F003: Deployment

### Raspberry PI Zero Packaging
- [ ] `T012` Create a Makefile target named `build-rpi-zero` for cross-compiling the application.
- [ ] `T013` Configure the `build-rpi-zero` target to use the correct `GOOS=linux` and `GOARCH=arm` environment variables for Raspberry PI Zero.
- [ ] `T014` Add comments or documentation within the Makefile explaining how to build and run the application on the target device.

### Android Packaging
- [ ] `T015` Research and select the best method for packaging a Go web server application as an Android library (e.g., Gomobile).
- [ ] `T016` Create a proof-of-concept Android application that successfully integrates and calls a simple function from the compiled Go library.
- [ ] `T017` Refactor the Go application's `main.go` to expose functions for starting and stopping the server, making it controllable as a library. Pass the configuration (server port, database connection) from the Android app to the Golang lib.
- [ ] `T018` Create a build script or Makefile target to automate the process of building the Go library and packaging it into an Android APK.
- [ ] `T019` Ensure the asset embedding process works correctly for the Android build, so the final APK is self-contained.
