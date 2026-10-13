---
version: 1
slug: "library"
primary_target: "web/public/templates/library.templ"
related_targets: ["web/frontend/src/common/style/components/library.scss", "/library", "/library/archive", "/library/collections"]
---

# Reading gallery

Mode: Operate. The user approved a cover-led reading gallery across My library,
Archive, and Collections. Code-first follows the recorded project preference.

## Direction contract

THESIS: Readers recognize a saved book and open it. Real covers establish identity;
serif headings and flat ruled sections give reading statuses a calm hierarchy.

FIRST VIEWPORT: A Literata heading above compact underlined tabs on a warm surface.
A reading-status heading and fine rule introduce regular covers with fully readable
titles. Mobile retains two covers per row. Both themes use shared semantic tokens.

SIGNATURE: Cover shadows deepen and titles underline on hover. Keyboard focus
outlines the complete book. Reduced motion disables cover and tab transitions.

FORM: Two-to-three cover proportions, small cover corners, generous separation
between reading statuses, compact labels, flat collection groups, and wrapping
collection cover strips. No generated artwork or fictional metadata.

## Scope and behavior

Preserve reading statuses, archive groups, book previews, censorship, collection
creation, pagination, routes, and permissions. Collection creation reveals the
existing form with a visible input label and focused input. Profile collection
cards retain their compact styling; gallery variants apply only in the library.
Dashboard design and the reading interface are outside this surface's scope.

## Evidence and verification

Source: `web/public/templates/library.templ` and
`web/frontend/src/common/style/components/library.scss`. Token and component
reference: `DESIGN.md` and `.impeccable/design.json`. Approved direction:
`.impeccable/library/direction.md`.

Finish review reported a ship verdict across desktop, mobile, and dark theme,
including synthetic populated collection fixtures. Synthetic fixtures validate
layout, not production collection data. Real covers remain the primary visual
material; long titles wrap completely and grid rows align by their tallest item.
