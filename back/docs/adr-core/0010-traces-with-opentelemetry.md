# 10. Traces with OpenTelemetry

Date: 2026-10-15
Status: Accepted

## Context

A request may read Redis and call Wikimedia Commons several times; logs alone do not
show where its time goes, nor which log belongs to which request.

## Decision

- Requests are traced with **OpenTelemetry**, the vendor-neutral standard, and exported
  over **OTLP/HTTP**. Context crosses services with the **W3C Trace Context** headers.
- Instrumentation comes from the official libraries rather than hand-written spans: the
  HTTP server and the Commons client through `otelhttp`, Redis through `redisotel`.
- The exporter is configured by the standard `OTEL_*` environment variables; tracing is
  off when `OTEL_EXPORTER_OTLP_ENDPOINT` is unset, so tests and tools run without it.
- The `slog` handler ties logs to traces: a log written during a traced request carries
  `trace_id` and `span_id`, and is copied as an event onto the current span; an error
  log marks the span as failed.
- In development, **Jaeger** receives and displays the traces, keeping them in memory.

## Consequences

- One view per request, from the HTTP call down to each Redis command and Commons call,
  with the logs written along the way.
- Any OTLP backend (Grafana Tempo, a hosted service) can replace Jaeger without code
  changes.
- Traces are lost when Jaeger restarts; that is accepted for development.
