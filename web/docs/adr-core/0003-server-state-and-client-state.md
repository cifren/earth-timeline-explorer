# 3. Server state with TanStack Query, shared choices in the URL

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
| Shared client state | The page URL (search parameters, through React Router) | selected zone, decade and photo |
| Local client state  | React state in the component | open viewer, zoom, drag in progress |

- TanStack Query is used in `application/` hooks; query functions call the context's
  infrastructure gateways, which return domain models.
- Choices shared by several contexts live in the URL, read and written by one hook in
  `core/` (`useSelection`): every screen has an address, to share or to come back to
  through the browser history, without a client store.
- **No mirroring**: server data is never copied into the URL or a component state, and
  client choices are passed to queries as keys (e.g. `['images', zone, decade]`).
- Components read state only through hooks.

## Consequences

- Caching, retries and loading states come from one library instead of hand-written code.
- Each piece of state has one owner; nothing needs to be kept in sync.
- Shared choices must fit in a URL: identifiers and numbers, never data.
