# 3. Server interfaces generated from OpenAPI

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The API contract `api/openapi.yaml`, at the root of the back, is the source of truth (ADR functional 0003). Hand-written
request and response types would drift from it.

## Decision

Server types and interfaces are generated with **oapi-codegen** in **strict server**
mode: one typed method per operation, typed request and response objects.

- Generated code lives in `internal/shared/apigen`, committed, never edited.
- Each context's `presentation` implements the operations of its OpenAPI tag; the
  composition root assembles them into the single generated interface.
- Generated types are wire types: hand-written mappers in `presentation` convert them
  to and from domain models; generated types never reach `domain`.
- Regeneration is the `api` Makefile target.

## Consequences

- A contract change that the code does not follow fails to compile.
- Handlers contain no parsing or serialization code.
- The domain stays independent of the wire format.
