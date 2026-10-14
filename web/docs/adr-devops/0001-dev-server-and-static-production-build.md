# 1. Development server and static production build

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The web must reload instantly in development and ship as plain static files.

## Decision

- **Development**: the Vite development server runs in a Node container with the source
  mounted, with hot module replacement.
- **Production**: multi-stage build; the static bundle is served by **Caddy**, with
  long-lived caching on hashed assets and a fallback to `index.html` for client routes.
- The web's **Makefile** catalogues its commands (`dev`, `build`, `lint`, `format`,
  `api`), run in its container through `./compose make web <target>` (ADR devops 0001).

## Consequences

- No Node installation needed on the host.
- Production serves static files only: no Node process in production.
