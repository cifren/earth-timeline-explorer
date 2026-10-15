# 8. Scaffold with official generators

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Hand-written boilerplate (project setup, tool configuration, UI components) drifts from
upstream defaults and costs time. Most tools of the stack ship an official generator that
produces the expected structure and stays aligned with their releases.

## Decision

Whenever a tool provides an **official generator** (scaffolding command, "schematic"), it
is used to create files instead of writing them by hand:

| Need                       | Generator                                                 |
|----------------------------|-----------------------------------------------------------|
| React + TypeScript project | `npm create vite@latest -- --template react-ts`           |
| ESLint configuration       | `npm init @eslint/config@latest`                          |
| UI setup and components    | `npx shadcn@latest init`, then `npx shadcn@latest add <component>` |
| API client                 | `npx orval` from its configuration (ADR web core 0004)    |

- Generators run inside the web container (`./compose make web <target>` or
  `./compose bash web`), never on the host.
- The generated output is then adapted to the ADRs (layout, strict TypeScript, lint
  presets, path aliases) and committed.
- A generator that is re-run (e.g. the API client, a new UI component) has a Makefile
  target holding the exact command.
- Files with no official generator (Dockerfile, Prettier configuration) are written from
  the tool's documented reference.

## Consequences

- Project files start from upstream defaults and follow upstream upgrades.
- Adding a UI component is one command, not a copy from documentation.
- Generated files may need a first adaptation pass to match the ADRs.
