# 2. UI components with Tailwind CSS and shadcn/ui

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The interface needs a side panel, buttons, skeleton placeholders, a full-screen viewer
and later a slideshow, accessible and consistent, without fighting a heavy component
library's styles.

## Decision

- **Tailwind CSS** for styling.
- **shadcn/ui** components (built on Radix primitives) copied into `shared/`,
  owned and adapted in the code base: sheet (side panel), button, skeleton, dialog.

## Consequences

- Accessible primitives (focus, keyboard, ARIA) without writing them.
- Components are ordinary project code: no upstream styling to override.
