# Verification: classification service and editor contract

2026-10-08:

- `GOCACHE=/tmp/openlibrary-go-cache go test ./...`: passed.
- Classification tests cover every supported rating, invalid/unrated values,
  non-adult warnings, and an adult tag raising every rating to adult.
- Custom-provider tests verify both public and manager details use the injected
  service. A warning regression verifies tag rules cannot override its result.
- HTTP tests cover legacy `isAdult` true/false at G, R, and NC-17. Each returns
  HTTP 400 before reaching persistence, avoiding partial saves.
- `git diff --check`: passed.
- `pnpm run build`: blocked by existing errors in unchanged files:
  `web/frontend/src/block-editor/state.ts:5` imports missing `./contracts`;
  `web/frontend/src/islands/header-user-menu/UserMenu.tsx:1` has unused `useEffect`.
  Both are present in HEAD. No production frontend build or browser verification
  was completed.

The `adult` response now includes assigned adult tags and excludes R by itself.
Legacy update clients must omit non-null `isAdult`. The service can be replaced
through Fx to customize classification; HTTP handlers contain no rating/tag rule.

Broader account filtering, descriptor selection, discovery/random consistency,
direct chapter protection, interstitials, scoped exceptions, saved-item handling,
and preference migration remain pending in tasks.md.
