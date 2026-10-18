# Tasks

- [x] Wire additional index filters through all search layers and cache.
- [x] Add page-size and filter controls, translations and URL serialization.
- [x] Verify focused regressions and frontend build.

Verification: Go tests for elasticstore, app and public handlers; node search-filter serialization regression; production frontend build; live HTTP combined filters and page-size clamping. Browser inspection could not run because the shared browser profile was in use.
