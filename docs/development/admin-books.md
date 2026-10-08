# Admin book navigation

`/admin#/books` lists all persisted books, including private, banned, shadow
banned, trashed, and permanently removed records. Name search is partial and
case-insensitive. Positive int64 IDs and `/book/<id>` or `/book/<slug>-<id>` URLs
resolve an existing book before opening `/admin#/books/<id>`. The table uses
server-side pagination with 20 records per page; `q` and `p` live in the hash
route. IDs remain decimal strings in the API and browser.

Admin details link to `/moderation#/books/<id>` and
`/book/<id>?admin.override=1`. The public handler checks the current admin role,
and BookService independently verifies the persisted actor role before applying
an override. Override permits restricted books and hidden chapters, carries
through slug normalization, TOC, reading links and chapter navigation, and shows
localized restriction reasons. When an admin opens a restricted book without the
override flag, the public page offers an explicit override action before reading
restricted content. Unrestricted books have no override warning.
Permanently removed content cannot be recovered; only the remaining persisted
record can be inspected.

The admin table also links to `/search?admin.link=1`. This flag adds admin detail
links to regular search results only for current admins; it does not change the
public search's visibility filters. Search submissions and pagination retain the
flag. Moderation uses DashboardShell and provides an admin dashboard link only
when the current user is an admin.
