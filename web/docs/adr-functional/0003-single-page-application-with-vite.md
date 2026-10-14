# 3. Single-page application with Vite and React Router

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The application is one interactive screen (map, timeline, panels) with no need for
server-side rendering in its first versions.

## Decision

- **Vite** builds a static single-page application.
- **React Router** handles routes; each page and heavy feature (e.g. the full-screen
  viewer) is lazy-loaded.
- Server-side rendering is out of scope until public shareable pages require it.

## Consequences

- Fast development server and small production bundles.
- The build output is plain static files, served from any static host.
