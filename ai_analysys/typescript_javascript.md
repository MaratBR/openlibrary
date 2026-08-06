# TypeScript and JavaScript findings

## High

### TS-01: Runtime response validation is optional and therefore usually bypassed

Evidence: `OLAPIResponse.create` defaults its schema to `z.any()` (`web/frontend/src/http-client/OLAPIResponse.ts:16-22`). Many API modules provide only a TypeScript generic, for example comment endpoints (`web/frontend/src/api/comments.ts:14-52`) and moderation user endpoints (`web/frontend/src/api/moderation/users.ts:14-55`). A generic is erased at runtime, so malformed or drifted server data passes through unchanged.

Impact: contract drift fails later inside UI code, far from the network boundary, with less useful context. Generated DTO types provide compile-time documentation but not runtime protection.

Recommendation: remove the default schema and require a Zod schema (or a generated runtime schema) for every JSON response. Validate error envelopes separately. Generate TypeScript types from those schemas or generate both from the Go contract so the representations cannot drift independently.

### TS-02: Error responses can be parsed with the success schema before callers can inspect them

Evidence: `_loadData` identifies 4xx/5xx responses but then always parses the same JSON through the success schema (`web/frontend/src/http-client/OLAPIResponse.ts:68-78`). With real schemas, an expected API error can throw a Zod error instead of producing `OLAPIError`. `throwIfError` is opt-in (`:62-66`), and different API modules return raw Ky responses, parsed JSON, or `OLAPIResponse`.

Impact: callers see inconsistent failure types; error notifications and status metadata can be lost; retries and user-safe reporting become hard to centralize.

Recommendation: define a discriminated result pipeline. Parse non-2xx bodies with one error schema and throw a typed transport/API error automatically; parse 2xx with the endpoint schema. Normalize invalid JSON, timeout, abort, offline, and schema mismatch into distinct error classes carrying request ID and endpoint/operation metadata.

### TS-03: Browser failures are not observable outside developer consoles

Evidence: the shared logger is only a styled console wrapper (`web/frontend/src/logger.ts`). Island loading/data errors are logged and swallowed (`web/frontend/src/lib/island/island-element.ts:76-110,168-177,196-208`). React's error boundary also only calls `console.error` (`web/frontend/src/preact/wrapper.tsx:31-35`). No global `error`/`unhandledrejection` reporting or telemetry client was found.

Impact: production-only chunk-load failures, contract errors, and render crashes are invisible unless a user supplies console output. Server request IDs are not carried into browser reports.

Recommendation: define a small frontend telemetry interface and production sink with release, route, island/component, request ID, stack, and sampled breadcrumbs. Report error-boundary, unhandled rejection, dynamic import, API schema, and asset-load failures. Keep console as the development sink and scrub user content.

### TS-04: Overriding `window.fetch` creates a global hidden dependency

Evidence: importing the HTTP client replaces `window.fetch` for the entire page to add CSRF (`web/frontend/src/http-client/client.ts:5-33`). It mutates `input.headers` when input is a `Request`, catches cookie access failure silently, and affects third-party and future code that never opted into the OpenLibrary client.

Impact: behavior depends on import order, tests require a browser-global patch, cross-origin calls may receive an application header and trigger CORS preflight, and debugging native fetch behavior becomes confusing.

Recommendation: put CSRF injection in Ky hooks on a configured first-party client, restrict it to same-origin state-changing methods, clone headers, and expose a separate raw/external client when needed. Do not patch the platform global.

## Medium

### TS-05: API organization has parallel conventions

Evidence: endpoints are split among `src/api`, `src/public.api`, and feature-local `api.ts`/`api` directories. Return styles vary among raw `Response`, `r.json<T>()`, and `OLAPIResponse<T>`. Some functionality is also attached to `window` for Alpine/server markup.

Impact: a contributor has to discover which layer and response convention a new endpoint should use; cross-cutting behavior such as auth, validation, notifications, and telemetry is applied unevenly.

Recommendation: define one endpoint module convention: URL construction, request/response schemas, normalized result/error, and ownership by domain. Keep the `window` surface as a thin compatibility adapter over those modules and document every stable global.

### TS-06: The island lifecycle has race and cleanup gaps

Evidence: when an island becomes inactive before its async load completes, the animation-frame branch returns without resetting `_isCreating` (`web/frontend/src/lib/island/island-element.ts:95-109`), preventing a later reactivation from creating it. `disconnectedCallback` destroys the island but does not remove the event listener added by `connectedCallback`, and reconnecting adds it again. `_preload` starts an import without rejection handling (`:196-201`). Mount exceptions happen inside an animation-frame callback and are not caught by the async `_create` catch.

Impact: WIP navigation and dynamic markup can produce islands that never remount, duplicate listeners, and unreported asynchronous exceptions.

Recommendation: model lifecycle states explicitly, use a generation/abort token for stale loads, always reset state in a complete `finally`, pair listener registration/removal, catch mount/dispose/preload failures, and add DOM tests for connect-disconnect-reconnect and active toggling during load.

### TS-07: There is no project-owned frontend test or enforced lint command

Evidence: no test files outside dependencies were found. `package.json:6-9` exposes only dev, build, and preview. Oxlint and Prettier are dependencies (`:28-30`) but have no scripts or checked configuration in the main workflow.

Impact: strict TypeScript catches type errors but not runtime API parsing, DOM lifecycle, state transitions, or formatting/lint regressions. Complex files such as the ~927-line moderation page and ~454-line comments island have no focused safety net.

Recommendation: add a fast unit/DOM test runner and target high-value seams first: HTTP response/error parsing, island lifecycle, editor state reducers, and moderation mutations. Add `lint`, `format:check`, and `test` scripts and run them in CI.

### TS-08: Production build policy is embedded as development-friendly constants

Evidence: Vite always sets `SOURCEMAP = true`, `minify: false`, and stable unhashed entry filenames (`vite.config.ts:8,123-150`). The verified production build emitted very large readable chunks and public map files.

Impact: deploy artifacts are larger, implementation is easier to inspect, cache invalidation depends on external discipline for entry files, and accidental old/new asset mixing is more likely. Sourcemaps are valuable for observability but should normally be uploaded privately to the error service.

Recommendation: make production minification, source-map publication/upload, and asset naming explicit environment policy. Use a manifest and content hashes or guarantee atomic deploy/cache headers for stable entry names. Add bundle budgets and a build report in CI.

### TS-09: Shared entry initialization assumes browser APIs without fallbacks

Evidence: `web/frontend/src/common/index.ts:25-27` calls `requestIdleCallback` directly despite an ES2020 browser target, and global initialization is performed via import side effects.

Impact: unsupported browsers or non-browser test environments can abort common initialization; implicit ordering complicates tests.

Recommendation: wrap optional browser APIs with tested feature-detection fallbacks and expose an explicit idempotent bootstrap function.

## Existing strengths worth preserving

- `strict`, unused-symbol checks, isolated modules, and side-effect import checking are enabled in `tsconfig.json`.
- The `@/` alias and domain-oriented directories make most imports understandable.
- React islands have a shared QueryClient and visual error boundary.
- Zod is already adopted in selected boundaries; expanding that pattern is incremental.
- Dynamic imports allow islands to remain separate from the common page bundle.
