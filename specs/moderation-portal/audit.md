# Moderation portal audit

Date: 2026-10-14. Source baseline: `9f4b980`; working tree was clean.
Scope: existing `/moderation`, its APIs, services, queries, and UI. Read-only
review; no enforcement actions, report decisions, or database cleanup performed.

## Evidence and limits

Inspected authenticated local pages using the existing admin session. Visited
Overview, Search, Chapters, Comments, Reports Overview/Search, Users, Books,
Login history, Audit log, report `40082`, book `6708568723978426170`, its chapter
list, and chapter `6708568771223066434`. Checked Users, Login history, and report
detail at 1440 × 1000 and 390 × 844; these had no document horizontal overflow.
Inspected desktop report and mobile report-search screenshots. Sample data had
502 users, 58,421 default-visible books, and 30,070 reports; these are observations,
not production counts or performance benchmarks.

Source review covered frontend routing/actions, report decision execution and
persistence, content moderation transactions, report queries, user report reads,
and the IP-location implementation. Runtime API probes confirmed missing chapter
and comment list APIs, ignored report status filtering, and invalid role filters
returning an empty successful response. Ordinary-user/anonymous permission tests,
destructive actions, full keyboard interaction, contrast measurement, dark-theme
rendering, database failure injection, and load testing remain unverified.

Focused checks passed:

- `go test ./internal/app -run 'Test.*(Moderat|Report|Login)' -count=1`
- `go test ./web/public -run 'TestModeration' -count=1`

These mostly verify service rules with stubs and filter parsing; they do not
establish database transaction safety. No build was needed for documentation-only
work. Earlier handoff claims about broader checks were not revalidated here.

## Workflow inventory

| Workflow | Current evidence | Remaining gap |
| --- | --- | --- |
| Start a shift | Overview is a heading-only placeholder | Real backlog, queue links, oldest pending work |
| Find an entity | Users and Books load searchable paginated directories | Global Search is empty; no chapter/comment directories |
| Review a user | Source has overview, content, history, reports, actions | Report reads create demo data; linked comments are placeholders |
| Review a book | Overview and chapter list loaded; source has actions/activity | Chapter links return 404; direct actions omit shadow/permanent controls |
| Review a chapter | Hide/restore service and POST API exist | No standalone read/detail route, list, content review, or history UI |
| Review a comment | Remove/restore service and POST API exist; reports load context | Directory/detail pages missing; thread context unavailable in portal |
| Triage reports | Search and ticket detail load real rows | Reports Overview empty; no status/priority/date queues or owner workflow |
| Decide a report | Persisted dispositions and target action dispatch exist | Concurrency and partial-failure gap; action options ignore current state |
| Inspect evidence | Book excerpt stored; other context comes from current rows | UI incorrectly labels all context as a submission-time snapshot |
| Review account sessions | Global history and filters load sessions | Displayed locations are fabricated; several controls lack accessible names |
| Follow accountability | Global audit and target history exist | Report events show JSON and cannot link directly to their ticket |

## Findings

Severity: P0 blocks trustworthy operation; P1 is major; P2 has a workaround.
Ten prioritized findings: one P0, six P1, three P2. Source-inferred risks are
identified separately from live reproductions.

### F-01 — P0: reading reports generates accusations

`internal/app/moderation_user.go:GetReports` creates four random reports whenever
the target has fewer than 50; `moderation_user_repository.go` persists them. A
unit test explicitly expects this behavior. An investigation can manufacture
additional complaints against its subject. Verified in source, intentionally not
triggered in the browser. Remove seeding from all production request paths; keep
fixtures opt-in and identify historical demo rows without guessing from wording.
Do not delete existing records automatically. Category: implementation integrity.

### F-02 — P1: core navigation leads to empty or missing workflows

`features/moderation/components/Portal.tsx` maps Overview, Search, Chapters, and
Comments to `PlaceholderPage`; `ReportsPage.tsx:ReportsOverview` returns null.
Comment detail is also a placeholder; chapter detail is absent. Live chapter
links returned `404`; chapter/comment list APIs returned 404. Moderators cannot
independently investigate these targets. Implement directories, contextual detail
and action/history views. Category: implementation integrity.

### F-03 — P1: no actionable report queues

Report search supports only text and target type, omits status/priority from list
DTOs and rows, and sorts newest first. The live `status=unreviewed` probe still
returned resolved report `40082`. Moderators cannot separate unresolved cases or
reliably work oldest/highest-priority first. Add validated queue filters, explicit
ordering, visible state, and counts. Category: implementation integrity.

### F-04 — P1: enforcement can precede a rejected case decision

`moderation_report.go:DecideReport` reads state, executes enforcement, then calls
`reports_repository.go:AddDecision`. That last operation atomically changes
status and records the decision, but enforcement has already committed separately.
Two reviewers can both enforce before one loses the final status update; a failed
case write can leave enforcement committed without a matching case event. Source
risk, not race-reproduced. Add a shared transaction/locking or durable operation
protocol covering the whole decision, with repeat-submission protection and tests.
Category: implementation integrity.

### F-05 — P1: location evidence is invented

`internal/app/ip_location.go` selects a random sample city for every IP. The live
session list displays international cities for loopback addresses without a demo
label. This can mislead abuse investigations. Represent provider/provenance and
unknown location honestly; do not trust existing unlabeled coordinates. Choosing
and paying for a geolocation provider is a separate product decision. Category:
implementation integrity.

### F-06 — P1: evidence/history labels overstate available facts

Report context queries read current user, comment, and book records. Only the
selected book excerpt is stored with the report; the UI labels the whole block
“Snapshot at time of report.” `ReportPage.tsx:BookContext` always renders a TODO
edit banner, a no-op Compare changes button, and “No previous enforcement” without
loading that history. Live report `40082` showed these plus the unresolved `{1}`
reporter-count interpolation. Distinguish saved evidence from current state;
show real history and unavailable revisions explicitly. Category: implementation
integrity; no claim that historical snapshots can be reconstructed.

### F-07 — P1: visible field labels are not accessible names

`LoginHistoryPage.tsx` uses `FormControl` without `htmlFor`; its fallback emits a
span. Live search/date inputs had neither an associated label nor an ARIA name.
Action forms use the same pattern. `ReportPage` also nests a main landmark inside
the shell main, confirmed in the DOM. Fix named controls and landmark structure;
verify keyboard/dialog focus and error announcements. Category: accessibility;
WCAG 1.3.1 and 4.1.2 apply to the verified unnamed inputs.

### F-08 — P2: standalone actions and audit navigation are incomplete

Book Actions exposes rating, summary, and restriction toggling; services and
report actions also support shadow restriction and permanent removal. Backend
report actions are selected by target type/scope, not current state. Global audit
has only target-type filtering; `ModerationLog.tsx` lacks report target links and
a report-decision renderer. Provide eligible standalone controls and actor,
action, target-ID/date audit filters with readable event links. Category:
implementation integrity.

### F-09 — P2: request and list behavior needs hardening

An invalid user role returns HTTP 200 with empty results; report status is silently
ignored. Book chapter lists fetch everything and paginate locally. Validate known
filters, use bounded server-side chapter/comment queries, retain URL state and
consistent empty/error/retry behavior. Category: performance/integrity. No latency
or bundle-size regression has been measured.

### F-10 — P2: small-screen navigation delays the task

At 390px, the wrapping sidebar with recent-record links occupies roughly 570px
before report-search content begins. No horizontal overflow was observed, but
repeated navigation dominates the viewport. Use a compact accessible mobile
navigation while preserving desktop rail and current styling. Runtime console
also reported missing `#client-flashes` initialization and favicon 404; treat
these as baseline integration defects, not proof actions fail. Category: responsive
and implementation integrity.

## UI audit assessment

Implementation integrity: fail for a complete moderator workstation, due to
fabricated data and incomplete investigation/triage. Existing visual language is
coherent and should be preserved.

| Dimension | Provisional score / 4 | Evidence and limits |
| --- | --- | --- |
| Accessibility | 2 | Unnamed controls and nested main; keyboard/contrast unmeasured |
| Performance | 2 | Bounded primary directories; unbounded chapter fetch; no profiling |
| Responsive | 2 | Sample pages avoid overflow; mobile navigation dominates viewport |
| Theming | 3 | Shared semantic palette with literal report-status colors; dark mode unverified |
| Implementation integrity | 1 | Real actions/audit foundations, but false evidence and missing workflows |
| Total | 10 / 20 | Provisional; significant work required, not a WCAG or performance certification |

The bundled detector found one advisory `side-tab` border on the quoted comment
in `ReportPage.tsx`. This is intentional quotation structure and matches incumbent
evidence styling; it is not a functional defect or a requested redesign.

Preserve: shared DashboardShell, URL-based filters/pagination, service-level
moderator authorization, reason validation, confirmation dialogs with pending
states, persisted decision activity, and transactional target mutation/audit logs.

Recommended UI follow-ups after domain fixes: Impeccable harden for field names,
evidence copy and errors; adapt for mobile navigation; polish as a final bounded
pass. Domain and transaction changes follow the implementation plan, not a visual
command. Re-audit after implementation using the same routes and widths.
