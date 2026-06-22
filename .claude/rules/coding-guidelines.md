# Coding guidlines

- Keep responses and modifications minimal; avoid speculative changes.
- Focus on relevant code areas; avoid unrelated changes.
- Analyze patterns in the codebase and try to make the new changes fit into existing patterns.
- Write short, focused functions with a single responsibility.
- Avoid overdocumenting; comment only when necessary.
- Explain "why" not "what" for complex logic.
- Avoid global state; use constructor functions for dependency injection.
- Ensure code is environment-aware (dev, test, prod).
- Use Context7 to get the most up-to-date documentations.- Avoid introducing new patterns or technologies without necessity.
- File names support fuzzy search (e.g. `page/repository.go`, `controller/router`, `schema/controller.go`, `page/model.go`, `schema/service.go`, `page/edit.hbs`).