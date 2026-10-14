# 5. Context and timeouts everywhere

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The back mostly waits on other systems (Commons, Redis). Without deadlines, a slow
dependency piles up requests and takes the service down.

## Decision

- Every function doing I/O takes a `context.Context` as its first parameter and passes it
  down.
- Every outbound call has a **deadline** (e.g. 3 seconds for a Commons call), derived
  from the incoming request's context.
- The HTTP server sets read, write and idle timeouts, and shuts down gracefully on
  `SIGTERM`, letting in-flight requests finish.

## Consequences

- A slow dependency degrades one request, not the service.
- A client disconnecting cancels the work done for it.
