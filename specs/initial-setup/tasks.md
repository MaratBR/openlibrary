# Tasks

- [x] Add persisted SiteSetup configuration, defaulting IsSetup to false.
- [x] Add and generate the administrator setup-state query.
- [x] Resolve the administrator condition.
- [x] Implement SetupService and register its provider.
- [x] Remove middleware-owned caching in favor of service caching.
- [x] Test setup, state checks, caching, and Fx wiring.

- [x] Render the translated setup form and password requirements.
- [x] Handle CSRF-protected submissions, validation, failures, and completion.
- [x] Run handler tests for new and existing admins, CSRF, failures, redirects, and method handling.
- [x] Verify Vite asset generation and final checks. Full pnpm build is blocked by the unrelated unused useEffect import in UserMenu.tsx.
