---
paths:
  - "**/*.go"
---

# Backend coding guidlines

- Follow [Effective Go](https://go.dev/doc/effective_go).
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
- Shared functionalities under `pkg/` must be public functions, but can have private functions to hide details.
- Place domain entities to `domain/<domain>/model.go`; they are public.
- Place domain services to `domain/<domain>/service.go`.
- Place controllers to `controller/<domain>/controller.go`.
- Controller methods represents actions (e.g. `page.List()`).
- Place DTOs to `controller/<domain>/dto.go`.
- Keep DTOs private.
- Make DTOs have meaningful receiver methods when it owns the data (e.g. `func (dto *pageDto) toModel() page_domain.Page {`)
- Use Go's context for request-scoped values and cancellations.
- Defer closing resources to avoid leaks.
- Use pointer receivers for methods that modify state.
- Use value receivers for methods that don't modify state.
