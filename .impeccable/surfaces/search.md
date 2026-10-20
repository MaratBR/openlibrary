---
version: 1
slug: "search"
primary_target: "web/public/templates/search.templ"
related_targets: ["web/frontend/src/common/style/components/search.scss", "web/frontend/src/common/style/components/book.scss", "web/frontend/src/islands/search-filters/SearchFilters.tsx", "/search"]
---

# Public book search

Mode: Operate. The user requested a redesign of `/search`, moving filters to the
side and then giving the filter container a surface background. The confirmed
direction extends the existing warm paper and forest green visual world.

## Direction contract

THESIS: A full-width query introduces a clear choice: refine in a contained left
panel or scan cover-led results. Existing covers, serif headings, and flat ruled
book rows connect discovery to the library and book detail surfaces.

FIRST VIEWPORT: The page caps at 1120px. A query field spans the layout above a
296px filter rail and flexible results column, separated by 40px. The filter
container uses the shared surface background, a one-pixel border, 12px corners,
and 24px padding. At 1000px and below, the rail becomes 264px with a 24px gap
and 20px padding. At 767px and below, filters collapse behind a labeled button
above results; the panel retains its surface and uses 12px padding.

SIGNATURE: Covers lead quiet, divided result rows. Desktop covers are 112px
wide with a 24px text gap; intermediate widths use 88px covers and a 20px gap,
and widths up to 480px use 72px covers and a 16px gap. Covers retain a subtle
spine shape and shadow while the rows remain flat.

FORM: The title uses the existing Literata stack, with sans-serif labels and
supporting text. Surface, foreground, primary, accent, border, and ring colors
come from shared semantic tokens in both themes. Existing control components
remain the basis of the form; no new fonts, palette, or artwork are introduced.

## Scope and behavior

Preserve the query, sorting, numeric and date ranges, included/excluded tags,
age rating, page size, pagination, censorship, highlighted matches, and
permission-dependent administration links. Keep real titles, authors, tags,
summaries, and covers in the results.

Place advanced metrics and date fields in a native disclosure, initially open
when those fields have active values. Reset clears the filter draft to defaults;
the next submission applies it. Submitting the top query also submits the current
filter draft. Filter submission starts at page one. The mobile toggle exposes its
expanded state and names its controlled panel; keyboard focus stays visible.

The top query form submits on its current route, including when reused on a
profile, rather than forcing navigation to `/search`. The search-specific result
modifier is applied when search match data is passed to the book card. In those
rows, the preview disclosure button remains in normal flow below the text, with
an 8px top margin and a 36px minimum height; the preview's bottom margin is reset
so the button cannot overlap the excerpt. Other book-card previews retain their
existing disclosure positioning.

## Evidence and design drift

Implementation evidence is in the source targets listed above. The existing
`DESIGN.md` and `.impeccable/design.json` remain the visual authority for shared
tokens and components; this brief records the scoped extension.

`DESIGN.md` still describes search results as retaining incumbent styling and
does not list public search among its approved reader surfaces. That scope text
is stale relative to this approved redesign. Report the discrepancy without
rewriting the global design document or sidecar during this scoped change.

This brief records source behavior and confirmed direction; it does not claim
browser or automated verification. The implementation task owns those checks.
