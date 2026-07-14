---
paths:
  - "**/*_test.go"
---

# Testing Rules

- THIS PROJECT IN MVP PHASE, IGNORE TESTS FOR THE MOMENT.
- Write tests using Golang stdlib, Testify assertions and Uber Gomock.
- Every public function MUST have unit tests.
- Use table-driven tests for multiple cases.
- TODO mocks
- Regenerate mocks when interfaces change.
- Test repositories with integration tests.