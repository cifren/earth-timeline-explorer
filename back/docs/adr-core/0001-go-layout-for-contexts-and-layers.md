# 1. Go layout for bounded contexts and layers

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

ADR common 0004 defines bounded contexts, the `BC → Core → Shared` pillars and four
layers per context. Go needs a concrete package layout for them that the compiler and the
linter can protect.

## Decision

```
api/openapi.yaml               API contract, source of truth (ADR functional 0003)
cmd/
├── server/                    entry point of the HTTP service
└── <tool>/                    entry points of command-line tools, e.g. zones
internal/
├── app/                       composition root: wires contexts, HTTP server, middleware
├── bc/
│   └── <context>/             e.g. geography, imagery
│       ├── domain/
│       ├── application/
│       ├── infrastructure/
│       └── presentation/      what the outside sees: HTTP handlers, CLI commands
├── core/<concept>/            shared business, e.g. timespan, zoning
└── shared/<package>/          business-free code, e.g. apigen, problem, httpclient
```

- Everything lives under `internal/`, so no package is importable from outside the
  module.
- Package names are short, lowercase, single words; a package is named after what it
  provides, never `utils` or `common`.
- Each `cmd/` entry point only calls the composition root; a command-line tool is one
  more entry point on the same contexts, not a separate project.

## Consequences

- The tree mirrors the common architecture one to one.
- Go's `internal/` rule plus the boundary linter (ADR back core 0006) keep imports
  in the allowed directions.
