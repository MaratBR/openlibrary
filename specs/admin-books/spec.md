# Admin books

The requested behavior is the agreed feature scope.

## Acceptance criteria
1. Admin books lists all statuses in a paginated table and supports name search. A valid positive int64 ID or book URL is resolved first; an existing book immediately opens its admin details. Missing IDs show an empty result, and invalid URLs do not redirect.
2. Admin details display persisted book information and visibility/status flags, including deleted and banned records, and link to moderation and the public page with `admin.override=1`.
3. Only current admins can override public book and chapter restrictions. Restricted pages explain all applicable reasons at the top. Publicly accessible pages have no override warning. Override survives slug normalization, TOC loading and chapter navigation; hidden chapters can be inspected.
4. Moderation uses DashboardShell and offers an admin link only to admins.
5. Admin books links to `/search?admin.link=1`; search results show admin links only for admins with that flag and preserve it through searches/pagination.
6. Anonymous users, ordinary users and moderators cannot gain access by supplying either query flag. Missing book records remain not found.

7. Public-search admin links use an actionable chip style. Opening an existing restricted book without override offers current admins an explicit “Override and view book” action; restricted content remains unread until they choose it. Missing books and non-admin visitors retain the normal unavailable page.

8. Search admin chips align with the text column, stay on one line, and have space above and below. Pagination uses previous/next controls, first/last pages, ellipses, and an accessible current-page indicator; it preserves search filters and fits mobile screens.
