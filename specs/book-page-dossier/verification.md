# Verification

- `templ generate`: passed; installed CLI reports it is older than the module.
- `go test ./web/public/...`: passed, including optional first-chapter action tests.
- `pnpm run build`: blocked by existing TS6133, unused `useEffect` import in
  `web/frontend/src/islands/header-user-menu/UserMenu.tsx:1`.
- `pnpm exec vite build --mode production`: passed.
- Live Playwright: 1440×1000 desktop, 390×844 mobile, and desktop dark theme;
  no document overflow, three actual chapter links, review tab opens, ArrowRight
  selects/focuses reviews, Home selects/focuses chapters, Start reading navigates
  to the correct chapter URL. Both themes use the actual website theme API.
- Existing `#client-flashes element not found` initialization error observed;
  its initialization code is outside the changed files.
- Independent screenshot/code review: ship; requested keyboard recheck passed.
- Impeccable detector: advisory typography, cover radii, shadow values; intentional
  book-surface choices documented separately. No mandatory detector findings.

Screenshots: `.impeccable/review/book/{desktop,mobile,dark}.png` (local artifacts).
Browser evidence uses the live anonymous page. Authenticated and empty states
remain covered by template tests rather than these live captures.
