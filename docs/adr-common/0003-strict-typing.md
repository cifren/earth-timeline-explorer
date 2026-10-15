# 3. Strict typing, no escape hatch

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Permissive typing lets null dereferences and untyped values reach runtime. Data crossing
the applications is meant to be typed end to end; escape hatches break that chain.

## Decision

- Every application runs its compiler in its **strictest** mode.
- No untyped escape hatch (`any` in TypeScript, `interface{}` / `any` in Go) in
  application code; unknown data is narrowed or decoded into a declared type at the
  boundary.
- Gaps in third-party types are typed once, at the boundary.

## Consequences

- Type errors fail the build instead of production.
- Data entering the system is validated where it enters, once.
