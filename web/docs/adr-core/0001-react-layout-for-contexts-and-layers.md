# 1. React layout for bounded contexts and layers

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

ADR common 0004 defines bounded contexts, the `BC → Core → Shared` pillars and four
layers per context. A React application needs a concrete folder layout for them.

## Decision

```
src/
├── app/                        composition root: routes, providers, page shells
├── bounded-contexts/
│   └── <context>/              e.g. geography, chronology, imagery
│       ├── domain/             types, value objects, pure functions
│       ├── application/        hooks orchestrating server and client state
│       ├── infrastructure/     gateways calling the API client, mappers
│       ├── presentation/       components
│       └── index.ts            public API of the context
├── core/<concept>/             shared business, e.g. exploration (selected zone, decade)
└── shared/
    ├── api-client/             generated from the OpenAPI contract
    └── <package>/              business-free UI primitives, hooks, helpers
```

- Outside code imports a context only through its `index.ts`, and only from `app/`.
- A page in `app/` composes components of several contexts; contexts exchange state
  through `core/`, never by importing each other.

## Consequences

- The tree mirrors the common architecture and the back's context names.
- Each context is lazy-loadable and replaceable on its own.
