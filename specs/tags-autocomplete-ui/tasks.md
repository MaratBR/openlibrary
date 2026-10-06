# Tasks

- [x] Update control layout, suggestions, and keyboard/focus behavior.
- [x] Add scoped styles and translated feedback at all existing call sites.
- [x] Regenerate templates, run the asset build and template tests, and check whitespace.
- [ ] Pass the full frontend build: blocked by existing missing block-editor/contracts module and unused UserMenu useEffect import.

Verification: `pnpm exec vite build --mode production`, `templ generate`,
`GOCACHE=/tmp/openlibrary-go-cache GOTMPDIR=/tmp go test ./web/public/templates`,
and `git diff --check` passed. `pnpm run build` reports only the two unrelated
TypeScript errors above. Keyboard and focus paths were reviewed in code;
interactive browser verification was not performed.

- [x] Accept genre tags in both runtime schema and generated DTO types.
- [x] Keep suggestions closed for empty or whitespace-only input.
- [x] Verify Martial Arts and all server categories with `node scripts/tags-search.test.mjs` (2 checks passed).
