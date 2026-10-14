# 5. Strict TypeScript configuration

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

ADR common 0003 requires the strictest compiler mode and no untyped escape hatch.

## Decision

- `tsconfig.json` enables `"strict": true` plus `noUncheckedIndexedAccess`,
  `noImplicitOverride` and `noFallthroughCasesInSwitch`.
- No `any`: `unknown` narrowed at the boundary, or type parameters; enforced by the lint
  rule `no-explicit-any` (ADR web core 0006).

## Consequences

- Null dereferences, implicit `any` and unchecked indexing are compile errors.
