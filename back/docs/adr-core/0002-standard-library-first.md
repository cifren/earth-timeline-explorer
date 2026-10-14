# 2. Standard library first

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Each dependency is code to audit, upgrade and learn. Go's standard library covers most
needs of a small HTTP service: routing with methods and path parameters, JSON, HTTP
client, structured logs, context.

## Decision

Use the standard library by default: `net/http` (server, router, client), `encoding/json`,
`log/slog`, `context`, `errors`. A third-party module is added only when the standard
library has no equivalent, and is wrapped in `infrastructure` or `shared`.

Accepted third-party modules: the OpenAPI code generator and its runtime (ADR back core
0003), the Redis client (ADR back functional 0002).

## Consequences

- Few dependencies, small binary, easy upgrades.
- No web framework to learn or follow.
- Some glue code is written by hand (middleware, error responses), kept in `shared`.
