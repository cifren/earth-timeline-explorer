# 1. Go layout for bounded contexts and layers

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

ADR common 0004 defines bounded contexts, the `BC → Core → Shared` pillars and four
layers per context. Go needs a concrete package layout for them that the compiler and the
linter can protect.

## Decision

```
api/                       API contract (ADR functional 0003)
cmd/<command>/             entry points: the HTTP service, command-line tools
internal/
├── app/                   composition roots
├── bc/<context>/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   └── presentation/      what the outside sees: HTTP handlers, CLI commands
├── core/<concept>/        business shared by several contexts
└── shared/<package>/      business-free code
var/<artefact>/            generated code only (ADR common 0008)
```

- Everything lives under `internal/`, so no package is importable from outside the
  module.
- Package names are short, lowercase, single words; a package is named after what it
  provides, never `utils` or `common`.
- Each `cmd/` entry point only calls the composition root; a command-line tool is one
  more entry point on the same contexts, not a separate project.

## Consequences

- The tree mirrors the common architecture one to one.
- Go's `internal/` rule keeps every package private to the module; the direction of
  imports between pillars and layers is left to a boundary checker (ADR common 0005).
