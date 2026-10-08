# Censored tags: agreed behavior

Requested 2026-10-08. This slice implements the existing settings modes:

- Readers select any existing tag, search by name, remove selections, and reload
  their saved list. IDs travel as decimal strings; legacy names remain resolvable.
- `none` does nothing. `hide` excludes books matching any selected tag before
  pagination/counts/random selection. Direct book/chapter/preview content is
  blurred in both hide and censor modes. Explicitly including a banned tag in
  a search URL overrides exclusion for its synonym family, keeping matches blurred.
  Banned tags do not appear in the public tag selector except in admin mode.
- `censor` retains matching books and blurs their cards and direct content until
  the reader deliberately reveals that item. Revealing does not alter settings.
- Matching includes synonyms. Invalid modes/IDs fail before saving; empty is valid.
- Anonymous readers have no tag preferences. Signed-in preferences load freshly
  on each request; failures do not silently disable filtering.
- Admin routes and management tools bypass this policy. Public `admin.link=1`
  and `admin.override=1` bypass it only for authenticated admins/system users.
  Ordinary admin browsing follows their preferences. Availability checks remain.
- Apply to search, home books, profiles, collections, library, random selection,
  book/TOC/preview and chapter routes. Saved membership is never removed.

Plan: resolve and validate IDs in app; attach request preferences after auth;
apply shared hidden-tag constraints in SQL/OpenSearch; record blur decisions in
request context for reusable templ wrappers; enable searchable settings controls.
Verify matching, modes, validation, direct access, query/count parity and admin
bypasses with focused tests; generate templ, build frontend, run Go checks.

Tasks:
- [x] Settings selector and validation.
- [x] Shared matching and authorized request scope.
- [x] SQL/OpenSearch exclusions and consistent counts/random selection.
- [x] Blur/reveal on every public book surface; blurred direct-page response.
- [x] Regression tests, generation, build, and documentation.

Verification: focused Go tests and `git diff --check` pass. SQLC and templ
sources regenerated. Frontend build remains blocked by the existing missing
`block-editor/contracts` module and unused `useEffect` in `UserMenu.tsx`.
Browser/deployment checks against PostgreSQL and OpenSearch were not completed.
The reveal control uses a compact centered button on narrow cards.
