# Verification

Passed:

- `make db_sqlc`; generated code comes from `query.notifications.sql`.
- `go test ./internal/app/notifications ./internal/app/userdata ./web/public/... ./web/admin/... ./cmd/server ./internal/i18n` with a writable temporary Go cache.
- `TestReadStatePostgres` against local PostgreSQL using `NOTIFICATION_TEST_DATABASE_URL`: temporary tables inside a rolled-back transaction, without reading or modifying application data. Covered missing read state, recipient isolation, metadata round-trip, a newer notification remaining unread, stale acknowledgements and future-time clamping. Without the environment variable, this integration test skips.
- Focused unit/HTTP tests cover the formatting allowlist and malicious markup, category defaults and emoji overrides, freeform types, invalid icons, read API validation/errors, authenticated recipient scoping, private user-data key protection, admin field forwarding and CSRF.
- `pnpm run build` and `git diff --check`.
- An isolated browser fixture rendered the actual React components against mocked APIs. Checked numeric unread count, successful read acknowledgement, failed history preserving unread state, failed acknowledgement/retry, duplicate/live delivery, safe formatting, category/emoji icons, Escape/focus, and panel bounds/scrolling at 1440px and 390px in both themes. Captures use production Poppins and Font Awesome assets and include a numeric badge and notification-center footer.

Scope: browser checks used synthetic content and API fixtures, not a restarted authenticated application. PostgreSQL behavior was checked independently through the real service and generated queries. No persistent application data was changed.

The design detector reported only advisory type-ramp differences for compact badges, timestamps and urgent text. Existing semantic theme tokens remain the source of truth; unrelated design metadata drift is unchanged.

Final visual review disposition: **ship**, scoped to the shared notification menu/record fixture in the four desktop/mobile light/dark captures. No material visual fixes were identified. This does not claim live application integration verification.
