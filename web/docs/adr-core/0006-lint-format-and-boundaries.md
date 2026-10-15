# 6. Lint, format and boundaries with ESLint and Prettier

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

ADR common 0002 and 0005 require a standard linter with recommended presets, a standard
formatter and a boundary checker.

## Decision

- **ESLint** with recommended presets only: `typescript-eslint` (type-checked),
  `eslint-plugin-react-hooks`, `eslint-plugin-jsx-a11y`. Their `no-explicit-any` rule
  enforces ADR web core 0005.
- **Prettier** is the only formatter; ESLint does not check formatting.
- **eslint-plugin-boundaries** enforces the architecture: no import between contexts,
  contexts imported only through their `index.ts`, `core` never imports a context,
  `shared` imports neither, layers point inward.
- Generated code (`var/`) is excluded from lint and formatting.

## Consequences

- `make lint` is the single quality gate of the web.
- Accessibility smells are caught as early as type errors.
- ESLint follows the newest major supported by every preset: `eslint-plugin-jsx-a11y`
  holds it on ESLint 9.
- Boundaries resolve the path aliases through `eslint-import-resolver-typescript`.
