# 3. Contract first with OpenAPI

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Back, web and mobile must agree on every field of the API. Hand-written types on each side
compile independently and fail at runtime when they drift.

## Decision

The API is described first, in one OpenAPI document owned by the back,
`back/api/openapi.yaml`, which is the **single source of truth** for the wire contract:

- the back generates its server interfaces from it;
- each client generates its types and HTTP calls from it;
- generated code is committed and never edited by hand;
- regeneration is a deliberate command, not a build step: the compile errors that follow
  are the drift detection.

API conventions: version in the URI (`/api/v1/...`), errors as Problem Details
(RFC 9457), HTTP caching headers on reference data (RFC 9111), retirement announced with
the `Deprecation` header (RFC 9745).

## Consequences

- Contract drift becomes a compile error instead of a production bug.
- The contract is reviewed as a document before code is written.
- The contract's quality (tags, `required` fields, enums) directly shapes the generated
  types on every side.
