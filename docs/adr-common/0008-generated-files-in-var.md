# 8. Generated files live in `var/`

Date: 2026-10-15
Status: Accepted

## Context

Several files are produced by tools rather than written by hand: server and client code
from the API contract, zone outlines from Natural Earth. Mixed with hand-written sources,
they get edited by mistake, blur what to review, and hide which files a regeneration
overwrites.

## Decision

Every generated file lives under a **`var/` folder at the root of its application**,
one sub-folder per generated artefact (`var/<artefact>/`), never next to hand-written
sources. Examples of artefacts: code generated from the API contract, geographic
outlines generated from a public dataset.

- Each artefact has a Makefile target holding its generator command.
- `var/` holds **only** generated files: no hand-written file, not even a wrapper.
- Generated code is not part of the Shared pillar of ADR common 0004, even when it is
  business-free: it is separated by its nature, not by its content.
- Generated files start with a "generated, do not edit" header naming their generator.
- They are committed, so that a checkout builds without running generators; a
  regeneration is a reviewable diff.
- They are excluded from linting and formatting, but stay compiled and type-checked.
- Generated code depends on no hand-written package; sources import it.

## Consequences

- What may be edited is obvious from the path.
- A regeneration only ever touches `var/`.
- Imports of generated code read as such, through the `var/` path segment.
