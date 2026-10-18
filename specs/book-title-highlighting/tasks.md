# Tasks

- [x] AC1: Request/map highlights and version cached results.
- [x] AC2: Render marked title segments in search cards.
- [x] AC3: Test safe parsing, fallback, and rendering; generate templates and run checks.

Validation: affected Go packages pass, including request/mapping, parser and escaped template rendering tests. Templates regenerated; diff whitespace checks pass. Live OpenSearch and browser checks were not run. The installed templ CLI reports an older version than go.mod.
