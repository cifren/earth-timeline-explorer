# 3. Zone geometry generated offline

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The map needs the outline of the 9 zones. Computing it on each request, or depending on a
live geometry service, is slow and fragile.

## Decision

Zone geometry is **generated once** by a command of the `geography` context, run through
its own entry point (`cmd/zones`), and committed:

- source: Natural Earth Admin 0 countries, 1:50m (public domain);
- country → zone: the UN M49 region and sub-region recorded by Natural Earth for each
  country, grouped into the 9 zones; overseas territories listed apart by M49 (French
  Guiana, Réunion…) join the zone they lie in, not their country's;
- countries dissolved into zones and simplified with GEOS (Go bindings), output as
  GeoJSON (RFC 7946); GEOS is used only by this command, never by the HTTP service;
- Commons country names are placed in zones by a separate table of the back, since
  Commons names countries differently and includes historical entities.

The command writes the outlines as a generated Go package in `var/` (ADR common 0008),
so they are compiled into the binary; the back serves them with long-lived HTTP caching.
This generated package groups every geographic outline, so further outlines (e.g.
sub-zones) join it next to the zones.

## Consequences

- Geometry costs nothing at runtime.
- Changing a zone definition is a regeneration plus a reviewable diff.
