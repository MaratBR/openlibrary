# Debugging and visual checks

Use this document for developer and coding-agent diagnostics. Feature behavior
belongs in the corresponding feature documentation.

## Debug controller

Public diagnostic routes belong in
[`web/public/routes_debug.go`](../../web/public/routes_debug.go), registered by
`web/public/handler.go` through Fx. Keep startup data in
[`web/webinfra/server_runtime.go`](../../web/webinfra/server_runtime.go);
`cmd/server/server_module.go` populates it after the HTTP listener binds and
before serving requests. Add debug endpoints to the controller rather than
intercepting requests in root middleware.

| Endpoint | Behavior |
| --- | --- |
| `GET /debug/launch-time` | Plain-text UTC ISO timestamp (`RFC3339Nano`) for the current HTTP listener startup. Unchanged between requests, with `Cache-Control: no-store`. |
| `HEAD /debug/launch-time` | Same headers, without a response body. |
| `/debug/500` | Deliberately renders the existing generic error page for visual testing. |

`/debug/launch-time` is publicly readable. Compare it before and after a server
restart to confirm which process the browser is reaching:

```sh
curl -fsS http://localhost:8080/debug/launch-time
```

Other HTTP methods on `/debug/launch-time` return 405 through the normal router.
These routes use the normal public request middleware.

## Playwright workflow

Use Playwright to inspect actual rendered UI when changing layouts or styles.
Begin with the current browser session; it may already be signed in for admin
checks. Take screenshots at desktop and phone widths (for example, 1440 × 1000
and 390 × 844), then inspect element bounds, wrapping, margins, keyboard focus,
and horizontal overflow. Use accessible locators to exercise navigation and
confirm that the URL, active page, and required query flags update together.

For admin book search, open `/search?admin.link=1`. Check chip alignment with the
book text, spacing above and below, and single-line labels. Exercise pagination
at its beginning, middle, and end; verify previous/next controls, ellipses, and
preservation of `admin.link=1` and search filters.

If the DOM still contains old markup, first check `/debug/launch-time`. Template
changes require `templ generate` and a Go server rebuild/restart; frontend
changes require the Vite watcher or a fresh asset build. A manually started
server does not reload when files change. Recheck the browser after the server
is current rather than treating a stale screenshot as a completed verification.

For deterministic book-page fixtures, the existing template tests can export
HTML into a temporary directory:

```sh
mkdir -p /tmp/openlibrary-book-fixtures
BOOK_PAGE_FIXTURE_DIR=/tmp/openlibrary-book-fixtures go test ./web/public/templates
```

Use fixtures for layout-only checks instead of modifying real book status just
to obtain a screenshot. Report whether checks used the live application or a
fixture, and record any baseline build failures separately from the change.

## Verification

```sh
go test ./cmd/server ./web/public/... ./web/webinfra
pnpm run build
git diff --check
```

When Go cache directories are not writable, use `GOCACHE` and `GOTMPDIR` under
`/tmp`, as described in the repository contributor guide.
