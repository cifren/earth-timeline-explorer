# 1. OpenLayers for a flat, horizontally wrapped map

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The UX specification requires a very simple flat map in which latitude and longitude are
straight evenly spaced lines, repeating endlessly from left to right but bounded at the
poles, with clickable zones.

## Decision

The map uses **OpenLayers** with the **EPSG:4326** projection (equirectangular),
`wrapX` enabled for the horizontal repetition, the view constrained in latitude only, and
zones loaded as a GeoJSON vector layer with hover and click handling.

The map engine is used only inside the `geography` context's presentation layer; the rest
of the application knows zones, not OpenLayers objects.

## Consequences

- EPSG:4326 and world wrapping are native, not hand-built.
- Replacing the engine (e.g. for the 3D stage) only touches one context's presentation.
