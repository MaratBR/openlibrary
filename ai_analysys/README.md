# OpenLibrary engineering audit

Date: 2026-08-06

## Scope and method

This is a maintainability, code-quality, reliability, and production-observability review of the Go server, TypeScript/JavaScript frontend, and SCSS/CSS organization. It deliberately does not evaluate UX and does not treat placeholders or incomplete product features as defects merely because they are unfinished.

The review sampled request routing and response helpers, application services, persistence boundaries, background workers, external clients, telemetry, frontend entry points, API contracts, island lifecycle code, build configuration, and the shared style graph. It also ran the repository's standard Go and frontend checks.

## Executive summary

The repository has a sound intended architecture: HTTP, application, and SQLC persistence code are visibly separated; Go contexts are commonly propagated; database, Redis, and HTTP tracing has begun; TypeScript is strict; server DTO generation exists; styles have shared tokens and components. Those are useful foundations.

The largest risk is not unfinished feature code. It is that failures are difficult to classify and safely diagnose in production. Panic and ordinary error responses can disclose internals, there is no readiness/health contract, startup-critical work is detached in goroutines, logging is split across three mechanisms, and traces have a hard-coded exporter but no metrics/alerting story. Several background jobs are in-memory and best-effort, so a process crash can create silent search/index drift.

On the frontend, compile-time typing is stronger than runtime contract safety: most API calls invoke a generic helper whose default schema is `z.any()`. Browser-side exceptions are generally written only to the console. On styling, the token/component direction is good, but the global style entry imports every component and page stylesheet, has no automated linting, and currently produces an unusually large shared CSS artifact.

## Priority order

1. Fix unsafe and inconsistent server error handling, especially panic recovery and 500 responses.
2. Add health/readiness endpoints and make startup dependency initialization explicit.
3. Standardize structured logging and correlate logs with request and trace IDs.
4. Make secondary-index work durable and observable rather than detached best-effort goroutines.
5. Repair the Go verification baseline and make both Go and frontend checks mandatory in CI.
6. Enforce runtime API schemas and centralize frontend error reporting.
7. Establish explicit frontend/style boundaries and budgets as the codebase grows.

Detailed findings are in:

- [Go, reliability, and observability](go_reliability_observability.md)
- [TypeScript and JavaScript](typescript_javascript.md)
- [SCSS and CSS](scss_css.md)

## Verification snapshot

- `pnpm run build`: passed. Vite reported an unresolved runtime image URL and stale Browserslist data. The production build was unminified, emitted sourcemaps, and produced `dist/common.css` at 972.64 kB (595.97 kB gzip).
- `go test ./...`: failed because Go's vet phase rejects dynamic format strings passed to `errorx.Type.New/Wrap` in `internal/olhttp/query.go`, `internal/app/comments_impl.go`, and `internal/app/moderation_book_impl.go`. Other packages that reached execution passed.
- No project-owned frontend test files were found. The frontend build script only runs TypeScript and Vite; package metadata contains Oxlint and Prettier but exposes no lint, formatting-check, or test script.

## Suggested implementation sequence

### First milestone: safe failure handling

Introduce one response/error boundary that maps typed application errors to status codes, logs unexpected errors once with request/trace correlation, returns a stable public error envelope, and never returns stack traces or raw internal errors. Add middleware tests for panic-before-write, panic-after-partial-write, API versus HTML responses, and request ID propagation.

### Second milestone: operable service lifecycle

Define liveness and readiness separately. Readiness should cover database connectivity, completed search schema setup, and required object-store buckets. Move required initialization into Fx lifecycle hooks with bounded contexts. Make optional dependencies explicitly degraded and expose that state.

### Third milestone: durable asynchronous work

Replace detached reindex goroutines with a durable outbox/job table written in the same transaction as source data. Track queue depth, oldest-job age, attempts, failures, and last successful completion. Apply the same worker lifecycle pattern consistently.

### Fourth milestone: contract and boundary enforcement

Make response schemas mandatory in the frontend API layer, add tests around the HTTP wrapper and island lifecycle, split global/page/island styles, and add lint/test/bundle-budget commands to CI.
