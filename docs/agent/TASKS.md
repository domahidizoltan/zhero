# Task List

**Status Legend:**
- `[ ]` = ready to implement
- `[x]` = completed
- `[~]` = in progress
- `[!]` = blocked or failed

---

## F001: Custom Route Management Feature

### Foundation Tasks

- [x] T001: Create slug package with Slugify function
  - Details: Create `pkg/slug/slug.go` with regex-based slugification (lowercase, alphanumeric + hyphens). Add unit test file `pkg/slug/slug_test.go`.
  - Dependencies: None
  - Comment: Changed to pkg/url/url.go per user request

- [x] T002: Register slugify Handlebars helper
  - Details: Add `"slugify": slug.Slugify` to helpers map in `pkg/handlebars/handlebars.go`. Call `InitHelpers()` on app startup (already done).
  - Dependencies: T001
  - Comment:

- [x] T003: Update database schema with route table
  - Details: Append SQL to `data/db/sqlite/0000_init_schemas.sql` to create `route` table with indexes `idx_route_route` and `idx_route_page_version`. No migration needed; table created on startup if not exists.
  - Dependencies: None
  - Comment:

### Domain & Repository Layer

- [x] T004: Define Route domain model and repo interface
  - Details: Create `domain/route/model.go` with Route struct, RouteRepo interface (Create, GetByRoute, GetLatestVersion), and Service struct. Document method contracts.
  - Dependencies: None
  - Comment:

- [x] T005: Implement Route repository
  - Details: Create `repository/route/repository.go`. Implement all interface methods using `database.InTx` for Create (read existing version + insert atomically). Add proper SQL queries with indexes in mind.
  - Dependencies: T004, T003 (table must exist)
  - Comment:

- [x] T006: Implement Route service with AssignRoute logic
  - Details: Create `domain/route/service.go` (or add to model.go). Implement AssignRoute with: (1) validation (starts with `/`, no trailing slash, length, URL-safe chars), (2) check if route exists for other page (duplicate check), (3) get latest version for page, (4) if unchanged return nil, else increment version and create. Return meaningful errors.
  - Dependencies: T005
  - Comment:

### Server Wiring

- [x] T007: Add route.Service to Services struct and wire in server
  - Details: Modify `controller/router/router.go` Services struct to include `Route route.Service`. In `server/server.go`, instantiate route repository and service, inject into Services.
  - Dependencies: T006
  - Comment:

### Page Service Integration

- [x] T008: Integrate route assignment in page save workflow
  - Details: Modify `domain/page/service.go` Save method. After page is saved (Insert or Update), if DTO contains non-empty Route field, call `routeService.AssignRoute(pageKey, customRoute)` where `pageKey = schemaName + "/" + identifier`. Ensure transaction handling if needed.
  - Dependencies: T007, T004 (need route service)
  - Comment:

### Admin UI - Data Layer

- [x] T009: Update page DTO with Route field
  - Details: Add `Route string` field to `pageDto` in `controller/adminpage/dto.go`. Update `EnhanceFromForm` to read `c.FormValue("route")`. Update `enhanceFromModel` to fetch latest route via routeService.GetLatestVersion and populate DTO.Route. Update `ToModel` to export Route value.
  - Dependencies: T007 (route service available), T004 (page service can call route service?)
  - Comment: Need to inject route service into DTO methods or pass it in

- [x] T010: Update admin controller to use DTO.Route
  - Details: In `controller/adminpage/controller.go`, the `edit` function already uses DTO. Ensure that after enhancement, DTO.Route is populated for edit form. Ensure Save passes DTO with Route to pageSvc.Update/Create (via ToModel).
  - Dependencies: T009
  - Comment:

### Admin UI - Template

- [x] T011: Add collapsible Custom Route section in edit template
  - Details: Modify `template/admin/page/edit.hbs`. Add `<details>` element with DaisyUI classes. Include `<input name="route" value="{{page.Route}}">`. Show placeholder with default: `/{{class}}/{{slugify page.secondaryIdentifier}}`. Include help text. Section should be open if `{{page.Route}}` is non-empty.
  - Dependencies: T010, T002 (slugify helper available)
  - Comment:

### Public Router

- [x] T012: Implement resolveRoute middleware
  - Details: In `controller/router/router.go`, add `resolveRoute` function. Logic: (1) lookup requestedPath in route table; if found, check if latest version, redirect if outdated, else load page; (2) if not custom route, try parsing as default `/:class/:identifier`; check page enabled; check if page has custom route (redirect to it if yes); load page; (3) 404 if no match. Use existing `loadPage` pattern or call dynamicPageCtrl.LoadPage after setting params.
  - Dependencies: T007 (route service available), T004 (page service)
  - Comment:

- [x] T013: Update SetPublicRoutes with catch-all route
  - Details: In `controller/router/router.go`, modify `SetPublicRoutes` to add `router.GET("/*path", resolveRoute)` as the last route (after `/preview` and `/:class`). Need to capture route service from Services. Update Services struct to include Route service.
  - Dependencies: T012
  - Comment:

### Validation & Error Handling

- [x] T014: Add duplicate route and validation error handling
  - Details: Already implemented in T006, but ensure errors propagate to admin UI. In admin controller, catch routeService errors during page save and display to user (FlashMsg or ErrorMsg). Update template to show error if present.
  - Dependencies: T006, T008
  - Comment:

### Testing & Polish

- [ ] T015: Manual testing and bug fixes
  - Details: Run through manual testing checklist: custom route creation, route change (301 redirect), removal, disabled pages 404, duplicate route rejection, slugify correctness, preview functionality. Fix any bugs discovered.
  - Dependencies: T014, T011, T013
  - Comment:

---

## Task Dependencies Summary

```
T001 (slug) → T002 (register helper) → T011 (template slugify)
         ↘
          T004 (route model) → T005 (repo) → T006 (service) → T007 (wire) → T012 (resolver) → T013 (router)
                              ↘
                               T009 (DTO) → T010 (controller) → T011 (template)
                              ↘
                               T008 (page service) → T014 (error handling) → T015 (testing)

T003 (DB schema) can be done any time before T005
```

**Key dependency chains:**
- Route feature: T001 → T002 → T004 → T005 → T006 → T007 → T012 → T013
- Admin UI: T007 → T009 → T010 → T011
- Integration: T007 → T008 → T014
- Final: All above → T015

---

## Notes

- Tests (`*_test.go` files) are skipped for MVP per user request, but should be added later.
- Preview functionality should work without changes (uses `/preview/:class/:identifier`).
- Enabling/disabling pages uses existing mechanism; router only serves enabled pages (checked in `GetPageBySchemaNameAndIdentifier` with `onlyEnabled=true`).
- Default routes (`/:class/:identifier`) bypass route table entirely unless custom route exists (then redirect).
- All old route versions retained for 301 redirects.
