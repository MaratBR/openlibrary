# Book Manager React island guide

This file applies to the Book Manager island under this directory. Follow the repository-level contributor guide as well.

## Entry point and routing

- `BM.tsx` is the root component for the `BMIsland` exported by `index.ts`. Keep the component name exactly `BM`; server-side island registration depends on the exported island wiring.
- The Book Manager is a hash-routed React application. Define top-level routes in `BM.tsx`, render them inside `BMLayout`, and use React Router navigation components instead of full-page links for routes within this island.
- Put route data fetching in exported loaders next to the routed screen. Validate route parameters before using them, and type `useLoaderData` from the loader's return type.
- Preserve the `/books` default redirect and the catch-all route when changing the route tree.

## Data and state

- Use the typed Book Manager clients under `@/api/bm` for book operations. Keep their Zod schemas synchronized with the complete server DTOs consumed by these screens; Zod strips undeclared fields.
- Call `throwIfError()` before treating mutation responses as successful. Keep pending, success, and failure behavior explicit for user-triggered mutations.
- Prefer loader data for route-level reads, React Query for server mutations/caching, and local React state for view-only state. Add Zustand state only when data genuinely needs to be shared across multiple Book Manager screens.
- Keep island-specific contracts in `contracts/`. Reuse shared DTOs and controls from `@/api`, `@/components`, and `@/lib` rather than defining parallel versions here.

## UI conventions

- Use `DashboardContent` and `BMLayout` for the dashboard shell. Reuse shared controls such as `Modal`, `Pagination`, `Tabs`, `BookCover`, and established `.btn`, `.card`, `.chip`, and `.table` classes.
- Use UnoCSS utilities for local layout and spacing. Add island-scoped styling only for a repeated Book Manager pattern that shared components do not cover.
- All visible copy must use `window._(...)`. Add or update English keys under the `bookManager` namespace in `translations/en.toml` and keep other catalogs consistent with the repository workflow.
- Preserve accessible control semantics: use buttons for actions, router links for navigation, labels for form fields, and meaningful loading/disabled states during submissions.

## Verification

- Run `pnpm run build` after changing routes, components, contracts, or frontend API usage.
- Run `git diff --check` and report unrelated baseline failures without modifying unrelated files.
- For behavior with nontrivial state transitions or regressions, add a focused test when the surrounding frontend test setup supports it; otherwise describe the manual route and interaction checked.
