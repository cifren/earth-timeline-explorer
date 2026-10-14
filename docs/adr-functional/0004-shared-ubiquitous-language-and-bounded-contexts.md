# 4. Shared ubiquitous language and bounded contexts

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Back and web are partitioned by bounded context (ADR common 0004). If each picks its
own names, the same concept ends up called differently on each side and in the
specifications.

## Decision

The **ubiquitous language** of the functional specification is used everywhere: code,
API, documentation. Back and web use the **same bounded context names**:

| Bounded context | Responsibility                                               | Starts |
|-----------------|--------------------------------------------------------------|--------|
| `geography`     | Zones and sub-zones, their geometry                          | v0     |
| `chronology`    | Periods: decades, eras, timeline scale                       | v0     |
| `imagery`       | Images, kinds, provenance, certification, random draw, sources | v0   |
| `place`         | Identified places, viewpoints, key events                    | N3     |
| `moderation`    | Reference registry, reports, error notes                     | Control |
| `identity`      | Accounts and roles                                           | Control |

Core vocabulary: *Image*, *Zone*, *TimeSpan*, *Decade*, *Kind*, *Provenance*,
*Certification*, *Authority*, *Place*, *Viewpoint*. A context exists in an application
only when that application has something to do in it.

## Consequences

- A concept has one name from the specification down to the code.
- A context can be located at once in both applications.
- Adding a context is a deliberate, documented act.
