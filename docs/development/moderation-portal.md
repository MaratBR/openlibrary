# Moderation portal implementation notes

This is the cross-session handoff for ongoing work in `web/frontend/src/islands/moderation`.

## Current goal and state

The `/users` route is now a functional, paginated user directory. It supports case-insensitive partial-name search, exact UUID lookup, banned/active filtering, and role filtering. Results show avatar, name, UUID, role, joined date, last visit, and latest ban time/reason, and link to the existing `/users/:userId` moderation view.

The `/books/:bookId` route is now implemented with Overview, Actions, Chapters, and Activity tabs. Overview includes author/state/report context and the latest pending report. Actions support age-rating changes, restricting/restoring a book, and summary changes, all with required moderation reasons and audit entries. Chapters are fetched as one list and searched/paginated client-side at 100 rows per page; chapter links intentionally point to the not-yet-implemented `/chapters/:chapterId` route.

The `/audit-log` route lists all moderation actions with server-side pagination and target-type filtering. User history, book activity, and the global audit log share the same log-entry and payload renderer. Action-specific payload renderers are registered in `ModerationLog.tsx`; unknown payloads fall back to a formatted JSON or plain-text `<pre>`.

The `/reports/:reportId` route now has a persisted decision workflow. Outcomes are no violation, request changes, enforcement action, or escalation. Target-specific actions are supplied by the backend and dispatched through existing user/book/chapter/comment moderation services. Successful decisions update report state and append a versioned `report_decision` entry to `moderation_logs`; the ticket activity renders those persisted events. User and comment reports include real target snapshots rather than an ID-only placeholder.

The global `/login-history` route now uses server-side pagination and supports a reusable asynchronous multi-user selector, free-text filtering across IP, user agent, and location, inclusive date bounds, and active/terminated session filtering. Results identify the user and include IP location and session state. `/users/:userId/login-history` redirects to this page with that user preselected, keeping one canonical history UI.

The `/books` route is a paginated moderation book directory. Numeric search is an exact ID lookup; text search is partial by default and may be switched to an exact case-insensitive name match. Banned/shadow-restricted and trashed/permanently removed books are excluded by default and can be included independently.

The implementation is currently uncommitted. Relevant files are:

- `web/frontend/src/islands/moderation/ModerationUsers.tsx` and `Portal.tsx`
- `web/frontend/src/islands/moderation/BookModerationPage.tsx`
- `web/frontend/src/api/moderation/users.ts` and `index.ts`
- `web/public/routes_moderation_api.go`
- `internal/app/moderation_user.go`, `moderation_user_repository.go`, and tests
- `internal/store/query.moderation-user.sql` plus SQLC output
- `translations/en.toml` and generated `web/frontend/src/backend-types.ts`

## Data and API decisions

- Endpoint: `GET /_api/moderation/users` with `search`, `banned`, `role`, `page`, and `pageSize` query parameters.
- A valid UUID search is treated as an exact ID lookup; other search text uses `ILIKE` on the user name.
- “Last visit” is the newest `sessions.created_at` for the user.
- Ban metadata is the newest `user_bans` row (`created_at` and `note`). No migration or bogus data was necessary.
- Authorization remains in `ModerationUserService` through `ModerationAuthorizer`.
- Role filter options are sourced from `app.AllRoles`, serialized into the moderation island data, and passed into `ModerationUsers`; do not duplicate the role list in TypeScript.
- Nullable timestamps use `app.Nullable[time.Time]` in the HTTP response so generated TypeScript correctly emits `Nullable<string>`.
- Sessions persist placeholder IP geolocation (`country`, `region`, and `city`) through `IPLocationService`. The service intentionally returns random sample locations and contains the TODO for a real provider.
- The overview loads `GET /_api/moderation/users/:userId/login-locations`, which returns the three newest distinct non-empty session locations.
- Recent book names open `/book/:id` in a new tab and have a moderation link to the TODO `/books/:bookId` route. Recent comments link to the TODO `/comments/:commentId` route.
- Report rows link to `/reports/:reportId`. Report workflow fields and book-report scope are persisted on `reports`; the detail query loads actual book, author, chapter, warning tags, publication state, update time, and related-report count. `up.sql` upgrades an existing development database and `tmp.sql` creates whole-book, chapter, and selected-text samples.
- Reports are a sidebar section. `/reports/search` provides paginated search across report number, reason, description, and target ID, plus target-type filtering; report detail supports decisions and ticket activity.
- Report decisions reuse `moderation_logs` (`target_type = 'report'`, `type = 'report_decision'`) rather than introducing a parallel event table. The report status update and event append share a repository transaction.
- Moderation list filters use the shared Select and Radix/shadcn-style Checkbox controls; the portal contains no native checkbox controls. Text/date fields continue to use the established `.input` component class because the repository has no shared React input wrapper.
- Endpoint: `GET /_api/moderation/login-history` with optional comma-separated `users` (at most 20), `search`, `status`, `dateFrom`, `dateTo`, `page`, and `pageSize` parameters. An empty user filter searches all users.

## Verification baseline

- `pnpm run build` passes (with existing Browserslist and unresolved runtime-image warnings).
- Focused moderation/report/login-history tests pass, and `go test -vet=off ./web/public` passes.
- `go test -vet=off ./...` compiles all packages but fails the unrelated existing `internal/app` case `TestValidateChapterName/too_long` because `validateChapterName` currently accepts that fixture.
- `git diff --check` passes.
- Runtime checks against the local development stack passed with an existing moderator session: `GET /_api/moderation/books` returned paginated catalog data; login history combined active-state and user-agent filters and returned location/expiry fields; and a real user-report detail returned its target snapshot, activity, related-report count, and backend-derived actions. The unauthenticated `/moderation` route correctly redirects to `/login`.
- Run `make db_sqlc` after SQL changes and `go run ./cmd/go2tsdef` after generated API DTO changes.

## Likely next work

Visually inspect `/moderation#/users` against real seeded data, especially narrow/mobile layouts and long ban reasons. Consider API validation for unknown filter values and integration-level query tests if a test database fixture is introduced.

## Report decision workflow design

The report detail page is the case-management surface. A decision and an enforcement action are related but distinct:

- A **disposition** explains the case outcome: `no_violation`, `request_changes`, `action_taken`, or `escalated`.
- An optional **action** changes the reported target through an existing moderation service. Actions are target-specific and typed; the report repository must not implement user/book/comment mutations itself.
- A decision always records a policy reason and may include an internal note. Author notification is stored as requested delivery metadata until a real notification transport exists; the UI must say when delivery is not yet available.
- `no_violation`, `request_changes`, and `action_taken` resolve the report. `escalated` changes it to an escalated state so another moderator can later resolve it. Resolved reports are immutable in the MVP; reopening can be added as a new transition rather than rewriting history.

Supported action registry for the MVP:

| Target | Actions |
| --- | --- |
| User | no target mutation, temporary ban, permanent ban, unban, rename, change profile/about |
| Book (whole book) | no target mutation, restrict, restore, shadow restrict, restore shadow restriction, permanently remove, change age rating, change summary |
| Book chapter | no target mutation, hide chapter, restore chapter; whole-book actions remain available when appropriate |
| Comment | no target mutation, remove comment, restore comment |

`request_changes` and `escalated` do not require a target mutation. Warning/message delivery is not currently an application capability, so the MVP records the requested author notification in the ticket activity without falsely claiming that a message was sent.

### Persistence and activity

Use the existing append-only `moderation_logs` table for report events with `target_type = 'report'`, the report ID as `target_id`, and `type = 'report_decision'`. The policy reason uses the existing `reason` column; disposition, action, internal note, notification request, and action payload use a versionable JSON payload. The initial submitted event may continue to be synthesized from the report row; every later case transition comes from persisted log entries. These entries are the authoritative ticket activity stream and also appear in the global audit log.

Updating the report status and inserting its event must occur in one repository transaction. Target enforcement uses the existing transactional moderation services first; only a successful action may be followed by report resolution. A failed action leaves the report open and produces no decision event. This is not a distributed transaction, but it prevents a case from being marked resolved when enforcement fails and keeps target rules/audit logging in their existing services.

### API and validation

`POST /_api/moderation/reports/{reportID}/decisions` accepts:

```json
{
  "disposition": "action_taken",
  "action": "book.change_age_rating",
  "policyReason": "Incorrect age rating",
  "internalNote": "Reviewed the reported chapter and surrounding context.",
  "notifyTarget": true,
  "payload": { "value": "mature" }
}
```

The application layer owns disposition transitions, required fields, action/target compatibility, IDs and payload validation. The HTTP layer only parses input and maps responses. The detail response exposes the available actions derived from the report scope and current target state so the frontend does not duplicate backend capability rules.

### UI behavior and acceptance criteria

- Use the shared Select, Checkbox, FormControl, modal/confirmation, and button components/styles; do not introduce native selects.
- Choosing a disposition narrows the action choices. Choosing an action renders only its required value fields (for example ban expiry or age rating).
- Policy reason is required; internal note is optional and length-limited. Destructive actions require confirmation and all submissions expose pending/error states.
- After success, route data revalidates, the header status changes, and ticket activity shows actor, disposition, action, reason, note, notification request, and time.
- The backend rejects unknown dispositions/actions, incompatible targets, missing action payloads, and attempts to decide an already resolved report. Focused service tests cover these rules and verify that failed enforcement does not resolve a report.
