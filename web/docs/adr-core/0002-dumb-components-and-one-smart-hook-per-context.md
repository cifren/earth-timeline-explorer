# 2. Dumb components and one smart hook per context

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Components that load data, hold state and render grow heavy and cannot be reused.
Splitting everything into containers and views adds plumbing where a component is
naturally cohesive.

## Decision

- **Dumb component**: receives props, emits callbacks, holds only local UI state (hover,
  open/closed), contains no business rule and no data fetching.
- **`shared/` components are always dumb**, which makes them safely reusable.
- **Inside a context, dumb is the default**, not an obligation: a cohesive component may
  call its context's hook directly.
- **One smart hook per context** (`application/`) exposes the context's state and intents
  to its components, e.g. `useImagery()` returning the current image, loading state and
  `renew()`.
- A hook that only re-exposes another hook is a pass-through (ADR common 0006).

## Consequences

- Business orchestration has one obvious place per context.
- `shared/` stays reusable and checkable by the boundary linter.
