# Implementation plan

Status: proposed; depends on review of [spec.md](spec.md). Findings and baseline
verification are in [audit.md](audit.md). Execute complete slices in this order.

## Affected areas

| Area | Ownership |
| --- | --- |
| Portal routes, loaders, views, API clients and styles | `web/frontend/src/features/moderation/` |
| Existing React island entry/registration | Follow actual registration; old handoff island paths are stale |
| Shared field/dialog/navigation fixes | `web/frontend/src/components/`, narrowly scoped and checked across consumers |
| HTTP DTOs, input parsing and API responses | `web/public/routes_moderation_api.go` |
| Moderation rules, services, report decisions and repositories | `internal/app/moderation*.go`, `reports*.go`, authentication/ban owners as needed |
| Query source and schema | `internal/store/query.*.sql`, ordered migrations |
| Portal host template and permissions | `web/public/routes_moderation_portal.go`, `templates/moderation_portal.templ` |
| Localized copy and API definitions | `translations/en.toml`, generated frontend DTO definitions |

Read nearest AGENTS.md before edits. Keep controllers thin and action services
responsible for policy. Do not broadly move existing code or replace dashboard
visual identity. Follow existing Effect/API and shared control conventions.

## Slices and design decisions

### 1. Trustworthy reads and baseline protections (AC-01, AC-02, AC-16)

Remove request-time random-report generation; replace the existing test that
expects it with a read-purity regression. Demo fixtures live in explicit tooling.
Investigate provenance of historical demo rows; do not delete or infer origins
from description strings. Replace random geolocation with truthful unavailable
state and plan provenance for previously persisted locations. Verify account ban
state, expiry and existing-session enforcement before claiming those actions work.
Agree staff/self-target policy and test service authorization before adding reads.

### 2. Report decision consistency (AC-11–AC-13)

Prefer a shared database transaction that locks the report, revalidates state,
executes target mutation plus target audit through transaction-aware services,
and appends the case decision/status before commit. Reuse rules; do not duplicate
user/book/content mutation logic in the report repository. Locking only the final
status update is insufficient. Establish consistent report/target lock order.
Add a persisted operation key/uniqueness rule for retry deduplication; define the
API behavior for reused keys with different payloads. If existing architecture
prevents one transaction, document a durable operation/reconciliation design and
meet AC-13 before shipping. No silent two-commit failure window.

### 3. Report queues and work overview (AC-03, AC-04)

Extend report query/DTO/API with status, priority, dates and a validated sort enum;
use one shared predicate for counts and lists, parameterized SQL and stable ties.
Default pending view includes unreviewed and escalated. Derive summary/oldest age
from the same authoritative records. Add status/priority badges and URL filters;
Reports Overview becomes the queue surface; portal Overview summarizes that work.
Assess query plans with realistic fixtures before choosing indexes. Do not add
cosmetic metrics without a concrete moderator action.

### 4. Chapter/comment investigation (AC-07–AC-10)

Add typed list/detail/read-history endpoints and authorized services/repository
queries. Chapter reads must include restricted content without making it public;
comment reads must include removed content and bounded thread context. Reuse the
existing hide/restore/remove actions and audit format, verifying that reader
visibility actually changes. Paginate server-side; do not retrieve all chapters
or full comment threads for a directory page. Use string IDs end-to-end to retain
int64 precision. Ensure removed/missing parents produce meaningful context states.

### 5. Search and complete target actions (AC-05–AC-11)

Compose bounded grouped global search from feature services. Preserve existing
exact-ID, partial text and role/status semantics. Complete user/book detail links
and expose state-appropriate actions including shadow restriction/permanent book
removal. Derive supported action capabilities on the server and validate them
again on submit. Show before/after values and explicit irreversible consequences.
Resolve role hierarchy before enabling privileged-account mutations.

### 6. Evidence and lightweight case collaboration (AC-10, AC-12, AC-14)

Distinguish live context from persisted excerpts. Load real prior enforcement and
related reports; remove the fake edit warning, comparison and reporter-count copy.
Provide explicit unavailable historical evidence; do not backfill guessed snapshots.
Extend existing append-only moderation logs with versioned report_note and
report_priority_changed events, preserving report_decision compatibility. Note and
priority writes revalidate unresolved state transactionally. No new parallel event
store. Retain requested/unavailable notification metadata. If review adds real
communication or full snapshots, update this plan before implementing them.

### 7. Audit and account investigation (AC-15, AC-16)

Extend audit filters and indexes only as query evidence requires. Add report target
links and readable known-event renderers; preserve unknown JSON fallback. Ensure
chapter-related events retain meaningful chapter and book context. Session history
stays in its canonical existing route, with truthful location state and validated
filters. No provider subscription, session termination feature or export is implied.

### 8. Operational UI and release verification (AC-17, AC-18)

Fix label associations and page landmarks in touched workflows; reuse dialogs
with focus return, confirmation summaries, inline errors and pending feedback.
Make mobile navigation compact while preserving desktop rail and both themes.
Verify state restoration, safe rendered content, long names/reasons and missing
assets/data. Keep all visible copy translated and generate API definitions when
DTOs change. Bound the visual review to desktop/mobile in both themes, one fix
batch and one confirmation pass.

## Data changes and rollout

Prefer additive migrations for new log/operation metadata and required indexes.
Existing report statuses and version-1 decisions remain readable. New report
list fields require aligned Go DTOs/client schemas; existing URLs remain valid.
Generate SQLC after query edits, templ after template edits, and TypeScript DTOs
after response changes; do not edit generated sources by hand or commit ignored
build output. Never assume old locations are trustworthy or demo reports can be
identified without provenance. Run race/mutation scenarios against isolated test
data; do not enforce actions on the user's existing local records for verification.

## Verification matrix

| Criteria | Evidence required |
| --- | --- |
| AC-01 | Repeated list/detail reads preserve report counts; fixture tooling isolated; unknown-location/history rendering |
| AC-02 | Anonymous/member/moderator/admin read/action cases; agreed staff/self-target matrix; no unauthorized mutation/log |
| AC-03–04 | Query integration tests for mixed states/priorities, inclusive dates, deterministic sorting, counts/pagination; browser queue links |
| AC-05 | Exact UUID/int64 and text searches, grouped bounds, missing/invalid input; live links into every target type |
| AC-06–07 | User/book detail fixture cases with restricted/removed states; safe reads and functional report/content/history links |
| AC-08–09 | Database-backed directory/detail tests; restricted/removed authorized inspection; hide/restore/removal effects in reader views |
| AC-10 | Live vs saved evidence fixture assertions, missing parents, related-report/enforcement history, no fake comparison |
| AC-11–12 | Action payload/reason/state validation; disposition rules; confirmation cancel; successful UI refresh and failed input retention |
| AC-13 | Database-backed parallel decisions and retry tests; failure injection before/after target mutation and case event; consistent audit/state |
| AC-14 | Empty/oversized notes, priority reason/enum, immutable events, resolved-state conflict and internal-note confidentiality |
| AC-15 | Actor/action/target/date filter combinations; known/unknown payload rendering; each ticket/target link works |
| AC-16 | Date/enum validation; truthful location provenance; ban expiry and existing-session effects with controlled clock |
| AC-17–18 | Isolated-fixture browser scenarios at 1440/390px in both themes, keyboard/labels/focus/status, URL/back/pagination, long text and errors |

Start with focused Go service/HTTP tests and database integration tests for new
query/transaction behavior. Run `pnpm run build` for frontend changes. At release,
run `go test ./...`, frontend build and `git diff --check`; report unrelated baseline
failures separately. Use writable `/tmp` Go caches if needed. Existing focused
checks pass; old notes mention an unrelated chapter-name test failure, but it has
not been revalidated. Database integration harness/fixtures must be established
before marking concurrency and visibility criteria complete.
