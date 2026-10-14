# 8. Scaffold with official generators

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Hand-written boilerplate (module files, tool configuration, generated server code) drifts
from upstream defaults and costs time. Most tools of the stack ship an official generator
that produces the expected structure and stays aligned with their releases.

## Decision

Whenever a tool provides an **official generator** (scaffolding command, "schematic"), it
is used to create files instead of writing them by hand:

| Need                         | Generator                                              |
|------------------------------|--------------------------------------------------------|
| Go module                    | `go mod init`                                          |
| Pinned code-generation tools | `go get -tool` (recorded in `go.mod`, run with `go tool`) |
| Server code from the contract| `go tool oapi-codegen` (ADR back core 0003)            |
| Development reload config    | `air init`                                             |

- Generators run inside the back container (`./compose make back <target>` or
  `./compose bash back`), never on the host.
- The generated output is then adapted to the ADRs (layout, lint rules) and committed.
- A generator that is re-run (e.g. code generation) has a Makefile target holding the
  exact command.
- Files with no official generator (Dockerfile, golangci-lint configuration) are written
  from the tool's documented reference.

## Consequences

- Project files start from upstream defaults and follow upstream upgrades.
- Tool versions are pinned in `go.mod`; every developer runs the same generator.
- Generated files may need a first adaptation pass to match the ADRs.
