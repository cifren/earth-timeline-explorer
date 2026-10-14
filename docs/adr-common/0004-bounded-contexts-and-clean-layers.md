# 4. Bounded contexts and clean layers

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Splitting code first by technical layer (`models/`, `services/`, `handlers/`) scatters one
business area across the tree and invites coupling between areas. Within an area, UI or
transport churn must not ripple into the business model.

## Decision

**Simplified DDD**, the same in every application:

1. **Bounded contexts first.** The top-level partition of the code is the bounded context
   (names in ADR functional 0004). Contexts never import each other; they coordinate through
   Core or through the back.
2. **Three pillars**, one-way dependencies `BC → Core → Shared`:
   - **BC**: business contexts;
   - **Core**: business shared by several contexts (e.g. *TimeSpan*, *Zone* identifiers);
     never depends on a BC;
   - **Shared**: generic, business-free code (helpers, UI primitives, generated API
     code); depends on neither.
3. **Four layers per context**, dependencies pointing inward:
   `presentation → application → infrastructure → domain`.
   - `domain`: entities, value objects, typed errors; no framework, no I/O;
   - `application`: use cases orchestrating the domain and the infrastructure;
   - `infrastructure`: gateways to the outside world (HTTP, cache, sources) and mappers
     (anti-corruption layer); returns domain models;
   - `presentation`: what the outside sees (HTTP handlers for the back, components for
     the web).
4. **No ports by default.** Upper layers depend on concrete infrastructure. An interface
   is introduced only when a second implementation or a test double actually exists.
5. A **composition root** per application wires the contexts together; it is the only
   place that knows them all.

## Consequences

- A business area lives in one folder, in both applications, under the same name.
- The domain has no dependency on frameworks, transport or storage.
- Less boilerplate than full hexagonal architecture; an interface appears when needed.
- Rules hold by convention until enforced by a boundary linter (ADR common 0005).
