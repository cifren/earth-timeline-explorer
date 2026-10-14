# 4. Typed domain errors, Problem Details at the edge

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Errors cross several layers: a source timeout, an unknown zone, an invalid decade. Each
must reach the client in a standard format without leaking internals.

## Decision

- `domain` and `application` return **typed errors** (sentinel errors or error types),
  wrapped with `%w` to keep the cause; callers inspect them with `errors.Is` / `errors.As`.
- Only `presentation` translates errors into HTTP responses, as **Problem Details**
  (RFC 9457, `application/problem+json`): status, `type`, `title`, `detail`.
- Unexpected errors become a generic `500` problem; the cause is logged, never sent.
- An empty result is not an error (e.g. no image for a zone and decade returns an empty
  list).

## Consequences

- One error format for every client.
- Business code never knows about HTTP status codes.
