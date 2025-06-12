# Codex Guidelines

This project uses Codex to help automate coding tasks. Follow these best practices when creating pull requests:

## Workflow
- Format Go code using `go fmt ./...` before committing.
- Run `go vet ./...` to detect common issues.
- Execute `make test` to ensure all tests pass.

## Commit Messages
- Keep commits focused and descriptive.
- Use the imperative mood ("Add feature" not "Added feature").

## Pull Requests
- Summarize the change and reference relevant issues if any.
- Include a **Testing** section showing the commands run and their results.

Feel free to update this file with additional project-specific instructions.
