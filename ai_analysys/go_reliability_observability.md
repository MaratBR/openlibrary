# Go, reliability, and observability findings

Severity indicates production and maintenance impact, not implementation difficulty.

## Critical

### G-01: Panic recovery leaks stack traces and may return HTTP 200

Evidence: `internal/olhttp/recovery_middleware.go:18-37` logs a panic, constructs a response containing `debug.Stack()`, and writes it directly to the client. It never calls `WriteHeader(500)`. If headers were not already written, the response therefore defaults to 200; if they were written, it can append a stack trace to a partial response. It also prints the whole request with `fmt.Printf`, bypassing structured logging and potentially including sensitive request metadata.

Impact: implementation details and filesystem/package information are disclosed; monitoring based on status codes misses panics; clients can cache or interpret failures as success; one incident is split across structured and unstructured output.

Recommendation: buffer dynamic responses or use a response writer that tracks commitment, set 500 when possible, return a small content-negotiated public error with request ID, and emit the stack only to the structured server logger. Record the panic and error status on the active span. Never return stack traces outside an explicit development-only mode.

### G-02: Ordinary server errors expose raw internal error text

Evidence: `web/public/handler_response_html.go:7-24` writes `err.Error()` for 400, 409, and 500 responses. API helpers pass the original error to `olhttp.NewAPIError` at `web/public/handler_response_api.go:10-40`; the response convention should be audited so unexpected errors cannot serialize database or infrastructure detail. Similar direct writes exist in admin handlers.

Impact: database, validation, dependency, or filesystem details can escape to clients. Error wording becomes an accidental public API, making refactoring harder.

Recommendation: create a single typed error mapper. Public errors should have a stable code, safe message, request ID, and optional field violations. Log the wrapped internal cause once at the boundary. Reserve raw causes for development logs.

## High

### G-03: The process can be serving before required dependencies are initialized

Evidence: `cmd/server/server_module.go:83-85` launches `postInit` in a goroutine; `postInit` sleeps two seconds and then launches three more untracked goroutines (`:160-186`). Bucket initialization failure is logged but does not affect service readiness. OpenSearch setup is also launched asynchronously and only logs failure (`cmd/server/infra_module.go:191-196`). No health, liveness, or readiness handlers were found.

Impact: deployment reports healthy because the TCP listener is up while uploads or search fail. Operators have no machine-readable way to distinguish starting, ready, and degraded. Shutdown does not wait for these goroutines.

Recommendation: classify dependencies as required or optional. Put required initialization into Fx `OnStart` hooks with timeouts and returned errors. Add `/livez` and `/readyz`; readiness should remain false until migrations/index setup/buckets complete and should check DB connectivity. Represent optional dependency failure as explicit degraded state.

### G-04: Search indexing is best-effort, process-local, and can silently drift

Evidence: `internal/app/book_search_reindex_impl.go:48-55` ignores the caller context, starts an untracked goroutine, and only logs a failed reindex. There is no persistent retry or dead-letter state. Full reindex deletes the live index before rebuilding it (`:109-121`), ignores JSON encoder errors (`:172,189`), and only checks transport-level bulk errors (`:193-201`), not per-document bulk failures.

Impact: a crash between database commit and indexing loses work permanently. A failed full rebuild can leave search empty or partially populated. Operators cannot query backlog, age, or failed document IDs.

Recommendation: use a transactional outbox/job table, bounded worker pool, idempotent jobs, retry policy, and dead-letter visibility. Rebuild into a versioned index and atomically swap an alias. Parse bulk responses and report item failures. Add queue-depth, oldest-age, success/failure, duration, and index-freshness metrics.

### G-05: Logging is fragmented and weakly correlated

Evidence: the application mixes global `slog`, Zap, `fmt.Print*`, standard `log`, and direct stderr writes. `cmd/server/log.go` always configures a colorized development console encoder at debug level. Request IDs exist, but most log calls do not include them; trace IDs are not injected into log records. Some worker error messages omit the actual error, for example `internal/app/analytics/event_processor_worker.go:116,129,137`.

Impact: production ingestion and querying are inconsistent, noisy debug logs increase cost, and reconstructing a request or job failure requires guesswork.

Recommendation: choose one structured logging facade, configure JSON/level/output by environment, and derive request-scoped loggers from context with request ID, trace ID, route, component, and operation. Establish event names and required fields. Include errors consistently and avoid logging high-cardinality payloads or raw SQL by default.

### G-06: Telemetry configuration is hard-coded and incomplete

Evidence: `cmd/server/telemetry.go:16-35` always creates an insecure exporter for `localhost:4318` and fixes the service name to `openlibrary`. The process panics if exporter construction fails (`cmd/server/server_module.go:49-54`). HTTP, pgx, Redis, and some application spans exist, but no operational metric exporter, service version/environment resource attributes, sampling configuration, or alert-oriented health signals were found.

Impact: container and hosted deployments need code changes or sidecar-specific networking; telemetry can become a startup dependency unintentionally; traces cannot easily be separated by environment/version; absence of RED metrics makes regressions harder to detect.

Recommendation: support standard OTEL environment variables, optional telemetry, TLS/auth configuration, resource attributes (service version, environment, instance), and sampling. Add request rate/error/duration, DB pool saturation, external dependency latency/errors, worker health/backlog, and business pipeline freshness metrics. Define dashboards and alerts from SLOs rather than relying on raw logs.

### G-07: HTTP server lifecycle misses important failure and timeout handling

Evidence: `cmd/server/server_module.go:193-199` configures read and write timeouts but no `ReadHeaderTimeout` or `IdleTimeout`. `srv.Serve(ln)` is detached and its returned error is discarded (`:202-209`). Database connection retry uses `context.Background()` and an unbounded loop (`cmd/server/infra_module.go:76-101`).

Impact: slow-header attacks retain connections longer than necessary; an unexpected accept-loop failure can leave the process alive but non-serving; invalid database configuration can hang startup forever and defeat orchestration deadlines.

Recommendation: set explicit header/idle timeouts, propagate unexpected Serve failures to the application lifecycle, use bounded startup contexts/backoff, and expose startup state. Close DB, Redis, and other clients through lifecycle hooks.

### G-08: OpenSearch TLS verification is always disabled

Evidence: `cmd/server/infra_module.go:175-184` sets `InsecureSkipVerify: true` regardless of environment or URL.

Impact: production TLS connections are vulnerable to interception and configuration mistakes are masked.

Recommendation: use normal certificate validation by default. If local development needs insecure TLS, require an explicit development-only flag and log it prominently.

## Medium

### G-09: Background-worker lifecycle and concurrency patterns are inconsistent

Evidence: workers independently implement cancellation, done channels, cron locks, timers, and goroutines. The analytics worker creates an unbuffered wake channel and uses `len(channel)` followed by a send (`internal/app/analytics/event_processor_worker.go:34-37`), which is not an atomic nonblocking send and can block with concurrent callers. Initialization panics inside its goroutine (`:73-82`), escaping Fx startup error handling. `stop` ignores its context while waiting (`:153-161`). `internal/http-lim/http-worker.go:48-66` can block forever when called before `Run`, race with `Close`, or send to a closed queue; its state is unsynchronized and request cancellation is not used while queueing/rate-limiting.

Impact: shutdown hangs, goroutine leaks, rare races, and failures that bypass lifecycle supervision are likely as load and features grow.

Recommendation: provide a shared supervised-worker pattern: lifecycle-owned context, `errgroup`, buffered/coalescing signals with `select { case ch <- x: default: }`, bounded shutdown, panic recovery at worker boundaries, and explicit health state. Add `go test -race` tests around wake/stop/close behavior.

### G-10: Large application services concentrate unrelated responsibilities

Evidence: `internal/app/book_manager_impl.go` is about 1,020 lines, while several public API route files approach 500 lines. The book-manager implementation combines reads, uploads, draft mutation, chapter ordering, publishing, metrics, and search reindex scheduling.

Impact: changes have a broad conflict and regression surface; transaction and authorization invariants are harder to see; tests require large dependency setups.

Recommendation: split by cohesive use case (catalog queries, draft editing, publishing, cover management), keeping transaction orchestration in small command services. Prefer narrow interfaces owned by consumers. Do this incrementally when touching an area, not as a flag-day rewrite.

### G-11: Render and response-write errors are usually discarded

Evidence: many handlers call `Component.Render(ctx, w)` without checking the returned error, such as `web/public/routes_book.go:131,164,184,210`, `web/public/routes_library.go:70,168`, and multiple admin routes. Direct `w.Write` results are also commonly ignored.

Impact: client disconnects and template failures disappear or produce partial responses without an error event, making production diagnosis incomplete.

Recommendation: centralize template rendering in a helper that renders before committing when practical, records failures with request correlation, and distinguishes client cancellation/broken pipe from server rendering defects.

### G-12: The Go quality gate is currently red

Evidence: `go test ./...` fails vet checks for non-constant format strings in `internal/olhttp/query.go:20`, `internal/app/comments_impl.go:42,220,231`, and `internal/app/moderation_book_impl.go:36`.

Impact: contributors cannot use the documented full test command as a reliable regression signal; real failures become easier to ignore.

Recommendation: use constant format strings (`"%s"`) or the correct error constructor semantics, then enforce `go test ./...`, `go vet ./...`, and preferably `staticcheck` in CI. Keep baseline green and quarantine only truly external/flaky checks with ownership and expiry.

### G-13: Configuration is stringly typed and inconsistently named

Evidence: infrastructure constructors pull keys directly from Koanf. Redis creation reads `redis.addr` (`cmd/server/infra_module.go:65-68`) while the committed config defines `redis.url`. Required Mailgun-key validation has an empty block (`:115-116`). Several invalid configurations call `panic` or `os.Exit` from providers instead of returning errors.

Impact: typos become runtime defaults, validation is scattered, and configuration errors have inconsistent behavior.

Recommendation: decode once into a typed config struct, validate it at startup, redact secrets in any diagnostic dump, and have Fx providers return errors. Add config parsing tests for the committed example and production-required combinations.

## Existing strengths worth preserving

- The repository generally keeps HTTP, application, and SQLC persistence concerns in separate packages.
- Contexts are present on most service and query APIs.
- SQLC-generated persistence code reduces ad hoc scanning and mapping errors.
- OpenTelemetry instrumentation already covers inbound HTTP, pgx, Redis, and selected application work.
- Background analytics processing updates its cursor in the same DB transaction as metrics, which is a good atomicity property.
- Focused tests exist for nontrivial moderation, search, analytics, validation, CSRF, and generator rules.
