# OpenLibrary

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

The dashboard areas serve administrators, book managers, and moderators working
with library accounts, books, chapters, reports, and activity records. Readers use
their personal library to find saved books by reading status and manage collections.
Public book detail pages help readers assess a book, choose a chapter, and start reading.
Public search helps readers discover books through queries, filters, and cover-led results.

## Product Purpose

Support library management through existing searchable lists, detail screens,
editing workflows, and moderation actions. Support readers in recognizing saved
books, opening them, and organizing their collections. Public book details preserve
book information and reading status while making chapter navigation easy to use.

## Capabilities and Constraints

The dashboard redesign preserves existing workflows, copy, inputs, routing, and
permissions. The sidebar remains on the left on desktop; small screens retain
the existing responsive navigation arrangement. No additional constraints were
requested. The approved reader redesign covers `/library`, `/library/archive`,
and `/library/collections`: existing statuses, previews, censorship, collection
creation, pagination, routes, and permissions remain intact. The separately approved
public book detail redesign covers `/book/:slug`, preserving metadata, summaries,
tags, ratings, chapters, reviews, external links, status controls, and permissions.
Start reading links directly to the first chapter only when one exists. The approved
public search redesign covers `/search`, preserving queries, sorting, ranges, tags,
pagination, censorship, and permission-dependent links. Other public pages and the
reading interface are outside this scope.

## Evidence on Hand

Existing implementation and authenticated local pages at `/admin`,
`/books-manager`, and `/moderation`. Existing logos, fonts, and light/dark themes
are available in `web/frontend`. The library direction is recorded in
`.impeccable/surfaces/library.md` and implemented in `library.templ` and
`common/style/components/library.scss`.

The book detail direction is recorded in `.impeccable/surfaces/book.md`, implemented
in `web/public/templates/book.templ`, `book-toc.templ`, and
`web/frontend/src/common/style/components/book.scss`. Code-first remains the
recorded preference; existing covers supply the visual material.

The public search direction is recorded in `.impeccable/surfaces/search.md`,
implemented in `web/public/templates/search.templ`,
`web/frontend/src/common/style/components/search.scss`, `book.scss`, and
`web/frontend/src/islands/search-filters/SearchFilters.tsx`.
