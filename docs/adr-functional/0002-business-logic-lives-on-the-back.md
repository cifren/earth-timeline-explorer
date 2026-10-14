# 2. Business logic lives on the back

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Rules such as the random draw, date normalization (EDTF), zone assignment or trust levels
must not be re-implemented, and drift, in every client. Web and mobile are consumers of
the same rules.

## Decision

All business rules live in the back. Clients hold **no business rule**: they display what
the API returns and orchestrate the user interaction (selection, navigation, loading and
error states). The back stays authoritative and validates every input.

Example: the back returns a date both as EDTF (`195X`) and as a plain English label
("1950s"); clients display the label, they never parse EDTF.

## Consequences

- One source of truth for the rules; no drift across web and mobile.
- Clients stay thin and can be rewritten or added without touching the rules.
- Client-side logic is limited to presentation and interaction.
