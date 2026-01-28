# Development Plan - Public Facing Website Extension

This plan outlines the steps to extend the public-facing server with schema navigation, page listing, and individual page detail views.

## 1. Backend Changes

### 1.1. Schema Discovery for Menu
- **Service Layer:** Add a method `GetEnabledSchemaNames` to `domain/schema/service.go`. This method will query for schema names that have at least one enabled page.
- **Controller Layer:** Create a new handler `PublicSchemaNames` in `controller/schema/controller.go`. This handler will utilize `GetEnabledSchemaNames` and return a JSON list of schema names.
- **Server Configuration:** Register a new public API route in `server/server.go` (e.g., `/api/schemas`) that maps to the `PublicSchemaNames` handler.

### 1.2. Public Page Listing by Schema
- **Controller Layer:** Implement a new handler `PublicPageList` in `controller/page/controller.go` to manage the `/:schemaName` route.
- **Service Layer:** This handler will use `domain/page/service.go` to fetch all *enabled* pages for the specified schema.
- **Templating:** Render a new template `template/page/public_list.hbs` to display the list of pages.

### 1.3. Public Individual Page Detail
- **Controller Layer:** Implement a new handler `PublicPageDetail` in `controller/page/controller.go` to manage the `/:schemaName/:identifier` route.
- **Service Layer:** This handler will use `domain/page/service.go` to fetch a single *enabled* page by its schema name and primary identifier.
- **Templating:** Render a new template `template/page/public_detail.hbs` to display the full content of the page.

## 2. Frontend Changes

### 2.1. Menu Integration
- **Layout Template:** Modify the main public layout template (e.g., `template/layout.hbs`) to:
    - Fetch schema names from the `/api/schemas` endpoint.
    - Dynamically generate the navigation menu using these schema names.
    - Each menu item will link to the public page listing route (e.g., `/:schemaName`).

### 2.2. New Templates
- **`template/page/public_list.hbs`:** Create this new Handlebars template to display a list of pages for a given schema. Styling should be consistent with `template/temp_body.html`.
- **`template/page/public_detail.hbs`:** Create this new Handlebars template to display the detailed content of a single page.

## 3. Error Handling
- Non-existent schemas or pages will result in a `404 Not Found` HTTP response.

## 4. Styling
- Existing Tailwind CSS and DaisyUI classes will be used for consistency. No custom CSS will be introduced unless absolutely necessary.
