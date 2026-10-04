# Verification

- `templ generate`: passed; installed generator reports it is older than the module version.
- `GOCACHE=/tmp/openlibrary-preview-go-cache GOTMPDIR=/tmp go test ./web/public/templates`: passed.
- `pnpm exec vite build --mode production`: passed.
- `node specs/book-card-preview-refresh/checks/browser.cjs`: passed using local Chromium outside the filesystem sandbox. Checks stable short-content dimensions, viewport edges, resize, late responses after mouse leave, and disposal. Uses synthetic preview HTML with the actual island and built CSS.
- `git diff --check`: passed.
- `pnpm run build`: blocked by unrelated TypeScript errors in `src/block-editor/state.ts` (missing `./contracts`) and `src/islands/header-user-menu/UserMenu.tsx` (unused `useEffect`).

The preview retains up to five tags, omits absent author/summary/tags, and uses a decorative cover fallback. Long summaries are clipped in CSS instead of cutting through HTML or UTF-8 bytes.

Follow-up checks passed: right/left/top/bottom placement in controlled viewports; 500 ms touch hold; short tap and movement cancellation; centered card with computed blur; prevented navigation; native Chromium touch release keeps the card open; backdrop and Escape dismissal restore scrolling; library navigation fits 320/375/768 px with English and Russian labels and retains active styling. Vite build and whitespace checks passed again. Full type-check remains blocked by the same two unrelated errors above.
