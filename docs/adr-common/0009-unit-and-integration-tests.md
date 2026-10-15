# 9. Unit and integration tests

Date: 2026-10-15
Status: Accepted

## Context

The business rules (date normalization, zone assignment, random draw, trust levels) and
the mappers between layers are where regressions hide; the wiring between layers and
with real dependencies is where the rest hide. Static analysis checks neither, and manual
checks through the running application are slow and partial.

## Decision

Every application has two levels of automated tests, each run by its own command in the
application's container:

| Level           | Scope                                                     | Dependencies                                          | Command            |
|-----------------|-----------------------------------------------------------|-------------------------------------------------------|--------------------|
| **Unit**        | One rule, mapper or translation, in isolation             | None: no network, no database, no cache server        | `test`             |
| **Integration** | One use case through its real layers, from the entry point (HTTP route, page) to the gateways | Internal dependencies run for real (e.g. a cache server); external public services are faked | `test-integration` |

Rules for both levels:
- **External public services are never called** (image sources, datasets): they are
  replaced by local fakes serving recorded or hand-written answers, in-process.
- Tests sit **next to the code** they test, following the ecosystem's convention;
  integration tests are told apart by their name or a build marker.
- **Table-driven** when a rule has several cases: one table of inputs and expected
  outputs, one loop.
- A bug fix comes with the test that reproduces it.
- Generated code (`var/`) is not tested directly; it is exercised through the code using
  it.

What each level covers:
- **Unit**: domain and Core rules, use-case decisions, mappers (anti-corruption layers),
  presentation translations (validation, error formats).
- **Integration**: a request or a user action end to end inside the application — e.g.
  "drawing images of a zone returns images from the source and caches the source's
  country list".

End-to-end tests across applications (browser driving the web against the back) are a
separate decision, taken when needed.

## Consequences

- Rules change with confidence; the wiring is checked without starting the whole system
  by hand.
- Tests document expected behaviour with concrete examples.
- Unit tests stay fast enough to run on every change; integration tests are run before
  sharing changes.
