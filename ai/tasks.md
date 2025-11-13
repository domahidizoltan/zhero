# Page Management Tasks

## F01: Page Listing

### Backend

- [x] **T00001: Repository:** Implement `List` in `repository/page/repository.go` to fetch paginated and searchable page data (identifier, secondary identifier) for a given schema.
- [x] **T00002: Service:** Add a `List` method to `domain/page/service.go` that calls the repository.
- [ ] **T00003: Controller:** Update the `List` method in `controller/page/controller.go` to:
    - Call the new service method.
    - Handle query parameters for searching, sorting, and pagination.
    - Pass the list of pages and pagination data to the template.

### Frontend (`template/page/list.hbs`)
- [x] **T00004: Dynamic Rendering:** Replace the hardcoded table with a Handlebars `{{#each}}` loop to dynamically render pages.
- [ ] **T00005: Action URLs:** Ensure the "Edit", "Enable/Disable", and "Delete" links have correct, dynamically generated URLs (e.g., `/page/edit/{{class}}/{{identifier}}`).
- [ ] **T00006: Search & Sort:** Use HTMX to make the search input and sortable table headers trigger requests to the `List` endpoint with the appropriate query parameters.
- [x] **T00007: Pagination:** Implement dynamic pagination controls that reload the page list.

## F02: Page Editing & Creation

### Backend (`controller/page/controller.go`)
- [x] **T00008: Save Logic:** Refine the `Save` method to ensure it robustly handles both creation and updates.
- [x] **T00009: Redirects:** Correct the redirect on successful save to point to the list page for the current schema (e.g., `/page/list?=schema={{class}}`).
- [x] **T00010: Validation:** Implement comprehensive server-side validation for submitted page data and return clear error messages to the user. (client side validation is enough for the moment)

### Frontend (`template/page/edit.hbs`)
- [ ] **T00011: Dynamic Components:** Extend the template to render different HTML form components (e.g., `textarea`, `select`, `checkbox`) based on the `Component` type defined in the schema for each field.
- [x] **T00012: Button Logic:**
    - [x] Dynamically change the submit button text to "Create" for new pages and "Update" for existing ones.
    - [x] Fix the "Cancel" button's `onclick` handler to redirect to `/page/list?schema={{class}}`.
- [x] **T00013: Client-side Validation:** In `template/page/page.js`, add JavaScript to handle basic client-side validation for mandatory fields to provide instant feedback.

## F03: Additional Actions & Refinements

- [ ] **T00014: Enable/Disable Functionality:**
    - [x] Implement the backend logic (routes, controller methods, service, repository) to toggle a page's `IsEnabled` status.
    - [x] Update the UI in `template/page/list.hbs` to show the current status and provide the correct toggle action icon/link.
- [ ] **T00015: Delete Functionality:**
    - [ ] Implement the full backend stack for deleting a page.
    - [ ] Add a confirmation modal (using the existing `popup` function in `index.js`) to prevent accidental deletion.
- [ ] **T00016: Preview Functionality:**
    - [ ] Create a new route and controller to render a preview of a page. This could display the page content and the generated JSON-LD structured data.
