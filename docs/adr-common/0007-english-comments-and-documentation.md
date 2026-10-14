# 7. English comments and documentation

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The code base is read by different people and by AI assistants. Comments that narrate
history or state the obvious age badly and mislead.

## Decision

- Code, comments, ADRs and technical documentation are written in **English**.
- A comment exists only when the code's complexity requires it, and explains **how or
  why** the code works, based on facts in the current file.
- Comments never tell a change: no before/after, no migration status, no "replaced X".
- Documentation describes **what** the project does, **how** to run it and **why** it is
  built this way, so it can be deployed and improved without its authors.

## Consequences

- Comments stay true as long as the code they describe.
- History is found in version control, decisions in the ADRs.
