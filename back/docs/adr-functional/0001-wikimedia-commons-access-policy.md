# 1. Wikimedia Commons access policy

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Commons is a shared public service with an API etiquette and a User-Agent policy. Abusing
it gets clients blocked, and its latency drives the latency of the random draw.

## Decision

The Commons gateway (infrastructure of the `imagery` context):

- sends an identifying **User-Agent** with a contact address on every request;
- adds the **`maxlag`** parameter and backs off when Commons reports lag;
- sends requests **one at a time** per draw, never in parallel bursts;
- applies a **3-second deadline** per call (ADR back core 0005);
- retries a draw at most **3 times** when a file is missing, then returns no image.

## Consequences

- The back stays a well-behaved Commons client.
- A slow Commons bounds the draw time instead of blocking it.
