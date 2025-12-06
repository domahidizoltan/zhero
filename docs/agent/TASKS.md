# Task List

~ in progress, x completed, ! failed or blocked

## F001: Admin Interface

  - [x] T001: **Create a new Gin router group dedicated to the `/admin` path.**
    - Details: Implemented in `controller/router/router.go`.
    - Dependencies:
    - Comment:

  - [x] T002: **Relocate all existing schema management routes (e.g., `/schema/list`, `/schema/new`) to be under the new `/admin` router group.**
    - Details: Implemented in `controller/router/router.go`.
    - Dependencies: T001
    - Comment:

  - [x] T003: **Relocate all existing page management routes (e.g., `/{schemaName}/list`, `/{schemaName}/new`) to be under the new `/admin` router group.**
    - Details: Implemented in `controller/router/router.go`.
    - Dependencies: T001
    - Comment:

  - [x] T004: **Update all frontend templates, links, and HTMX attributes (`hx-get`, `hx-post`, etc.) to point to the new `/admin` prefixed URLs.**
    - Details: Updated `template/page/list.hbs`, `template/page/edit.hbs`, `template/schemaorg/edit.hbs`, `template/schemaorg/schemaorg.js`, `template/index.hbs`, and `template/page/main.hbs`.
    - Dependencies: T002, T003
    - Comment: All relevant frontend templates have been updated.

## F002: Additional Actions & Refinements

  - [ ] T005: **Create a new controller in a new package (`controller/preview`) for handling public page previews.**
    - Details: Here come the implementation details.
    - Dependencies:
    - Comment: Post-implementation comments about failures or impediments.

  - [ ] T006: **Implement a new public-facing route, `/preview`, that accepts POST requests containing page data.**
    - Details: Here come the implementation details.
    - Dependencies: T005
    - Comment: Post-implementation comments about failures or impediments.

  - [ ] T007: **Create a new service and repository for transforming page data into JSON-LD format.**
    - Details: Here come the implementation details.
    - Dependencies:
    - Comment: Post-implementation comments about failures or impediments.

  - [ ] T008: **Update the "Preview" button of the page create/edit form in the admin UI.**
    - Details: Here come the implementation details.
    - Dependencies: T006, T007, T009
    - Comment: Post-implementation comments about failures or impediments.

  - [ ] T009: **Implement client-side JavaScript to serialize the page form data and POST it to the `/preview` endpoint, displaying the result in a new browser tab.**
    - Details: Here come the implementation details.
    - Dependencies: T006
    - Comment: Post-implementation comments about failures or impediments.

  - [ ] T010: **Create a new package (`pkg/renderer`) that transforms the JSON-LD data into a user-friendly HTML representation.**
    - Details: Here come the implementation details.
    - Dependencies:
    - Comment: Post-implementation comments about failures or impediments.

  - [ ] T011: **The preview controller will use the new renderer package to process the posted data and return the final HTML.**
    - Details: Here come the implementation details.
    - Dependencies: T006, T010
    - Comment: Post-implementation comments about failures or impediments.

## F003: Deployment

### Raspberry PI Zero Packaging
  - [ ] T012: **Create a Makefile target named `build-rpi-zero` for cross-compiling the application.**
    - Details: Here come the implementation details.
    - Dependencies:
    - Comment: Post-implementation comments about failures or impediments.

  - [ ] T013: **Configure the `build-rpi-zero` target to use the correct `GOOS=linux` and `GOARCH=arm` environment variables for Raspberry PI Zero.**
    - Details: Here come the implementation details.
    - Dependencies: T012
    - Comment: Post-implementation comments about failures or impediments.

  - [ ] T014: **Add comments or documentation within the Makefile explaining how to build and run the application on the target device.**
    - Details: Here come the implementation details.
    - Dependencies: T012, T013
    - Comment: Post-implementation comments about failures or impediments.

### Android Packaging
  - [ ] T015: **Research and select the best method for packaging a Go web server application as an Android library (e.g., Gomobile).**
    - Details: Here come the implementation details.
    - Dependencies:
    - Comment: Post-implementation comments about failures or impediments.

  - [ ] T016: **Create a proof-of-concept Android application that successfully integrates and calls a simple function from the compiled Go library.**
    - Details: Here come the implementation details.
    - Dependencies: T015
    - Comment: Post-implementation comments about failures or impediments.

  - [ ] T017: **Refactor the Go application's `main.go` to expose functions for starting and stopping the server, making it controllable as a library. Pass the configuration (server port, database connection) from the Android app to the Golang lib.**
    - Details: Here come the implementation details.
    - Dependencies: T016
    - Comment: Post-implementation comments about failures or impediments.

  - [ ] T018: **Create a build script or Makefile target to automate the process of building the Go library and packaging it into an Android APK.**
    - Details: Here come the implementation details.
    - Dependencies: T017
    - Comment: Post-implementation comments about failures or impediments.

  - [ ] T019: **Ensure the asset embedding process works correctly for the Android build, so the final APK is self-contained.**
    - Details: Here come the implementation details.
    - Dependencies: T018
    - Comment: Post-implementation comments about failures or impediments.
