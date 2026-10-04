# Implementation tasks

- [x] Review scoped guides and existing route loading/dirty-form behavior (AC-1–4).
- [x] Add translated masthead site link with unsaved-change handling (AC-1, AC-2).
- [x] Build responsive skeleton compositions using the shared primitive (AC-3, AC-6, AC-7).
- [x] Wire initial and destination loading, keeping revalidation nonblocking (AC-3, AC-4).
- [ ] Confirm success, empty, failure, and retry transitions (AC-5).
- [ ] Verify keyboard, announcements, reduced motion, themes, and viewport sizes (AC-1, AC-2, AC-6, AC-7).
- [ ] Run implementation checks and record evidence for every criterion.

## Implementation and verification evidence

Implemented in `BM.tsx`, `BMLayout.tsx`, `ui.tsx`, and new `loading.tsx` /
`loading.scss`; link copy is in `translations/en.toml`. Existing form and manager
stylesheet edits from the concurrent UI task were preserved.

- AC-1: Source review confirms the shell contains a real `/` anchor alongside
  Your books. Loader fallbacks and route errors remain inside that shell.
- AC-2: The existing dirty-form hook handles the masthead anchor with its
  Stay/Discard dialog and unload bypass after Discard. Primary keyboard clicks
  and modified primary clicks follow the existing guarded editor-link convention;
  clean anchors keep native behavior. Browser interaction verification remains.
- AC-3: A temporary server-rendering fixture passed for `/books`, `/books/one`,
  `/books/one/edit`, and `/books/one?t=chapters`, asserting the corresponding
  grid, metadata, form workspace, and chapter composition. These are fixture
  checks, not live API or persistence checks.
- AC-4: Source review confirms pathname/search changes select the destination
  skeleton and hide/inert the previous route and its portal dialogs. Same-route
  revalidation retains mounted content with a loading status. Book overview routes
  are keyed by book ID to reset chapter panel state between books; the existing
  details form already has a book-ID key. Slow navigation and retained edits
  still need browser verification.
- AC-5: Existing loaders still throw failures into RouteError; Retry reloads the
  document and re-enters initial loading. Empty states render only after loader
  success. No artificial delay was added. Success/error/retry transitions still
  need browser verification.
- AC-6: The rendering fixture verified one status, `aria-busy`, hidden skeleton
  content, and no anchor/input/button elements for all four compositions. Source
  review confirms the reduced-motion rule disables animation. Accessibility tree,
  announcement timing, and keyboard checks still need a browser.
- AC-7: Skeletons reuse BM-grid, BM-workspace, BM-cover dimensions, Card, and theme
  tokens. Masthead wraps its navigation. Both-theme visual checks at 360/768/1280px
  remain outstanding.

Checks:

- `git diff --check`: passed.
- `pnpm run build`: blocked by unrelated TypeScript errors in
  `web/frontend/src/block-editor/state.ts` (missing `./contracts`) and
  `web/frontend/src/islands/header-user-menu/UserMenu.tsx` (unused `useEffect`).
  No errors were reported in the changed navigation/loading files.
- Browser fixture could not run: sandbox denied the local Vite listening socket
  (`EPERM`), and the shared Playwright browser reported it was already in use.
  The existing browser session was left untouched to avoid disturbing the
  concurrent UI task. Do not mark the remaining verification tasks complete until
  the browser checks and full build can run.
