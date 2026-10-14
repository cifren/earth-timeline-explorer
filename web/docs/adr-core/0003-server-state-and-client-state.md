# 3. Server state with TanStack Query, client state with Zustand

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The application holds two natures of state: data coming from the API (cached, refetched,
loading, failing) and choices made by the user (selected zone, decade, map view). Mixing
them in one store means re-implementing caching, or mirroring server data by hand.

## Decision

| Nature        | Owner          | Examples                                      |
|---------------|----------------|-----------------------------------------------|
| Server state  | TanStack Query | zones, geometry, images, availability by decade |
| Client state  | Zustand stores | selected zone, selected decade, map view      |

- TanStack Query is used in `application/` hooks; query functions call the context's
  infrastructure gateways, which return domain models.
- Zustand stores hold only client state; a store shared by several contexts lives in
  `core/`, a context's own store stays in the context.
- **No mirroring**: server data is never copied into a Zustand store, and client choices
  are passed to queries as keys (e.g. `['images', zone, decade]`).
- Components read state only through hooks.

## Consequences

- Caching, retries and loading states come from one library instead of hand-written code.
- Each piece of state has one owner; nothing needs to be kept in sync.
