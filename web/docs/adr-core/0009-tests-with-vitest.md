# 9. Tests with Vitest and Testing Library

Date: 2026-10-15
Status: Accepted

## Context

ADR common 0009 requires unit and integration tests, fast and isolated from external
services. The application is built with Vite.

## Decision

- **Vitest** runs every test, sharing the Vite configuration (aliases, TypeScript).
- **React Testing Library** tests components through what the user sees and does, not
  through their internals.
- **Unit tests** (`*.test.ts`, `*.test.tsx` beside the code) cover mappers, hooks and
  single components.
- **Integration tests** (`*.integration.test.tsx`) render a whole page and drive it like
  a user, with the back's API faked at the HTTP level by **Mock Service Worker**, so
  gateways and mappers run unchanged.
- `make test` runs the unit tests; `make test-integration` runs the integration tests.

## Consequences

- Tests reuse the build configuration: no separate setup to keep in sync.
- Page tests survive refactoring as long as the behaviour is unchanged.
