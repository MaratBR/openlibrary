# Feature: Complete moderation portal

Status: DRAFT FOR REVIEW. Requirements below are proposed, not yet agreed.
Audit: [audit.md](audit.md). Plan: [plan.md](plan.md). Progress: [tasks.md](tasks.md).

## Problem and outcome

Moderators need to discover pending work, inspect the reported target in context,
make a justified decision, and verify its effect and history. Existing services
and several portal screens support this, but incomplete investigation routes,
missing queues, demo data, and decision consistency gaps prevent dependable use.

Completion means every core acceptance criterion below has an implementation and
verification evidence. No advertised route remains a placeholder, no request
manufactures evidence, and an unsuccessful or competing decision cannot silently
leave contradictory enforcement and report state.

## Scope

### Core release — proposed finish line

- A real work overview, report queues/search, and entity search.
- User/book directories and contextual detail; chapter/comment directories,
  detail, review, standalone enforcement, related reports and audit history.
- Safe existing actions: user temporary/permanent ban, unban, rename and about;
  book restriction/restoration, shadow restriction/restoration, permanent
  removal, age rating and summary; chapter hide/restore; comment remove/restore.
- Persisted report decisions, escalation, independent internal notes, priority
  changes, and safe handling of concurrent/repeated submissions.
- Searchable audit history and trustworthy account/session investigation.
- Existing hash URLs, shared controls, themes, English localization, desktop
  rail, and responsive/keyboard operation.

### Deferred unless explicitly added

Ownership/claiming and reassignment; reopening resolved cases; bulk enforcement;
user-facing warnings/messages and notification delivery; appeals; formal SLA
alerts; automated abuse detection; report merging; exports; a new geolocation
provider; complete historical content snapshots and revision diffing.

The core must state when these capabilities/evidence are unavailable and remove
nonfunctional controls. Permanent removal retains existing semantics; it does
not introduce data purging. Role administration stays in the admin dashboard.
Historical demo-record cleanup requires provenance and a separately reviewed plan.

## Scenarios and acceptance criteria

### Trust and access

- AC-01: Loading or searching any portal/API view changes no domain records.
  Reports and metrics come from persisted submissions; demo fixtures are explicit
  development tooling. Unknown/unverified location or historical evidence is
  labeled unavailable, never filled with random values or invented history.
- AC-02: Anonymous callers cannot access moderation data/actions; ordinary users
  are denied; eligible moderators/admins can use core workflows. Every new read
  and mutation enforces permissions in the application layer. Sensitive account
  data remains restricted; admin links do not grant extra privileges.

### Find and investigate

- AC-03: Overview shows accurate unreviewed/escalated counts and oldest pending
  age, with links to the matching queues and recent decisions. Empty backlog
  and failed loads have explicit states; counts never include generated records.
- AC-04: Report queues support status (unreviewed, escalated, resolved), existing
  priority values, target type, submitted date range, text/number/target lookup,
  and explicit oldest/newest/priority order. Default is unresolved work ordered
  by priority then oldest submission, with a deterministic ID tie-breaker.
  Rows display status, priority, target, reporter, reason and time; pagination
  preserves filters. Counts use the same predicates as lists.
- AC-05: Global Search searches users, books, chapters, comments and reports with
  bounded, grouped results and links to usable detail views. Numeric IDs and UUIDs
  receive exact lookup where applicable; free text follows documented matching.
  Empty, malformed and missing results never send moderators to blank pages.
- AC-06: User detail exposes current profile/ban state, authored content, submitted
  reports against the user, enforcement history and session-history links. Existing
  searches/status/role filters continue to work without seeding reports.
- AC-07: Book detail shows current visibility/restriction/removal/rating state,
  author, summary, related reports, chapters and history. All chapter links open
  functional moderation views, including restricted content for authorized staff.
- AC-08: Chapter directory supports ID/name, book/author and visibility filtering
  with server pagination. Detail shows content, book/author, publication and
  moderation state, update time, related reports and audit history. Hide/restore
  updates reader visibility consistently without losing content.
- AC-09: Comment directory supports ID/text, author/book/chapter and removed-state
  filtering with server pagination. Detail shows the comment and available thread
  context, author, timestamps, parent content, reports and history. Removed
  comments remain inspectable by authorized staff; removal/restoration is audited.
- AC-10: Report detail distinguishes stored submitted evidence from live target
  state, links to usable target/context views, lists related reports and actual
  enforcement history, and clearly handles missing/deleted content. Revision
  comparisons appear only if supported by persisted evidence; otherwise the UI
  explains the limitation without a fake comparison button or snapshot claim.

### Act and collaborate

- AC-11: Direct and report-driven actions require a nonblank policy reason and
  valid target-specific fields (including future temporary-ban expiry). Available
  controls reflect current target state and actor capabilities, with authoritative
  server validation. Confirmation identifies target, change and consequence;
  cancel performs no mutation; pending submissions cannot be repeated locally.
  Errors retain entered values; success refreshes target state and history.
- AC-12: No violation, request changes and action taken resolve a report;
  escalation leaves it available in the escalated queue. Action taken requires a
  compatible enforcement action. Other dispositions perform no target mutation.
  Resolved reports stay read-only in core. Notification requests are recorded as
  requested/unavailable, never presented as delivered.
- AC-13: One successful decision records its actor/reason/time and changes report
  state together with the intended enforcement exactly once. Concurrent stale
  decisions and retries cannot apply additional actions or overwrite the winner.
  On any persistence/enforcement failure, either the whole operation rolls back
  or a durable, explicitly visible recoverable state identifies committed work;
  success is never reported for an incomplete operation.
- AC-14: Moderators can append internal notes and change priority on unreviewed
  or escalated reports without resolving them. Notes require nonblank text and
  have a 1,000-character limit. Priority changes require a reason. Actor, old/new
  values and time appear in immutable ticket/audit history; resolved reports
  reject these mutations. Reporter-facing content never exposes internal notes.

### Account history and operational quality

- AC-15: Audit search supports actor, action, target type/ID and date range, with
  deterministic pagination. Known events have readable summaries and valid
  target/ticket links; unknown payload versions retain a safe readable fallback.
  User/book/chapter/comment history agrees with the global audit log.
- AC-16: Login history retains user, search, inclusive dates and session-state
  filters. IP, browser, timestamps and location provenance are truthful; missing
  or unverifiable locations explicitly read unknown. Invalid filters and inverted
  dates return actionable validation errors. A temporary ban expires according
  to its recorded time; existing sessions respect active bans after enforcement.
- AC-17: Core routes support direct URLs, browser back/forward and filter-preserving
  pagination. Empty/loading/not-found/forbidden/error states are explicit; failed
  loads can be retried. Every field has an accessible name, keyboard controls and
  dialog focus work, one main landmark contains the page, and operation feedback
  is announced. At 1440px and 390px in both themes, navigation remains practical,
  text wraps and page content does not overflow; touch controls meet 44px targets.
- AC-18: Moderator-visible text uses existing English translations and semantic
  UI tokens; no TODO/demo text, broken interpolation, false success messages or
  no-op action links ship. Existing workflows remain intact. New nontrivial
  rules and race/failure regressions have focused tests, and end-to-end core
  scenarios pass against isolated fixtures.

## Open questions — resolve before behavior-changing implementation

1. Is this core release the agreed finish line? In particular, independent notes
   and priority editing are proposed additions; claiming/assignment, bulk actions
   and reopening are deferred. Which deferred capability is required for launch?
2. Should moderators be allowed to act on administrator/system/moderator accounts
   or themselves? Current role checks are broad. Agree the actor/target action
   matrix before changing permissions; recommended default is to reserve staff
   account enforcement for administrators and reject self-enforcement.
3. Confirm default queue ordering and existing priority labels. This draft proposes
   highest priority then oldest submission; no automatic SLA or priority scoring.
4. Is “request changes” acceptable as a recorded resolution with delivery marked
   unavailable, or must the release include real author communication? The latter
   expands scope to a notification channel, delivery state, and retry policy.
5. Existing report evidence cannot reconstruct past full content. Is honest live
   context plus saved excerpts sufficient for core, or are submission snapshots
   and actual revision comparisons required? Likewise, core proposes unknown
   location rather than choosing a paid geolocation provider.

## Review decision

Pending. Record agreed scope and answers here, then align plan/tasks. An elapsed
wait or an unchecked task is not approval. No application implementation is
included in this audit/spec change.
