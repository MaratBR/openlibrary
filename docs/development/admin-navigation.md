# Admin workspace

The admin UI is a React island at `/admin`, using the same hash routing approach
as Book Manager. The entry template is `web/admin/templates/spa.templ`; the root
router and screens are under `web/frontend/src/features/admin`. The admin Alpine
entry imports the shared island registry before starting Alpine and registers
`admin/App`. Styles live in `web/frontend/src/features/admin/styles`.

The shared dashboard shell places branding, section links, theme switching,
site return, and logout in a persistent left sidebar. On small screens the
navigation wraps above the content. Styling is shared with Book Manager in
`web/frontend/src/components/dashboard-shell.scss`.

Use React Router links for workspace navigation. Current routes are `#/`,
`#/users`, `#/users/:id`, `#/books`, `#/tags`, `#/tags/:id`,
`#/tags/:id/edit` and `#/debug`. Search, filters and pagination live in the hash
query string. Legacy server page URLs redirect to their equivalent hash routes.
Books retain the existing numeric-ID lookup and link to public book pages.

JSON endpoints under `/admin/api` use the existing user, tag and reindex services.
All endpoints are inside the authorization group; authorization checks the
freshly loaded account's admin/system role. POST endpoints retain the server's
CSRF middleware and the shared frontend client's CSRF header. DTO IDs are
strings; validate responses using the schemas in `features/admin/api/index.ts`.

Put route reads in exported loaders and pass the loader request's AbortSignal.
Keep mutations explicit: disabled pending controls, retained input on failure,
and success only after a confirmed response. Dirty forms use the local
Stay/Discard dialog and a native browser unload warning. Maintenance actions
require confirmation before scheduling. Do not advertise settings, tag creation
or deletion, or book search until backend operations exist.

After editing, run `templ generate`, `go test ./web/admin/...`,
`pnpm run build` and `git diff --check`. Browser fixture checks are in
`specs/admin-spa/checks/browser.cjs`; they cover navigation, payloads, loading and
failure states, dirty forms and both themes at 360/768/1280px. Fixture checks do
not establish real database persistence; verify that separately on local data.
