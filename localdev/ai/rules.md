#### Project Structure
```
/project-root
├── /cmd                  # Main applications for the project
│   └── /your-app
├── /internal             # Private application and library code
│   ├── /domain1          # Grouped by domain
│   ├── /domain2
│   └── /domain3
├── /pkg                 # Public libraries
├── /api                 # API definitions
├── /web                 # Frontend files
│   ├── /htmx
│   └── /js
├── /migrations          # Database migrations
├── /configs             # Configuration files
└── /tests               # Test files
```

#### Technical stack
- Backend: Golang | [Gin](https://github.com/gin-gonic/gin)
- Frontend: HTML5/HTMX/JavaScript | [Tailwind CSS](https://tailwindcss.com/)
- Database: SQL (separate databases for dev, test, prod)
- Testing: Golang tests, uber-go/mock for mocking,  [Playwright](https://playwright.dev/) for E2E testing

#### General Responsibilities:
- Guide idiomatic, maintainable, and high-performance Go code.
- Enforce modular design and separation of concerns.
- Promote test-driven development, observability, and scalability.


#### Coding Standards
- Follow [Effective Go](https://go.dev/doc/effective_go).
- Use gofmt for formatting; limit line length to 80 characters.
- Use descriptive names for variables, functions, and packages.

#### Libraries and Frameworks
- Use [Gin](https://github.com/gin-gonic/gin) for the backend framework.
- Use [Tailwind CSS](https://tailwindcss.com/) for styling the frontend.
- Utilize `uber-go/mock` for mocking in tests.

#### Development best practices
- Write short, focused functions with a single responsibility.
- Check and handle errors explicitly using wrapped errors.
- Avoid global state; use constructor functions for dependency injection.
- Use Go's context for request-scoped values and cancellations.
- Safely use goroutines; guard shared state with channels or sync primitives.
- Defer closing resources to avoid leaks.
- Implement retries and timeouts for external calls.
- Favor simple solutions and avoid code duplication.
- Ensure code is environment-aware (dev, test, prod).
- Avoid introducing new patterns or technologies without necessity.
- Keep the codebase clean, organized, and avoid scripts in files.
- Refactor files exceeding 500-600 lines.
- Mock data only for tests; never for dev or prod environments.
- Do not overwrite .env files without confirmation.
- Avoid overdocumenting; comment only when necessary.

#### Database management
- Use parameterized queries or ORM to safeguard against SQL injection.
- Never execute DB write operations (create, update, delete) without any approval.
- Never delete data from datasets (yaml, json) or entire files without any approval.

#### Testing Practices
- Use Go's testing package for backend tests; write unit tests for public functions.
- Use Playwright for frontend testing; employ table-driven tests and parallel execution.
- Separate fast unit tests from slower integration and E2E tests.
- Ensure test coverage for exported functions.

#### Security Guidelines
- Validate all user inputs to prevent invalid data from entering the system.
- Implement authentication and authorization using [Casbin](https://casbin.org/).
- Never leak sensitive information; use secure defaults for JWT and cookies.
- Isolate sensitive operations with clear permission boundaries.

#### Concurrency and Goroutines
- Ensure safe use of goroutines, and guard shared state with channels or sync primitives.
- Implement goroutine cancellation using context propagation to avoid leaks and deadlocks.

#### Key Conventions
- Prioritize readability, simplicity, and maintainability.
- Design for change: isolate business logic and minimize framework lock-in.
- Emphasize clear boundaries and dependency inversion.
- Ensure all behavior is observable, testable, and documented.
- Automate workflows for testing, building, and deployment.

#### Frontend Guidelines
- Use semantic HTML5 elements with HTMX attributes.
- Implement CSRF protection.
- Utilize HTMX extensions and hx-boost for navigation.

#### Self-Improvement Mechanism
- Periodically summarize potential improvements to the rules file for approval during code reviews.
- Implement a feedback mechanism (thumbs-up/thumbs-down) for AI suggestions.

#### Additional Notes
- Encourage AI to suggest optimizations based on code complexity.
- Keep the rules file dynamic for evolving best practices.
- Focus on relevant code areas; avoid unrelated changes.
- Write thorough tests for major functionality.
- Avoid major changes to established patterns unless instructed.
- Consider potential impacts on other code areas when making changes.
- When code examples, setup or configuration steps, or library/API documentation are requested, try to use context7.


#### Do not repeat the following mistakes
TODO

