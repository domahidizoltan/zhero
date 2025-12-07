# Development instructions
- Use `go mod tidy` to update the project dependencies.
- Use `go build` to check if any new change does not create compile time issues across the project.
- Use Context7 to get up-to-date technical documentations.

# Testing instructions
- Use `go test` to run unit tests.

# Development Workflow instructions
- TODO

# Agentic Coding Workflow instructions
- Refresh the `main` branch and create a new branch.
- Get details about new change request from the user.
- Use the change request details with `PRD.md` contents and create a plan into `PLAN.md`.
- Ask questions about the change to make sure the problem is fully understood and the best plan will be created.
- If needed create features in `PLAN.md` and use identifiers with format `F012`.
- Use `RULES.md` and `PLAN.md`, think hard, and create tasks into `TASKS.md` file:
  - Create an identifier for each task with format `T012`.
  - Create a checkbox in the front of the task title what indicates the task status.
  - Add a legend to the file for the following task status: ready to implement, in progress, completed, blocked.
  - Add a dependency field, where tasks are enumerated which has common files to change and can cause conflicts.
  - Add an implementation comment field, where implementation blockers can be summarised.
- Implement the tasks in the order requested by the user.
- Update the status of the tasks after implementation.
- Cover each change with tests.
- Get new dependencies and check for compile or runtime errors.
- In case of a bug fix, first add the tests.
- Review the changes and notify the user about potential security issues, technical debts and improvement possibilities.
- NEVER commit anything to Git; ask user to review and commit changes.

## TASKS.md template
- Add this status legend to the top of the file: ~ in progress, x completed, ! failed or blocked
- Use the statuses in the task title checkboxes.
- Update the task statuses after completing each task.
- Use this template below for the tasks (parts between {} are instructions for the field, don't copy them):
```
# Task List

## F012: Feature title

- [ ] T001: **Task Title**
  - Details: {Here come the implementation details.}
  - Dependencies: T000, T001
  - Comment: {Post-implementation comments about failures or impediments.}

```
