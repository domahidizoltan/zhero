# Security and code review guidelines

- Never hardcode or commit any secrets, passwords or keys.
- Extract them to `config.yaml` and `config/config.go` with some dummy values.
- Validate all user input.
- Use enum validation methods.
- Sanitize data before database operations.
- Never log passwords, tokens, or secrets.
- Never concatenate user input into database queries.