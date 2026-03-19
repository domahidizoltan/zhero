# Route Management Implementation Plan

## Overview
Add custom route support to pages with a collapsible field in the admin edit form, a database table to track routes with versioning, and router updates to handle both default and custom routes with proper redirects.

---

## Key Decisions Confirmed

1. **Route format**: Slug only (no ULID) - e.g., `/product/my-product`
2. **Slugify**: Entire secondary identifier converted to lower-kebabcase as a single slug
3. **Route table scope**: Only custom routes are tracked; default routes (`/:class/:identifier`) work without table entries
4. **Version history**: Keep all old route versions for 301 redirects (never delete)

---

## Phase 1: Database Schema

### 1.1 Update Database Schema
**File**: `data/db/sqlite/0000_init_schemas.sql` (append to existing)

```sql
-- Route table for custom URL routing with versioning
CREATE TABLE IF NOT EXISTS route (
    route TEXT NOT NULL UNIQUE,           -- URL path (e.g., '/product/my-product')
    page TEXT NOT NULL,                   -- Page key in format '{schema_name}/{identifier}'
    version INTEGER NOT NULL DEFAULT 1   -- Increments on each route change
);
-- Note: For MVP we skip foreign key constraints; use application-level integrity
-- The `page` column stores "{schema_name}/{identifier}" to match how pages are identified

-- Indexes for fast lookups
CREATE INDEX IF NOT EXISTS idx_route_route ON route(route);
CREATE INDEX IF NOT EXISTS idx_route_page_version ON route(page, version DESC);
```

---

## Phase 2: Domain Layer

### 2.1 Route Model
**File**: `domain/route/model.go` (new)

```go
package route

type Route struct {
    Page   string // Format: "{schema_name}/{identifier}"
    Route  string // URL path (e.g., '/product/my-product')
    Version int
}

type RouteRepo interface {
    Create(ctx context.Context, route *Route) error
    GetByRoute(ctx context.Context, route string) (*Route, error)
    GetLatestVersion(ctx context.Context, page string) (*Route, error) // page = "schema/identifier"
}

type Service struct {
    repo RouteRepo
}

// AssignRoute assigns a custom route to a page. If a route already exists for this page,
// it creates a new version. If the customRoute already exists for another page, returns error.
func (s *Service) AssignRoute(ctx context.Context, page, customRoute string) error
```

---

## Phase 3: Repository Layer

### 3.1 Route Repository
**File**: `repository/route/repository.go` (new)

Implement RouteRepo interface:

1. `Create(ctx, route *Route)`:
   - Check if this page already has a route (GetLatestVersion)
   - If existing route.Route == route.Route (unchanged), return nil (no-op)
   - If changed: increment version (existing.Version + 1), insert new record
   - Ensure route.Route is unique (UNIQUE constraint in DB)

2. `GetByRoute(ctx, route string)`: Query by route column, return the matching Route record

3. `GetLatestVersion(ctx, page string)`: Query highest version for given page key (format: "schema/identifier")

**Transaction handling**: Use `database.InTx` for Create operation to read current version + insert atomically.

---

## Phase 4: Page Service Updates

### 4.1 Route Assignment on Page Save
**File**: `domain/page/service.go`

Modify `Save` method:
- After page is saved, check if custom route was provided (from DTO)
- If provided (non-empty), call `routeService.AssignRoute(pageKey, customRoute)` where `pageKey = schemaName + "/" + identifier`
- If empty (field cleared), no action needed - old routes remain (historical redirects keep working, page accessible via default route)

No new service method needed for getting custom route - pages router can call routeService.GetLatestVersion directly.

---

## Phase 5: Admin UI - Edit Form

### 5.1 DTO Update
**File**: `controller/adminpage/dto.go`

Add single field to `pageDto`:
```go
type pageDto struct {
    // ... existing fields ...
    Route string // Custom route path (empty = use default)
}
```

Update:
- `EnhanceFromForm`: Read `route` form field (c.FormValue("route"))
- `enhanceFromModel`: Fetch custom route via `routeService.GetLatestVersion(schemaName + "/" + identifier)` and populate
- `ToModel`: Pass Route value to page service (not stored in page table)

### 5.2 Controller Update
**File**: `controller/adminpage/controller.go`

In `Edit` (GET):
- Load custom route from route service and pass to template (via DTO.enhanceFromModel)

In `Save` (POST):
- Read `route` from form in `EnhanceFromForm`
- Pass to service.Save via DTO
- After page save, service.Save will call routeService.AssignRoute

### 5.3 Edit Template
**File**: `template/admin/page/edit.hbs`

Add collapsible section (DaisyUI details component):
```hbs
<details class="collapse collapse-arrow bg-base-200" open="{{#if page.Route}}open{{/if}}">
  <summary class="collapse-title text-lg font-medium">
    <i class="fas fa-route mr-2"></i>
    Custom Route
  </summary>
  <div class="collapse-content">
    <div class="form-control">
      <label class="label">
        <span class="label-text">Route Path</span>
        <span class="label-text-alt">Optional custom URL path</span>
      </label>
      <input type="text" name="route" value="{{page.Route}}"
             placeholder="/{{class}}/{{slugify page.secondaryIdentifier}}"
             class="input input-bordered" />
      <label class="label">
        <span class="label-text-alt">
          Leave empty to use default route: /{{class}}/{{slugify page.secondaryIdentifier}}
        </span>
      </label>
    </div>
  </div>
</details>
```

**Note**: Need `slugify` Handlebars helper.

### 5.4 Slugify Helper
**File**: `pkg/url/url.go` (new)

```go
package url

import (
    "regexp"
    "strings"
)

var (
    validSlugCharactersRegex = regexp.MustCompile(`[^a-z0-9]+`)
    repeatingHyphensRegex    = regexp.MustCompile(`-+`)
)

func Slugify(text string) string {
    text = strings.ToLower(text)
    text = validSlugCharactersRegex.ReplaceAllString(text, "-")
    text = strings.Trim(text, "-")
    text = repeatingHyphensRegex.ReplaceAllString(text, "-")
    return text
}
```

Register in `pkg/handlebars/handlebars.go`:
```go
helpers = map[string]any{
    "concat":         concat,
    "beautify":       beautify,
    "use":            use,
    "compareAndUse":  compareAndUse,
    "htmxSortButton": htmxSortButton,
    "slugify":        url.Slugify,  // Add this
}
```

---

## Phase 6: Public Router Updates

### 6.1 Route Resolution Middleware
**File**: `controller/router/router.go` (new function `resolveRoute`)

```go
func resolveRoute(routeSvc route.Service, pageSvc page.Service, c *gin.Context) {
    requestedPath := c.Request.URL.Path

    // 1. Check if path is a custom route
    route, err := routeSvc.GetByRoute(c, requestedPath)
    if err == nil && route != nil {
        pageKey := route.Page // format: "schema/identifier"

        // Check if this is the latest version
        latestRoute, _ := routeSvc.GetLatestVersion(c, pageKey)
        if latestRoute != nil && latestRoute.Route != requestedPath && latestRoute.Version > route.Version {
            // Outdated route - 301 redirect to latest
            c.Redirect(301, latestRoute.Route)
            return
        }

        // Load page (only enabled)
        loadPage(c, pageKey, true)
        return
    }

    // 2. Check if path matches default pattern /:class/:identifier
    // Split path: first segment = class, second = identifier (ULID format)
    parts := splitPath(requestedPath)
    if len(parts) == 2 {
        class, identifier := parts[0], parts[1]

        // Check if page exists and is enabled
        page, err := pageSvc.GetPageBySchemaNameAndIdentifier(c, class, identifier, true)
        if err == nil && page != nil {
            // Check if page has a custom route pointing elsewhere
            customRoute, _ := routeSvc.GetLatestVersion(c, class+"/"+identifier)
            if customRoute != nil && customRoute.Route != requestedPath {
                // Has custom route but accessed default route - 301 redirect to custom
                c.Redirect(301, customRoute.Route)
                return
            }

            // No custom route or already on latest - load page
            loadPage(c, class+"/"+identifier, true)
            return
        }
    }

    // 3. No match - 404
    c.AbortWithStatus(404)
}

func splitPath(path string) []string {
    // Remove leading/trailing slashes, split
    path = strings.Trim(path, "/")
    if path == "" {
        return []string{}
    }
    return strings.Split(path, "/")
}

func loadPage(c *gin.Context, pageKey string, onlyEnabled bool) {
    parts := strings.SplitN(pageKey, "/", 2)
    if len(parts) != 2 {
        c.AbortWithStatus(404)
        return
    }
    class, identifier := parts[0], parts[1]
    // Reuse existing dynamicPageCtrl.LoadPage logic
    // Need to inject dynamicPageCtrl or call it
    // For simplicity in resolver, we could set gin params and call existing handler
}
```

**Better approach**: Instead of duplicating loadPage logic, extract page loading into a service and have resolver set gin params then call `dynamicPageCtrl.LoadPage`. Or have resolver call the service directly.

### 6.2 Router Configuration Updates
**File**: `controller/router/router.go`

Modify `SetPublicRoutes`:

```go
func SetPublicRoutes(router *gin.Engine, svc Services) {
    addCommonHandlers(router, false)
    registerPublicPageHelpers(svc)

    dynamicPageCtrl := dynamicpage_ctrl.NewController(svc.DynamicPageRenderer, svc.Schema, svc.Page)
    previewCtrl := preview_ctrl.NewController(dynamicPageCtrl)

    router.GET("/", func(c *gin.Context) {
        schemaNames, err := svc.Page.GetEnabledSchemaNames(context.Background())
        if err != nil {
            controller.InternalServerError(c, "failed to load page", err)
            return
        }

        if len(schemaNames) == 0 {
            template_ctrl.WithLayout(c, "empty")
            return
        }

        c.Redirect(http.StatusTemporaryRedirect, "/"+schemaNames[0])
    })

    // Schema-level list (e.g., /product)
    router.GET("/:class", dynamicPageCtrl.List)

    // Custom route resolver - catch-all for custom routes and default /:class/:identifier
    // Must be last route in public router
    router.GET("/*path", func(c *gin.Context) {
        resolveRoute(svc.Page, svc.Schema, dynamicPageCtrl, c)
        // Need to pass route service - add to Services struct
    })

    // Preview routes (should come before catch-all)
    router.POST("/preview/:class", previewCtrl.InFlightPage)
    router.GET("/preview/:class/:identifier", previewCtrl.LoadPage)
}
```

**Important**: Add `route.Service` to `Services` struct and inject during server setup.

### 6.3 Preview Integration
**File**: `controller/preview/controller.go` and `controller/dynamicpage/controller.go`

No changes needed. Preview already uses `/preview/:class/:identifier` and includes disabled pages. This bypasses route system intentionally.

---

## Phase 7: Template Updates for Preview

### 7.1 Edit Template Preview Button
**File**: `template/admin/page/edit.hbs`

The existing `submitPreview` JavaScript function needs to use identifier (ULID), not custom route. Should already work as it uses identifier. No change needed unless we discover bug.

### 7.2 List Template Preview Links
**File**: `template/admin/page/list.hbs`

Already uses `/preview/{{class}}/{{this.Identifier}}` - no change needed.

---

## Phase 8: Page Enabled Validation

### 8.1 Enforce Enable Checks
**File**: `resolveRoute` function in `controller/router/router.go`

Already calling `pageSvc.GetPageBySchemaNameAndIdentifier(c, class, identifier, true)` which filters to enabled only. This is correct.

---

## Phase 9: Migration & Data Backfill

### 9.1 Database Schema Update
**File**: `data/db/sqlite/0000_init_schemas.sql`

Simply append the route table creation SQL. No backfill needed - existing pages work via default routes. Route table only stores custom routes.

---

## Phase 10: Testing

### 10.1 Unit Tests
**Skipped for MVP** - will implement later.

### 10.2 Integration Tests
**Skipped for MVP** - will implement later.

### 10.3 Manual Testing Checklist
- [ ] Create page with custom route → loads at custom URL
- [ ] Change custom route → old URL redirects (301) to new
- [ ] Remove custom route (clear field) → page available at default route
- [ ] Access disabled page → 404
- [ ] Custom route uniqueness validation (duplicate route → error)
- [ ] Admin form: collapsible section works, default filled correctly
- [ ] Slugify: "My Product" → "my-product", "Special!@#" → "special"
- [ ] Preview button works with custom routes

---

## Phase 11: Error Handling & Validation

### 11.1 Duplicate Route
In `routeService.AssignRoute`:
```go
// 1. Check if customRoute already exists for ANY page
existing, _ := routeRepo.GetByRoute(ctx, customRoute)
if existing != nil && existing.Page != page {
    return fmt.Errorf("route already in use by another page")
}
// Then check if page already has this route (version check)
```

Admin UI: Show error message from service.Error() in template.

### 11.2 Route Validation
Validate in `AssignRoute` before DB insert:
- Must start with `/`
- No trailing `/` (normalize by trimming)
- At least 2 characters after leading slash (not just `/`)
- Only URL-safe characters after slugify: `[a-z0-9-]` and `/` separators
- Max length 255

If validation fails, return error.

### 11.3 Transaction Safety
`AssignRoute` uses transaction for:
1. Get latest version for page
2. Check uniqueness (route doesn't exist for other page)
3. Insert new version
All in one transaction with serializable isolation (or rely on UNIQUE constraint).

---

## Phase 12: Implementation Order

1. Create slug package + Handlebars registration
2. Database schema update (append to init SQL)
3. Route domain model + repo interface
4. Route repository implementation
5. Route service with AssignRoute logic
6. Add route.Service to Services struct and wire in server
7. Page service integration (call AssignRoute on save)
8. Admin DTO + controller updates (load/save route)
9. Admin template collapsible section
10. Public router: add catch-all with resolveRoute
11. Error handling & validation in route service
12. Manual testing & bug fixes

---

## Files to Modify/Create

### New Files:
- `domain/route/model.go`
- `repository/route/repository.go`
- `pkg/slug/slug.go`

### Modified Files:
- `data/db/sqlite/0000_init_schemas.sql` (append)
- `pkg/handlebars/handlebars.go` (register slugify)
- `domain/page/service.go` (route assignment on save)
- `controller/adminpage/dto.go` (add Route field)
- `controller/adminpage/controller.go` (use DTO.Route)
- `template/admin/page/edit.hbs` (collapsible section)
- `controller/router/router.go` (add route.Service to Services, SetPublicRoutes update, resolveRoute function)
- `server/server.go` (wire route service into Services)
- Possibly `controller/dynamicpage/controller.go` (if resolver uses it)

---

## Edge Cases & Considerations

1. **Slugify conflicts**: Multiple pages can have same slugified secondary identifier. Custom route uniqueness check prevents duplicates globally.
2. **Route changes on disabled pages**: Allowed - route version still increments for SEO redirects.
3. **Trailing slash**: Normalize - store without trailing slash, redirect `/foo/` to `/foo`.
4. **Preview with custom routes**: Preview uses default route, so it previews correctly regardless of custom route setting.
5. **Deletion cascade**: Not implemented via FK (since page is not a real FK), but when page deleted, route service should cleanup routes (or rely on manual cleanup). For MVP: no cascade, orphaned routes return 404.
6. **Performance**: Simple UNIQUE index lookup on route table should be fast. No caching needed for MVP.

---

## Success Criteria

- Admin can set custom route in collapsed section (default placeholder shows slugified secondary identifier)
- Custom routes stored with version increment on changes
- Public routing:
  - Custom routes load correct page (enabled check)
  - Old custom route versions 301 redirect to latest
  - Default `/:class/:identifier` routes work when no custom route set
  - Default route 301 redirects to custom if custom exists
  - Disabled pages return 404 for both default and custom routes
- Duplicate custom routes are rejected with user-friendly error
- Slugify helper works in templates
- All manual testing checklist items pass

---

This plan provides a complete implementation strategy for MVP that maintains backward compatibility while adding flexible custom routing with proper SEO handling via 301 redirects. Key simplifications: merged page key, minimal repository interface, skipped unit tests, no caching, and direct database schema update.
