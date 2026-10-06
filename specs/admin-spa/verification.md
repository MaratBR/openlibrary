# Verification — 2026-10-04

Implementation is on `feature/admin-1`, branched from the user's current branch.
Frontend implementation is under `web/frontend/src/features/admin`; island and
asset entry points are thin wrappers. Obsolete SSR pages, fixed sidebar and
legacy modal/password-reset islands have been removed. Legacy page GETs redirect
into the SPA; legacy user/tag POSTs redirect after success.

## Live application

Inspected `http://localhost:8080/admin` using the existing authenticated admin
session. The home route now renders the workspace. React navigation between
existing screens works without page reloads. `Island` is registered before
Alpine starts, and the previously missing client flash host is present. No
runtime console errors occurred in the completed route checks.

Checked home, users, user editing, books lookup, tags, tag details, tag editing
and debug actions at 360, 768 and 1280px, in light and dark themes (48 route/theme/
viewport combinations). No document-level horizontal overflow or raw translation
keys were found. Inspected desktop home and mobile tag-edit screenshots.

Saved the existing Action tag with its unchanged description and unchecked flags,
and saved an existing user's unchanged account details. Both SPA submissions
reported confirmed success; tag details were read back successfully. No passwords,
roles or tag semantics were changed. Maintenance scheduling was tested only in
fixtures, not against the live index.

## Regression checks

- `go test ./...` with writable temporary Go cache: passed.
- Focused admin tests cover current-role authorization (including cached stale
  roles), anonymous API denial, admin/system access, CSRF rejection/acceptance,
  redirects, filter forwarding, empty arrays, invalid IDs/payloads, precision-safe
  IDs, optional password/custom gender, unchecked flags/empty descriptions,
  missing records and explicit debug actions.
- Browser fixture suite exercises the actual `admin-alpinejs.js` registration,
  SPA navigation, query preservation, pagination, back/forward/reload, parent-tag
  selection, optional password generation/cancellation, failed/successful saves,
  pending controls, dirty-form Stay/Discard, tag-filter return navigation,
  maintenance confirmation/cancellation, loading, retry, empty/unauthorized/
  invalid-route states and 48 themed responsive route layouts. Fixture mutations
  do not establish database persistence; the live unchanged-data saves above
  separately check the real submission path.
- `templ generate`: passed; installed generator reports v0.3.960 while go.mod
  specifies v0.3.1020. Generated files remain ignored.
- `pnpm exec vite build --mode production`: passed; generated assets remain ignored.
- `git diff --check`: passed.

The standard `pnpm run build` is blocked by pre-existing TypeScript failures:
`web/frontend/src/block-editor/state.ts` imports missing `./contracts`, and
`web/frontend/src/islands/header-user-menu/UserMenu.tsx` imports unused `useEffect`.
No TypeScript errors are reported in the admin feature. These unrelated files
were not changed.

Run the fixture checks with an installed browser, for example:

```sh
PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH=/opt/google/chrome/chrome \
  node specs/admin-spa/checks/browser.cjs
```
