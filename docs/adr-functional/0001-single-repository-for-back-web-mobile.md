# 1. Single repository for back, web and mobile

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The system has one back and several clients (web now, mobile later) that evolve together
around one API. Separate repositories would make every contract change a multi-repository
release.

## Decision

One repository holds every application, each in its own top-level folder:

- `back/` — the Go service, which also owns the API contract and
  the generation of its reference data (e.g. zone geometry);
- `web/` — the React web application;
- `mobile/` — the mobile application, when it starts;
- `docs/` — whole-system documentation and ADRs.

Each application is self-contained: its own dependencies, container image, Makefile and
`docs/`. Applications never import each other's code; clients only read the back's API
contract to generate their client code.

## Consequences

- A contract change and its consumers land in the same commit.
- One clone gives the whole system.
- Each application keeps its own tooling and can be built and deployed on its own.
