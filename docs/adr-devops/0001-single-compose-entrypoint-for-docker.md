# 1. Single `./compose` entrypoint for all Docker usage

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Every Docker interaction needs the same boilerplate: host UID/GID so files created in
containers stay owned by the host user, the right service, first-run setup. With several
applications in one repository, typing `docker compose` by hand multiplies the ways to
get it wrong.

## Decision

**All Docker-related commands go through one executable script, `./compose`, at the
repository root.** It wraps `docker compose` and owns:

- **rights management**: exports the host UID/GID;
- **network**: every service joins the shared external `jacquet` network, created when
  missing, because Docker's default address range collides with the VPN;
- **project initialisation**: build images, install dependencies of every application;
- **lifecycle**: `up`, `down`, `logs`, `build`;
- **shells**: `./compose bash <service>` and `./compose bash-root <service>`;
- **app commands**: `./compose make <service> <target>` runs a target of that
  application's Makefile inside its container;
- **`./compose help`**: lists the sub-commands.

Each application's Makefile is only a **catalogue of its commands** (`lint`, `build`,
`api`…), holding the raw commands run inside its container. The host needs only Docker:
no Go or Node toolchain.

Services: `back`, `web`, `redis`.

## Consequences

- One source of truth for how Docker is invoked; rights and flags cannot be forgotten.
- Onboarding is `./compose init` then `./compose help`.
- Adding the mobile application means adding a service and a Makefile, not a new way of
  running things.
