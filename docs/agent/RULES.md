Coding and Development Rules


# 1. Project Structure
The project follows a layered architecture to ensure separation of concerns.
```sh
/
├── config/             # Configuration loading and object definitions.
├── controller/         # Gin controllers, handles HTTP requests, DTOs, and form binding.
│   ├── page/           # Controller and DTO for Page domain.
│   ├── router/         # Gin routing configuration.
│   ├── schema/         # Controller and DTO for Schema domain.
│   ├── template/       # Shared management functions.
│   └── controller.go   # Shared controller functions.
├── data/db/sqlite/     # SQLite database migration scripts.
├── docs/agent/         # AI agent related files.
├── domain/             # Core business logic, models (structs), and service interfaces.
│   ├── page/           # Service and Model for Page domain.
│   ├── schema/         # Service and Model for Schema domain.
│   └── schemaorg/      # Service and Model for Schema.org management.
├── pkg/                # Shared, reusable packages (e.g., database, logging).
├── repository/         # Data access layer (repositories) for domain Models.
├── template/           # Handlebars templates (.hbs) and static assets (.js, .css).
│   └── template.go     # Template name constant definitions.
├── config.yaml         # Appliction configuration file.
└── main.go             # Application entry point and dependency injection setup.
```


# 2. Technical Stack
- **Backend:** [Go 1.24+](https://go.dev/), [Gin Web Framework](https://gin-gonic.com/)
- **Frontend:** [HTMX](https://htmx.org/), JavaScript, [Tailwind CSS](https://tailwindcss.com/), [DaisyUI CSS Framework](https://daisyui.com/)
- **Templating:** [Handlebars](https://github.com/aymerick/raymond) (via `aymerick/raymond`)
- **Database:** [SQLite](https://sqlite.org/index.html) (via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite))
- **Data Formats:** [JSON-LD](https://json-ld.org/) for structured data, RDF for [Schema.org](https://schema.org/) graph processing.
- **Logging:** [Zerolog](https://github.com/rs/zerolog)
- **Configuration:** [Viper](https://github.com/spf13/viper)


# 3. General Responsibilities
- Adhere to the existing architecture (Controller -> Service -> Repository).
- Maintain the distinction between domain models and controller DTOs.
- Ensure frontend interactions are primarily driven by HTMX where possible.
- Guide idiomatic, maintainable, and testable Go code.
- Enforce modular design and separation of concerns.
- Follow [Effective Go](https://go.dev/doc/effective_go).
- Keep responses and modifications minimal; avoid speculative changes.
- Focus on relevant code areas; avoid unrelated changes.
- Analyze patterns in the codebase and try to make the new changes fit into existing patterns.


# 4. Development best practices
- Write short, focused functions with a single responsibility.
- Avoid overdocumenting; comment only when necessary.
- Explain "why" not "what" for complex logic.
- Avoid global state; use constructor functions for dependency injection.
- Use Go's context for request-scoped values and cancellations.
- Defer closing resources to avoid leaks.
- Ensure code is environment-aware (dev, test, prod).
- Avoid introducing new patterns or technologies without necessity.
- Use pointer receivers for methods that modify state.
- Use value receivers for methods that don't modify state.
- Use Context7 to get the most up-to-date documentations.


# 5. Development guidelines
- File names support fuzzy search (e.g. `page/repository.go`, `controller/router`, `schema/controller.go`, `page/model.go`, `schema/service.go`, `page/edit.hbs`).

## 5.1 Frontend coding guidelines
- Define templates in `template/<domain>/<action>.hbs` files; `<action>` is for responsibility (e.g. edit or list).
- Use kebab-case format for HTML tag IDs and class names.
- Add template names to `template/template.go` for backend reuse.
- Use FontAwesome icons; try to avoid generating SVG codes for icons.
- Make responsible layouts.
- Use theme classes what supports dark and light themes.
- Keep minimal JavaScript in `.hbs` files; extract JavaScript functions to `template/<domain>/<domain>.js` files.
- Use native and DaisyUI validators for HTML forms.
- Don't create custom HTML components; use native and DaisyUI components.
- Split complex templates into partials.

## 5.2 Backend coding guidelines
- Interfaces are defined on the consumer side (e.g. service defines the repository interface before injection).
- Interfaces should be private.
- Organize stdlib imports first, and dependency imports after.
- When import must be aliased, use snake_case in reverse order (e.g. `page_domain "../domain/page"`) 
- Add package documentation comment to the main file in package or to the file `<package>.go`.
- Add a `<package>.go` file to the package to store shared constants and functions.
- Controller and Repository layer depends on Domain layer, to avoid circular dependencies.
- Controller, Repository and Domain layer contains same subpackages for the given domains (e.g `controller/page`, `domain/page`, `repository/page`).
- Use package and variable names to give context, and method names to extend context (e.g. `page.NewService()`, `schema.NewRepo()`, `pageRepo.Insert()`).
- Move shared functionalities into `pkg/<functionality>/<functionality>.go`.
- Shared functionalities under `pkg/` must be private functions, but can have private functions to hide details.
- Place domain entities to `domain/<domain>/model.go`; they are public.
- Place domain services to `domain/<domain>/service.go`.
- Place controllers to `controller/<domain>/controller.go`.
- Controller methods represents actions (e.g. `page.List()`).
- Place DTOs to `controller/<domain>/dto.go`.
- Keep DTOs private.
- Make DTOs have meaningful receiver methods when it owns the data (e.g. `func (dto *pageDto) toModel() page_domain.Page` {

## 5.3. Data access guidelines
- All database access MUST go through repository interfaces.
- Test repositories with integration tests.
- Always use parameterized queries.
- Use query builders when possible.
- Run write operations within transaction.
- Extract SQL statements to constants.
- Have short method names where the context is given by the repo package name (e.g. `page.List`, `schema.Upsert`)

### 5.3.1 Database migrations
- Place all migrations in `data/db/sqlite/`.
- Each feature has a separate migration script with the file with pattern `YYMMDD_01_short_name`.
- Number migrations sequentially.
- Add the changes to the `sqlite.go` file.
- Never modify existing migrations.

## 5.4. Testing guidelines
- IGNORE TESTS FOR THE MOMENT, JUMP TO THE NEXT SECTION.
- Write tests using Golang stdlib, Testify assertions and Uber Gomock.
- Every public function MUST have unit tests.
- Use table-driven tests for multiple cases.
- TODO mocks
- Regenerate mocks when interfaces change.

- ## 5.5. Security and code review guidelines
- Never hardcode or commit any secrets, passwords or keys.
- Extract them to `config.yaml` and `config/config.go` with some dummy values.
- Validate all user input.
- Use enum validation methods.
- Sanitize data before database operations.
- Never log passwords, tokens, or secrets.
- Never concatenate user input into database queries.

