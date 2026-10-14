# 5. Images hotlinked, never stored

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

Image files are large and owned by their sources (Wikimedia Commons for v0). Copying or
proxying them would cost storage and bandwidth and blur who publishes them.

## Decision

The system never stores or proxies image files. The API returns the source URLs (a
thumbnail and the original), and clients load them directly from the source's CDN, always
displayed with their attribution (license, author, link to the source page).

The system only keeps references to images (source and identifier), as specified for the
control layer.

## Consequences

- No storage or bandwidth cost for images; the back only handles metadata.
- Availability of images depends on the source; a missing file is handled by the draw
  fallback.
- Attribution required by Creative Commons licenses is always shown.
