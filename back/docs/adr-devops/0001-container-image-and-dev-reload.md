# 1. Container image and development reload

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The back must start quickly in development and ship as a small, safe image.

## Decision

- **Production image**: multi-stage build; the Go binary is compiled statically and
  copied into a **distroless** static base image, run as a non-root user.
- **Development**: the source is mounted in a Go container and **Air** rebuilds and
  restarts the service on each change.
- The back's **Makefile** catalogues its commands (`run`, `build`, `lint`, `api`), run in
  its container through `./compose make back <target>` (ADR devops 0001).

## Consequences

- Production image of a few tens of megabytes, no shell, minimal attack surface.
- Edit-and-reload loop in development without a local Go installation.
