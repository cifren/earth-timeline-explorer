# 7. Path aliases for imports

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Folder-by-context plus folder-by-layer produces deep trees; relative imports such as
`../../../../shared/...` are noisy and hide which boundary they cross.

## Decision

TypeScript path aliases, mirrored in Vite: `@app/*`, `@core/*`, `@shared/*`, and
`@bc/<context>` pointing at each context's public `index.ts`. Imports inside a context's
own folder stay relative.

## Consequences

- Imports read by intent and survive file moves.
- Alias roots line up with the boundaries checked by the linter.
