# Implementation plan

Reuse the authorized moderation book search for the admin table, with all statuses included. Add admin read endpoints and React routes; serialize IDs as strings. Use BookService details for the admin page with a role-checked override, including status flags and reason codes.

Apply visibility checks in BookService and chapter reads. Verify override actor roles against persistence. Carry the authorized override through public templates and navigation, and render localized reason codes. Reuse existing SQLC chapter queries with an optional override for navigation and all-chapter TOC reads.

Replace moderation's outer frame with DashboardShell and pass current admin eligibility from the server. Add conditional search-result admin links. No migration is needed.

Verify parser and role/visibility regressions with focused Go tests, regenerate templ/SQLC, run affected Go tests, frontend build, and whitespace checks. Document navigation and access conventions.

When a normal book read returns a visibility denial, render an admin-only opt-in page with a same-book URL preserving query parameters and setting `admin.override=1`. Do not retry the restricted read until the user follows that link. Use existing Chip styles for search admin actions.

Use a scoped book-search action style for chip sizing and spacing. Shared server-rendered pagination shows a compact range with first/last pages, ellipses, and previous/next controls; retain AJAX attributes and query parameters. Verify real screenshots and page navigation in Playwright.
