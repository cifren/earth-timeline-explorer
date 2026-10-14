# 7. Environment configuration and structured logs

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The same image must run in development and production, and logs must be searchable.

## Decision

- Configuration comes only from **environment variables**, read once at startup into a
  typed configuration struct; a missing required value stops the service with a clear
  message.
- Logs are written to standard output as **JSON** with `log/slog`, one event per line,
  with a request identifier on every request-scoped log.

## Consequences

- One image for every environment; no configuration file to ship.
- Logs are readable by any log collector without parsing rules.
