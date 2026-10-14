# 1. Use concise Architecture Decision Records

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

We need a lightweight, durable trace of *why* each architectural choice was made,
readable by humans and AI assistants without digging through history. The repository
holds several applications built with different tools.

## Decision

Every structural decision is recorded as a short ADR, following
[Michael Nygard's format](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions.html),
in a folder chosen by **where** it applies and **what** it is about:

| Folder                  | Content                                                        |
|-------------------------|----------------------------------------------------------------|
| `docs/adr-common/`      | Conventions shared by every application                        |
| `docs/adr-functional/`  | Whole-system decisions driven by the business (repository, contract, shared language) |
| `docs/adr-devops/`      | Whole-system build and run                                     |
| `<app>/docs/adr-core/`  | Reusable baseline of the application's stack                   |
| `<app>/docs/adr-functional/` | Business-driven decisions specific to this application  |
| `<app>/docs/adr-devops/`  | Build and run of this application                            |

`adr-core` versus `adr-functional`: a decision another project on the same stack would
reuse unchanged goes to `adr-core` (e.g. error format, folder layout); a decision driven
by Frise Chrono's business goes to `adr-functional` (e.g. Wikimedia Commons access, the
flat wrapped map).

`<app>` is `back`, `web`, later `mobile`. References between folders name the folder,
e.g. "ADR common 0003", "ADR web core 0002".

**Format** — every ADR is a Markdown file `NNNN-kebab-title.md`, in English:

```
# N. Title

Date: YYYY-MM-DD
Status: Accepted

## Context
## Decision
## Consequences
```

**Rules**
- Only **structural** decisions: technology choices, project structure, build and run
  tooling, cross-cutting concerns that are hard to reverse.
- **Concise**: a few lines per section; a table or a short example rather than prose.
- Decision and principle level, not implementation: no class names or file paths that
  drift. A short *Implementation status* note may point at where a decision is realised.
- Status: `Proposed`, `Accepted`, `Accepted — not yet implemented`, `Superseded by NNNN`.
- Once implementation starts, ADRs are immutable: a change of decision is a new ADR that
  supersedes the old one.

## Consequences

- Decisions are traceable and live next to the code they govern.
- Shared conventions are written once; each application only states how it applies them.
- A small writing overhead is added to each architectural change.
