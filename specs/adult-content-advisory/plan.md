# Plan: Adult-content and warning preferences

Status: provisional; finalize against agreed spec before implementation.

1. Resolve the spec's open questions and inventory every existing public book
   retrieval path, cache, and equivalent API. Read scoped AGENTS.md files before
   editing internal, persistence, public, frontend, or translation code.
2. Define resolved preferences and a central policy in `internal/app`, separating
   authored rating, effective adults-only classification, and matching descriptors.
   Return reasons and surface-specific outcomes. Keep availability authorization
   ahead of preference exceptions. Cover the agreed boundary with focused tests.
3. Resolve the editor contract in `BookEdit.tsx`, public manager APIs, and manager
   DTOs. Remove the independent switch; implement agreed legacy-field handling.
   Retain existing rating storage unless an approved decision requires migration.
4. Design persistence and transport for canonical descriptor IDs. Use decimal
   strings in browser JSON. Prepare reversible migrations and explicit mappings
   for existing string preferences, censor modes, and adult flags; surface
   unresolved entries. Add settings validation and translated selector controls.
5. Derive SQL and OpenSearch constraints from the same resolved policy. Apply them
   before pagination/counting, retaining manual exclusions and existing visibility
   checks. Replace unconditional random R/NC-17 exclusion. Account for cached
   parent tag IDs and synonym resolution. Identify index fields/backfill needed
   for effective classification; invalidate on tag changes and preference updates.
6. Implement server-rendered direct-link interstitials for book, chapter, preview,
   and equivalent APIs before sensitive payload loading/rendering. Use server-
   validated per-book exceptions, bound to the agreed session/account lifetime.
   Remove legacy cookie authorization; provide translated reasons, spoiler reveal,
   safe back navigation, and accessible form submission with CSRF protection.
7. Apply the agreed placeholder policy to saved/history surfaces without deleting
   membership. Verify cover/summary payloads, chapter fragments, and cache keys.
   Keep author/moderator management flows under their existing authorization.
8. Stage rollout with migrations, search backfill, and cache invalidation. Enable
   enforcement only with a consistent index or a defined safe fallback; stale
   documents must not reveal blocked books. Document recovery and rollback,
   including schema/API compatibility and conservative preference mapping.
9. Verify the acceptance criteria with policy, transport, migration, backend,
   HTTP, and browser checks. Run appropriate Go tests, templ generation after
   template edits, frontend build after frontend edits, and git diff --check.
   Record baseline failures and remaining limitations in verification.md.

## Verification mapping

| Criteria | Evidence |
| --- | --- |
| AC-1–2 | Editor/settings transport regressions and browser save/reload/error checks |
| AC-3 | Rating/descriptor truth table and conflict tests |
| AC-4 | SQL/OpenSearch result/count parity; random empty state; surface inventory |
| AC-5 | Book/chapter/preview HTML and API responses contain no blocked payload |
| AC-6 | Exception scope, expiry, concurrent tabs, and account transition tests |
| AC-7 | Saved-item placeholder and preserved membership/history checks |
| AC-8 | Anonymous/account precedence, failure handling, and cache isolation tests |
| AC-9 | Migration fixtures for every legacy mode and cookie bypass regression |
| AC-10 | Rename/merge/delete fixtures, int64 transport, index/cache invalidation |
| AC-11 | Availability and author/moderator authorization regressions |
| AC-12 | Keyboard, mobile, theme, and JavaScript-disabled interstitial checks |
