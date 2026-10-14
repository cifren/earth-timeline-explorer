# 5. Enforce architecture boundaries with a linter

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The rules of ADR common 0004 (contexts isolated, inward dependencies, `BC → Core →
Shared`) only hold by discipline if nothing checks them.

## Decision

Each application makes these rules **executable** with a dependency rule checker run as
part of its `lint` command: a forbidden import is a lint error. The tool is chosen in each
application's `adr-core`.

## Consequences

- Architecture violations are caught mechanically, not by reviewers.
- One configuration file per application becomes the source of truth for boundaries.
