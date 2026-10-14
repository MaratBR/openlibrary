# Implementation tasks

Status: audit and draft complete; application implementation has not started.
Check a task only when its changes and verification are complete.

## Discovery and agreement

- [x] Inspect repository instructions, workflow, current routes/services/queries and initial git status.
- [x] Inspect authenticated portal and bounded desktop/mobile layouts; record findings and verification limits in audit.md.
- [x] Draft core scope, acceptance criteria, implementation plan and ordered tasks.
- [ ] Review spec.md open questions; record agreed finish line, actor/target policy and optional scope.
- [ ] Update spec/plan/tasks to reflect the agreed requirements before coding.

## Slice 1 — trustworthy data and access

- [ ] Remove report generation from reads; introduce explicit fixture tooling and a read-purity regression (AC-01).
- [ ] Replace random location generation/display with truthful unknown/provenance; investigate old data without deleting it (AC-01, AC-16).
- [ ] Verify and enforce agreed role/target permissions on reads/actions (AC-02).
- [ ] Verify temporary-ban expiry and existing-session enforcement with controlled fixtures (AC-16).

## Slice 2 — safe report decisions

- [ ] Implement transaction-aware enforcement, report serialization and operation deduplication (AC-11–13).
- [ ] Add database-backed competing-decision, retry and partial-failure tests (AC-13).
- [ ] Verify disposition rules, reasons/payloads and immutable resolved cases (AC-11–12).

## Slice 3 — triage

- [ ] Add validated report status/priority/date/sort filters and matched list/count queries (AC-04).
- [ ] Add queue DTOs/API clients and state-visible report rows; verify pagination and ordering (AC-04).
- [ ] Implement Reports Overview queues and portal Overview with real counts/oldest pending links (AC-03).

## Slice 4 — chapter/comment investigation

- [ ] Add authorized, server-paginated chapter list/detail/content/history reads with integration fixtures (AC-02, AC-08).
- [ ] Implement chapter directory/detail routes and audited hide/restore controls; verify reader visibility (AC-07–08, AC-11).
- [ ] Add authorized comment list/detail/thread/history reads with removed/missing-context fixtures (AC-02, AC-09).
- [ ] Implement comment directory/detail and remove/restore controls; verify reader effects (AC-09, AC-11).
- [ ] Replace broken links from user/book/report/audit surfaces (AC-06–10, AC-15).

## Slice 5 — search and action coverage

- [ ] Implement bounded grouped global search with exact ID/UUID and text behavior (AC-05).
- [ ] Complete book standalone actions and server-derived eligibility; verify direct/report consistency (AC-07, AC-11).
- [ ] Complete user investigation/action behavior against agreed staff policy; retain existing filters (AC-06, AC-11).

## Slice 6 — evidence and collaboration

- [ ] Show saved vs current evidence, real related reports/enforcement, and honest unavailable revisions (AC-10).
- [ ] Implement append-only internal notes and reasoned priority changes with unresolved-state validation (AC-14).
- [ ] Verify escalation/resolution activity and requested/unavailable notification copy (AC-12, AC-14).

## Slice 7 — accountability

- [ ] Implement audit actor/action/target/date search and readable ticket-linked event rendering (AC-15).
- [ ] Harden session filters and preserve canonical user-to-global history navigation (AC-16).

## Slice 8 — release quality

- [ ] Fix field names, landmarks, dialog focus and feedback; implement practical mobile navigation (AC-17).
- [ ] Complete translations, missing/empty/error/retry states, filter restoration and DTO generation (AC-17–18).
- [ ] Run isolated end-to-end moderator workflows, desktop/mobile light/dark review, and regression checks (AC-01–18).
- [ ] Record acceptance evidence and proportional checks, migrations/rollout notes, baseline failures and remaining limitations.
- [ ] Confirm every agreed acceptance criterion is verified; leave deferred features explicitly outside completion.
