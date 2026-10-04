# Implementation verification

Implementation was authorized on `feature/sdd/bm-1` on 2026-10-03. Verification
continued on 2026-10-04. The checks below use the current built Book Manager island,
common CSS, repository logo/fonts/icons, English catalog, and fixture API responses.
They do not exercise a database or a signed-in server session.

## Reproduce the browser checks

The suite requires Python 3.11 or newer, Playwright, and a compatible Chromium
binary. It can use a standalone `playwright` installation or the Playwright
package provided by the installed `@playwright/mcp` tool. No new application
runtime dependency is added. Build assets first:

```sh
pnpm exec vite build --mode production
node specs/book-manager-refresh/checks/browser.cjs
```

If the installed Playwright package and available browser versions differ, set
`PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH` to the available Chromium executable. The
suite starts its own temporary loopback server, intercepts API requests with
fixture responses, and closes the server/browser when finished. Screenshots are
written to a temporary directory whose path is printed at completion.

## Coverage

| Criteria | Verification |
| --- | --- |
| AC-1 | Library, overview, details, chapters, and wizard at 360/768/1280px in both themes; document width does not exceed the viewport |
| AC-2, AC-3 | One/20 cards, long titles, missing/broken covers, singular counts, empty library, pagination, failed list and retry |
| AC-4 | Trash failure retains state/dialog; pending disables submission; confirmed success refreshes; Restore and Cancel |
| AC-5 | Default overview and legacy Analytics links show summary/metadata; no-summary state |
| AC-6, AC-7 | Whitespace Enter validation, server-provided ratings, optional tags, review, completed-step navigation, retained values, and one explicit native creation POST |
| AC-8 | Changed title, Stay, failed save retaining input, confirmed save, clean Cancel; focused lower fields remain above the sticky action bar |
| AC-9, AC-10 | Chapter list, scheduled timezone copy, Add chapter validation parity, failed/successful creation and mobile popover fit |
| AC-11 | Panel fits all three widths; input focus, Escape, trigger focus restoration, changed-field Stay, failed and successful save |
| AC-12 | Named controls and native keyboard actions; shared toolbar buttons, tag options, and pagination use native semantics; both themes |
| AC-13 | Wizard discard, details route blocking, chapter dismissal prompts, clean forms after successful saves |
| AC-14 | Default redirect, catch-all, old tab links, list/detail recovery and unchanged editor URL construction; live authorization remains unverified |

The suite also simulates a server response that does not confirm the requested
Adult setting. The form retains input and reports the mismatch rather than
reporting all fields saved.

## Build and server checks

- `pnpm exec vite build --mode production`: passes. Existing asset-resolution,
  Sass import, and Browserslist notices remain.
- `go test ./web/public/...`: passes with writable temporary `GOCACHE` and
  `GOTMPDIR`.
- `templ generate -f web/public/templates/book_manager_home.templ`: completed.
  The installed generator reports an existing version mismatch with go.mod.
- `git diff --check`: passes.
- `pnpm run build`: blocked by unrelated existing branch errors:
  `web/frontend/src/block-editor/state.ts` imports missing `./contracts`, and
  `web/frontend/src/islands/header-user-menu/UserMenu.tsx` has an unused
  `useEffect` import. No errors were reported in the changed TypeScript files.

## Remaining live checks and existing contract limits

The local application was unavailable during implementation verification. A
signed-in session with disposable data is still needed to verify database
persistence, authorization redirects, and the actual block-editor destination.
The block editor implementation is outside this refresh.

Source review found that `bookDirectUpdate` accepts `isAdult` but does not pass it
into `UpdateBookCommand`. This existing backend behavior was not changed. The UI
reports when the response does not confirm that setting; other submitted fields
may have persisted. Reconcile the adult/rating contract in a separate backend
change before relying on the independent Adult switch.

Native creation failure still renders a separate server error page. Wizard input
recovery after that response and reconciliation of the client’s 70-character
chapter-name limit with the server’s 200-character limit remain outside this scope.
