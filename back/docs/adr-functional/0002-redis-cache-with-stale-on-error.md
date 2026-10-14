# 2. Redis cache, stale on error

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Reference data from Commons (countries having photos per decade) changes rarely but is
needed on every draw. Fetching it each time would be slow and abusive.

## Decision

- Reference data is cached in **Redis**: countries per decade for 7 days.
- When a refresh fails, the **last cached value is served** and the refresh is retried
  later.
- Random draws are never cached.
- Cache access goes through one `shared` cache package; contexts never talk to Redis
  directly.

## Consequences

- Cached endpoints answer in milliseconds; the cache survives restarts.
- A Commons outage does not prevent drawing among known countries.
