---
version: 1
slug: "book"
primary_target: "web/public/templates/book.templ"
related_targets: ["web/public/templates/book-toc.templ", "web/frontend/src/common/style/components/book.scss", "/book/:slug"]
---

# Literary dossier

Mode: Operate. The user requested an interesting, functional redesign of
`/book/cyberworld-alice-terror-rising-6707299548257058898` and approved
“Go ahead—surprise me.” Code-first follows the recorded project preference.

## Direction contract

THESIS: A cover-led literary dossier helps readers assess the book and begin
reading. A large Literata masthead, compact author/rating byline, and a cover
shelf establish identity; flat metadata and sequential chapter records support use.

FIRST VIEWPORT: A ruled masthead above a sticky desktop cover/action rail and
book details. The title reaches 60px on desktop, with a phone clamp beginning
at 32px. The page caps at 1120px with a 320px rail, 48px column gap, and 32px
row gap. Below 768px, a compact cover sits beside the reading actions.

SIGNATURE: The existing cover has asymmetric spine corners and a soft offset
shadow. These are intentional scoped additions, alongside the masthead and
section type sizes. No new raster assets or theme colors are introduced.

FORM: Two-column metadata below the reading controls, 17px summary text with 1.8 leading (15px on phone),
and single-column chapter records at least 88px tall. Real content and semantic
theme colors remain the material in light and dark modes.

## Scope and behavior

Limit implementation to the book detail templates and BookPage/BookToc styles.
Preserve content, ratings, tags, reviews, external links, existing save/status
controls, and permissions. Start reading opens the first chapter directly only
when one exists. Keep chapter search and reading state. Chapters/Reviews tabs
support keyboard focus and selection, with visible focus rings. Global theme,
navigation, search results, and the reading interface retain incumbent behavior.

## Evidence and verification

Sources and tokens: the related targets above, `DESIGN.md`, and
`.impeccable/design.json`. Detailed verification is recorded in
`specs/book-page-dossier/verification.md`. Screenshots:
`.impeccable/review/book/desktop.png`, `mobile.png`, and `dark.png`.

Finish review: ship after keyboard verification, now confirmed. ArrowRight
focuses/selects Reviews; Home focuses/selects Chapters. Desktop, phone, and dark
views show no page overflow; actual first-chapter navigation was verified. Public
Go tests and Vite assets pass. Full pnpm build is blocked by the baseline unused
`useEffect` import in `UserMenu.tsx`.

## Density refinement

The user liked the design but found it empty. The cover now measures 240 × 360px
on desktop, the synopsis aligns with its top, and the narrower content grid uses
tighter spacing. Book facts support the cover/actions instead of displacing the
synopsis. Mobile retains its compact cover and two-column facts.
