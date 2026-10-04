# Feature: Public book page visual polish

## Problem and outcome

The public book page has useful information and a suitable overall structure,
but its presentation needs stronger hierarchy and a more appealing finish.
Readers should quickly recognize the book, author, reading actions, and content
sections within the site's existing visual identity.

Status: visual proposal in `design-review.html` approved by the user on
2026-10-04 (“looks ok”); implemented; verification recorded in `verification.md`. Book page means the public
`/book/{slug}-{id}` page, not the manager overview.

## Scope

Polish typography, spacing, surfaces, borders, cover presentation, metadata,
ratings, tags, reading actions, summary, sources, and Chapters/Reviews navigation.
Retain the desktop cover/actions column and main information column, plus the
existing mobile stacking order and section order. Preserve all current data,
links, reading-list actions, reviews, chapter loading, and visibility warnings.
No new data, features, endpoints, or redesign of the reader, manager, search,
or global site shell is included.

## Scenarios and acceptance criteria

- AC-1: A visual proposal demonstrates desktop and mobile versions in both
  themes before implementation. It retains the current structure while using
  existing serif headings, cream/dark surfaces, green actions, and theme tokens.
  Record the chosen proposal so the visual outcome can be reviewed consistently.
- AC-2: The book title is the primary heading, author is clearly linked, and
  reading actions are easy to find beside/below the cover. Ratings and secondary
  metadata have a quieter hierarchy. Spacing, type sizes, and surface treatments
  consistently distinguish book identity, summary, and chapter/review content.
- AC-3: Every currently supported information item and action remains available:
  cover, author, rating/votes/views, chapter/word counts, publication date, age
  rating, tags, summary, sources, report, reading-list controls, chapters, and
  reviews. Existing URLs, visibility notices, adult warning, and request behavior
  retain their meaning. No illustrative ratings, links, or other data are added.
- AC-4: Long titles/authors, many tags, long summaries, missing covers, absent
  summaries/sources, zero counts, and empty chapters/reviews remain readable.
  Optional absent sections introduce no empty decorative panels; cover fallbacks
  preserve proportions. Existing empty and loading feedback remains visible.
- AC-5: At 360px, 768px, and 1280px the page has no horizontal overflow; titles,
  metadata, actions, and tabs wrap without clipping. Summary text has a readable
  line length, and cover proportions remain stable.
- AC-6: Both themes use legible text and controls with WCAG AA contrast, visible
  keyboard focus, and usable action targets. Chapters/Reviews tabs expose their
  selected state and support keyboard operation with associated panels. Tab
  selection does not change text size enough to move surrounding layout.

## Design decision

The user approved `design-review.html` on 2026-10-04 (“looks ok”). The
implementation follows that direction; see `design.md` and `verification.md`.
