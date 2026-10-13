# OpenLibrary

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

The dashboard areas serve administrators, book managers, and moderators working
with library accounts, books, chapters, reports, and activity records. Readers use
their personal library to find saved books by reading status and manage collections.

## Product Purpose

Support library management through existing searchable lists, detail screens,
editing workflows, and moderation actions. Support readers in recognizing saved
books, opening them, and organizing their collections.

## Capabilities and Constraints

The dashboard redesign preserves existing workflows, copy, inputs, routing, and
permissions. The sidebar remains on the left on desktop; small screens retain
the existing responsive navigation arrangement. No additional constraints were
requested. The approved reader redesign covers `/library`, `/library/archive`,
and `/library/collections`: existing statuses, previews, censorship, collection
creation, pagination, routes, and permissions remain intact. Other public pages
and the reading interface are outside this scope.

## Evidence on Hand

Existing implementation and authenticated local pages at `/admin`,
`/books-manager`, and `/moderation`. Existing logos, fonts, and light/dark themes
are available in `web/frontend`. The library direction is recorded in
`.impeccable/surfaces/library.md` and implemented in `library.templ` and
`common/style/components/library.scss`.
