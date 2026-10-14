# 6. Concise code, no pass-through abstractions

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Layers and contexts invite ceremony: objects that only forward a call, nested conditionals,
duplicated logic. They hide where the real logic lives and inflate the code base.

## Decision

- **No pass-through**: an object, function or layer exists only if it carries a
  **sequence** (several calls), a **decision** (guard, branch, error translation) or a
  **non-trivial transformation**. Otherwise the caller uses the underlying code directly.
  Example: a use case whose body only returns `repository.load()` is removed.
- **Reuse first**: use what exists in the project or the standard library before writing
  new code.
- **Guard clauses** (early return) instead of `if/else` nesting.
- **SOLID and DRY** as the default design checks.

## Consequences

- Every layer signals real logic; reading the code shows where decisions are made.
- Adding logic later (a guard, a second call) is the trigger for introducing a wrapper.
- "Concise" remains a review judgement, backed by these rules.
