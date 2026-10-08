# Tasks

- [x] Add admin table/detail API and SPA routes (1, 2).
- [x] Implement role-checked book/chapter override and warnings/navigation (3, 6).
- [x] Move moderation to DashboardShell with admin navigation (4).
- [x] Add public search admin links and flag persistence (5, 6).
- [x] Add focused regressions, document and verify (1–6).

- [ ] Complete `pnpm run build` after unrelated baseline TypeScript failures are resolved.

## Verification

- `go test ./...` passes with writable temporary Go caches, including ID/URL parsing, role enforcement, restriction reasons, hidden chapter override, admin API pagination parameters, and conditional public-search links.
- `templ generate` and SQLC generation completed. The installed templ CLI reports a version mismatch with go.mod.
- `pnpm exec vite build --mode production` passes.
- `pnpm run build` stops at existing TypeScript errors: missing `block-editor/contracts` and unused `useEffect` in `header-user-menu/UserMenu.tsx`. No new TypeScript errors are reported.
- `git diff --check` passes.

- [x] Style search admin chips and add/test the restricted-book override prompt (7).

Prompt regressions pass for admins/system users, moderators, ordinary users, anonymous visitors, and missing books. The tests verify that restricted content is not fetched again before opt-in and that existing URL parameters survive the override link.

- [x] Refine chip spacing and pagination; verify desktop/mobile layouts and next-page navigation in Playwright (8).

Playwright verified at 1440px and 390px widths: chip height 34px, no wrapping, text-column alignment, no horizontal overflow, and page 5 → 6 navigation retaining `admin.link=1`.
