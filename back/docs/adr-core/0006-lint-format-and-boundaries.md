# 6. Lint, format and boundaries with golangci-lint

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

ADR common 0002, 0003 and 0005 require a standard linter, a standard formatter and a
boundary checker.

## Decision

- **golangci-lint** with its standard linters plus `go vet`, no bespoke rules.
- **gofmt** and **goimports**, run through golangci-lint's formatters.
- **depguard** rules in the same configuration enforce the boundaries: no import between
  contexts, `core` never imports `bc`, `shared` imports neither, `domain` imports no
  infrastructure or framework package.
- Generated code (`apigen`) is excluded from linting.

## Consequences

- One tool and one configuration file for lint, format and architecture rules.
- `make lint` is the single quality gate of the back.
