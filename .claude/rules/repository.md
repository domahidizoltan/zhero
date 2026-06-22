---
paths:
  - "repository/**/*.go"
---

# Repository Rules

- Always use parameterized queries.
- Extract SQL statements to constants.
- Have short method names where the context is given by the repo package name (e.g. `page.List`, `schema.Upsert`)
