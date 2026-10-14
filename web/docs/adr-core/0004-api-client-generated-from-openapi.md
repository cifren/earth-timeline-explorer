# 4. API client generated from OpenAPI

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The API contract `back/api/openapi.yaml` is the source of truth (ADR functional 0003). Hand-written
request types and calls would drift from it.

## Decision

Types and HTTP calls are generated with **Orval** as plain typed `fetch` functions (no
generated hooks: queries belong to the application layer, ADR web core 0003).

- Generated code lives in `src/shared/api-client/`, split per OpenAPI tag, committed,
  never edited; excluded from lint and formatting, kept under strict type-check.
- Generated types are wire types: hand-written mappers in each context's
  `infrastructure` convert them into domain types; generated types never reach `domain`
  or `presentation`.
- Regeneration is the `api` Makefile target.

## Consequences

- A contract change that the web does not follow fails to compile.
- The domain stays independent of the wire format.
- The same generation can serve the mobile application.
