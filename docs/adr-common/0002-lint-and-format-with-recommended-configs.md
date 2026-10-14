# 2. Lint and format with recommended configs

Date: 2026-10-14
Status: Accepted — not yet implemented

## Context

The compiler catches type errors, not style or correctness smells. A hand-curated rule
set is a maintenance burden and a source of endless debate; manual formatting wastes
review time.

## Decision

Every application uses its ecosystem's **standard linter** with its **recommended
presets**, and its ecosystem's **standard formatter**:

- no bespoke rule set; a deviation from the presets needs an ADR;
- the formatter formats, the linter lints: separate lanes, no overlap;
- generated code is excluded from linting and formatting, but stays type-checked;
- one `lint` command per application, run in its container.

The concrete tools are chosen in each application's `adr-core`.

## Consequences

- A shared, low-noise baseline that follows upstream on upgrade.
- Formatting leaves code review.
- Each application adds its linter and formatter as development dependencies.
