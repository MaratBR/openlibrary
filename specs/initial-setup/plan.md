# Implementation plan

1. Add SiteSetup to the siteconfig snapshot and settings definitions.
2. Add a SQLC query reporting whether an administrator exists.
3. Implement and register SetupService, with mutex-protected cached checks and fresh state reads.
4. Create missing administrators through existing user creation helpers in a transaction.
5. Save setup configuration after admin creation; restore active config on save failure so retries are safe.
6. Remove middleware caching and verify rules, caching, errors, persistence, admin creation, and Fx wiring.

Admin creation and configuration persistence use separate transactions. A failed
configuration save leaves the created admin available for the next setup attempt.

7. Render a server-side setup form using the existing HTML layout and translations.
8. Handle GET/POST in setup middleware ahead of authorization, using fresh setup state, existing CSRF verification, typed validation errors, and a login redirect.
9. Test form visibility, submission, retained values, CSRF, completion, state/save failures, and method restrictions. Generate templ output and run affected Go tests and the frontend build.
