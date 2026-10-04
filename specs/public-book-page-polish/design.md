# Public book page proposal

Status: approved by the user on 2026-10-04 (“looks ok”); implemented; see verification.md.

Open [the four-board visual study](design-review.html) to compare 1280px desktop
and 360px mobile in both existing themes. Fixtures are explicitly labeled sample
content; no sample data or source will be added to the application.

## Proposed treatment

Retain the desktop cover/actions column, information column, and mobile order.
Use a 200 × 300 cover with a neutral themed fallback, subtle depth, and the
existing reading-list controls immediately below. Following implementation review,
the action area is 260px wide so Continue reading fits on one line; the cover
remains 200px wide and centered. Preserve all reading-list
states and destinations rather than replacing them with a generic Read action.

Make the serif title dominant (40px desktop, 30px mobile), allow full wrapping,
and keep the linked author directly below. Keep Report a visible 44px target.
Ratings, votes and views form a quieter wrapping row; retain existing stars.
Group chapter/word counts, publication date, age rating and tags in a restrained
surface with a thin border. Preserve the existing tag expansion behavior.

Keep About on the page background with approximately 65 characters per line.
Render sources only when provided by the DTO. Omit absent summary/source sections
rather than showing decorative empty panels. Remove the hard-coded Patreon link.
Use a stable serif size for both content tabs and indicate selection with an
underline and accessible state. Use buttons, roving focus, arrow keys, Home/End,
and named associated panels. Preserve chapter fragment loading and its slot ID;
the replaced fragment must retain panel semantics. Preserve all review actions.

## Inventory and verification to carry into implementation

Inspected book.templ, book.scss, book-toc.templ, book-reviews.templ,
reading-list.templ and theme.css. Preserve cover, author/profile URL, rating,
votes/views, chapter/word counts, created date, age rating, expandable tags,
rich summary, DTO sources, Report, reading-list controls, chapter loading,
review editing/rating/login, slug replacement, author visibility notice and
adult warning. No routes, DTOs or requests need changing.

Verify 360/768/1280px, long unbroken names, dense tags, rich/long summaries,
missing cover, absent optional sections, zero counts, empty reviews/chapters,
pending chapters, authentication and reading-list states. The proposal uses
existing palette values; AA contrast must be measured on final surfaces and
page-scoped adjustments made where existing muted tokens do not pass. Do not
change global theme tokens. Verify tab behavior after fragment replacement.

This artifact demonstrates the visual direction; it does not establish live
page correctness, final contrast, or completed acceptance criteria.
