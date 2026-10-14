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
- country → zone table: UN M49 regions and sub-regions, grouped into the 9 zones;
- countries dissolved into zones and simplified with GEOS (Go bindings), output as
  GeoJSON (RFC 7946); GEOS is used only by this command, never by the HTTP service;
- the same table maps Commons country names to ISO 3166-1 codes.

The back embeds the files in its binary and serves them with long-lived HTTP caching.

## Consequences

- Geometry costs nothing at runtime.
- Changing a zone definition is a regeneration plus a reviewable diff.
