---
paths:
  - "**/*_test.go"
---

# Testing Rules

- MVP pragmatism: NO blanket coverage requirement. "Every public function must have unit tests"
  does NOT apply to this project.
- Add a test ONLY where it buys fast, deterministic feedback on subtle logic that documentation
  warns breaks silently:
  - reference token round-trip + `RefPattern`/`ShortRefPattern` regexes
  - four-projection derivation in `pageDto.ToModel`
  - `determineComponent` type/name → component mapping
- If you change reference regexes, projection derivation, or the component mapping, extend the
  matching test in the same change.
- Style: stdlib `testing` + testify assert (already in go.mod); table-driven for mapping cases.
- Constraints: no network, no external services, no new dependencies. No mock framework exists
  in this repo and none should be added — test pure functions and DTO transforms directly.
  `go test ./...` must stay in the seconds range.
- Existing tests (`domain/schemaorg`, `pkg/url`) stay; do not delete them.