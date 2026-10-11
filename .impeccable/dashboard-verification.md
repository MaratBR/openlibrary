# Dashboard redesign verification

Verified against the authenticated local application on 2026-10-11.

- Six routes: administration home/users, book manager list/detail, moderation
  books/users. All captured at 1440 × 1000 and 390 × 844 in light and dark themes.
  Document width equals viewport width in all 24 combinations. Administration
  tables retain internal horizontal scrolling on phones.
- Administration, book manager, and moderation sidebar navigation updates the
  route and active state. The book action menu opens; Escape closes it and returns
  focus to the trigger. New-book navigation renders its existing form.
- Existing input components, input styles, copy, and domain behavior are preserved.
  The moderation user filter container switches layout at the wider breakpoint
  required by the sidebar. No mutations were exercised.
- `pnpm exec vite build --mode production`: passed.
- `pnpm exec tsc --noEmit --noUnusedLocals false`: passed.
- `pnpm run build`: blocked by a baseline unused `useEffect` import in
  `web/frontend/src/islands/header-user-menu/UserMenu.tsx:1`, also present in HEAD.
  That unrelated file was not changed.
- `git diff --check`: passed. Mechanical detector returned no findings.
- Independent finish reviewer inspected all 24 settled screenshots and the
  scoped source changes; disposition: **ship**, no material fixes required.

Screenshots are local evidence under the ignored `.impeccable/review/` directory.
