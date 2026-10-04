# Verification — 2026-10-04

The user approved the four-board proposal with “looks ok”. Implementation follows
that direction without changing DTOs, routes, reading-list requests or review
requests. The page now uses a stable portrait cover/fallback, wrapping identity,
quieter rating/statistics, a bordered metadata surface, bounded summary text,
DTO-only sources and stable keyboard-operable tabs. The age badge uses neutral
page-scoped tokens to avoid low contrast on its previous bright fill.

## Automated checks

- `templ generate`: passed. Installed CLI v0.3.960 reports that go.mod uses
  v0.3.1020; generated files compile and remain ignored.
- `GOCACHE=/tmp/openlibrary-go-cache GOTMPDIR=/tmp go test ./web/public/...`:
  passed, including source-only rendering and public page/fragment regressions.
- `pnpm run build`: blocked by unrelated TypeScript errors in
  `web/frontend/src/block-editor/state.ts` (missing `./contracts`) and
  `web/frontend/src/islands/header-user-menu/UserMenu.tsx` (unused `useEffect`).
  Neither file was changed for this spec.
- `pnpm exec vite build --mode production`: passed, compiling the page SCSS and
  template utilities. Existing Sass, Browserslist and asset warnings remain.
- `git diff --check`: passed.

## Browser evidence

[The reproducible browser check](checks/browser.cjs) uses Go-rendered templates,
compiled production CSS/JavaScript, real Alpine and real Alpine AJAX fragment
replacement. Authentication fixtures run through the real middleware with stub
session/user services; no account or database is changed.

Thirty cases cover five fixtures, both themes, and 360/768/1280px:

| Fixture | Content checked |
| --- | --- |
| Empty | Zero counts, missing cover, absent summary/sources/tags, empty chapters/reviews |
| Long | Unbroken long title/author, twenty expanded long tags, long summary, DTO source |
| Populated | Existing image cover, long chapter title and reviewer name, review content |
| Authenticated | Author visibility notice, reading-list controls, existing own review |
| Adult warning | Warning bounds and dismissal followed by the authenticated page |

All cases have no horizontal overflow; cover bounds remain 200 × 300px. Long tag
names originally overflowed at 360px and now wrap within the book page. Own-review
and top-review grids use shrinkable text columns and appropriately sized avatars.
The current chapter highlight class now matches the existing style selector.

ArrowLeft/ArrowRight/Home/End select and focus tabs, update aria-selected, and
retain tab width, height and document position. Panels expose labels and focus
stops, loading remains visible, and panel semantics/show/hide behavior survive
chapter fragment replacement. Both populated and empty fragments are exercised.
Authenticated reading-list states are exercised without persisting mutations.
All five status labels fit on one line without truncation at each target width.
Following user feedback, the reading-list area and desktop action column are
260px wide while the centered cover remains 200px wide. Continue reading keeps
its existing click behavior; no navigation change was made. The authenticated
fixtures use the reported book/chapter IDs to preserve 64-bit string handling.

Sampled title, author link, metadata, age badge, statistics, inactive tab and
primary reading-list control contrast all pass 4.5:1. The minimum sampled ratio
is approximately 7.39:1. Page links/buttons/focusable panels receive a visible
foreground outline. Light mobile and dark desktop screenshots were visually
inspected. Screenshots/results are temporary artifacts in the fixture directory.

The browser harness supplies the shell's missing `#client-flashes` mount so an
existing common-script startup error does not interfere with page checks.
This does not change the global shell. It serves SVG assets with their correct
MIME type and holds chapter responses to test pending feedback.

## Reproduce

```sh
mkdir -p /tmp/book-page-fixtures
BOOK_PAGE_FIXTURE_DIR=/tmp/book-page-fixtures GOCACHE=/tmp/openlibrary-go-cache GOTMPDIR=/tmp go test ./web/public/... -count=1
pnpm exec vite build --mode production
node specs/public-book-page-polish/checks/browser.cjs
```

The script defaults to `/usr/bin/chromium-browser`; set
`PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH` for another installed Chromium. It uses the
existing local Playwright MCP dependency. Localhost/Chromium execution requires
permission outside the execution sandbox in this environment.

## Scope of evidence

Anonymous login/profile/chapter/source links and authenticated review/reading-list
controls remain in their original templates and retain their URLs and handlers.
Report still uses the existing report island. Existing adult-warning wording,
cookie behavior and author-only visibility notice are preserved. No illustrative
source is rendered, and absent optional sections produce no empty panels.

This is fixture verification, not live backend mutation testing. Review saving,
report submission, reading-list persistence and real user data were not changed
or exercised against a live service. Contrast sampling covers the changed page
surfaces and controls rather than auditing the global shell or every possible
rich-text color. The unrelated TypeScript failures prevent reporting a successful
full frontend build.
