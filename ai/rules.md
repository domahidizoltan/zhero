# Zhero Project - AI Coding Assistant Guidelines

#### Project Structure
The project follows a layered architecture to ensure separation of concerns.
```
/
├── /controller      # Gin controllers, handles HTTP requests, DTOs, and form binding.
├── /domain          # Core business logic, models (structs), and service interfaces.
├── /repository      # Data access layer (repositories), interacts with the database.
├── /pkg             # Shared, reusable packages (e.g., database, logging).
├── /template        # Handlebars templates (.hbs) and static assets (.js, .css).
├── /config          # Configuration loading and structs.
├── /data/db/sqlite  # SQL database migration scripts.
└── main.go          # Application entry point and dependency injection setup.
```

#### Technical Stack
- **Backend:** Go 1.24+ | [Gin](https://github.com/gin-gonic/gin)
- **Frontend:** HTML5 | [HTMX](https://htmx.org/) | JavaScript | [Tailwind CSS](https://tailwindcss.com/) with [DaisyUI](https://daisyui.com/)
- **Database:** SQLite (via `modernc.org/sqlite`)
- **Templating:** [Handlebars](https://handlebarsjs.com/) (via `aymerick/raymond`)
- **Logging:** [Zerolog](https://github.com/rs/zerolog)
- **Configuration:** [Viper](https://github.com/spf13/viper)

#### General Responsibilities:
- Adhere to the existing architecture (Controller -> Service -> Repository).
- Maintain the distinction between domain models and controller DTOs.
- Write clean, idiomatic, and maintainable Go code.
- Ensure frontend interactions are primarily driven by HTMX where possible.

#### Coding Standards & Best Practices
- **Go:**
    - Follow [Effective Go](https://go.dev/doc/effective_go). Use `gofmt` for formatting.
    - Use constructor functions for dependency injection (e.g., `NewService(repo)`).
    - Handle errors explicitly. Use the helper functions in `controller/controller.go` for HTTP error responses.
    - Use `database.InTx` for all database write operations to ensure atomicity.
    - `context.Context` should be propagated from the controller down to the repository layer.
- **Frontend:**
    - Use DaisyUI classes for styling to maintain a consistent look and feel.
    - Use HTMX for server-driven UI updates. Keep client-side JavaScript minimal and for UI enhancements only (e.g., TomSelect, SortableJS).
    - Place JavaScript related to a specific feature in its own file within `template/[feature]/`.
- **Templates:**
    - Use Handlebars templates (`.hbs`).
    - Leverage existing helpers in `pkg/handlebars/handlebars.go` for common formatting tasks.
    - Split complex templates into partials.

#### Database Management
- Write new database schema changes as migration files in `data/db/sqlite/`.
- Use parameterized queries to prevent SQL injection.
- The repository layer is the only place where database queries should be executed.

#### Testing (do not add tests for the moment)
- Add unit tests for new business logic in services.
- Follow existing testing patterns using `stretchr/testify`.

#### Do not repeat the following mistakes
- TODO: Add items here as they are discovered.

