# Implementation plan

## Affected areas

- `web/public/templates/book.templ`: presentation markup and tab semantics.
- `web/frontend/src/common/style/components/book.scss`: BookPage styles.
- Inspect `book-toc.templ`, `book-reviews.templ`, and reading-list components for
  integration; change only presentation needed for a consistent book page.
- `translations/en.toml` only if accessible labels or presentation copy change.

## Design and data changes

Read public/frontend scoped guides. Build a reviewable visual proposal from the
current page, covering real content and edge cases; record accepted decisions in
this spec directory before coding. Refine existing columns and section grouping
with typography, spacing, restrained surfaces, and stable tab styling. Prefer
existing tokens and primitives; scope page presentation to BookPage rather than
changing unrelated shared components.

Keep DTOs, routes, fragment identifiers, Alpine behavior, and action contracts.
Preserve rich text rendering and current chapter/review data. Audit sources:
`bookExternalLinks` currently appends a hard-coded Patreon link. Remove this
illustrative addition while preserving DTO-provided sources to meet AC-3; no
backend change is needed. Coordinate form-control appearance with
[unified input styles](../unified-input-styles/spec.md) instead of page overrides.
No migrations or deployment coordination are required. Regenerate templ output
locally after template changes; follow repository generated-source rules.

## Verification

| Criterion | Verification |
| --- | --- |
| AC-1 | Record accepted desktop/mobile proposal and both theme treatments |
| AC-2 | Compare rendered page with proposal; assess heading/action hierarchy |
| AC-3 | Data/action inventory before/after; authenticated/anonymous states; sources match DTO |
| AC-4 | Long/missing/empty content fixtures and pending chapter fragment |
| AC-5 | 360/768/1280px renders and overflow/line-length checks |
| AC-6 | Contrast, keyboard focus, tab selection and panel association in both themes |

During implementation use `templ generate`, `go test ./web/public/...`,
`pnpm run build`, and `git diff --check`. Record unrelated baseline failures.
