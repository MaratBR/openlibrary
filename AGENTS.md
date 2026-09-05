# OpenLibrary contributor guide

This file applies to the whole repository. Read the nearest nested `AGENTS.md`
before changing files in a scoped area; it adds to this guide.

## Project shape

OpenLibrary is a Go web application with server-rendered templ pages and a
TypeScript frontend. Its primary technologies are Go 1.24, Chi, Fx, PostgreSQL
with pgx and SQLC, templ, Alpine.js, React, SCSS with UnoCSS, Vite, and pnpm.

| Area | Responsibility |
| --- | --- |
| `cmd/server` | Application entry point, configuration, root middleware, and startup |
| `internal/app` | Domain DTOs, service interfaces, business rules, and application errors |
| `internal/store` | SQLC query sources, migrations, and generated persistence code |
| `web/public` | Public controllers, handlers, routes, and templates |
| `web/admin` | Admin controllers and templates |
| `web/frontend` | TypeScript, Alpine, React islands, APIs, and styles |
| `translations` | TOML translation catalogs |

## Documentation and local development

- Project development documentation belongs in `docs/development`. Keep it
  current when a workflow or convention changes.
- Prerequisites and libvips setup are documented in `README.md`. Typical local
  startup uses `docker compose up -d`, `make migrate_db`, `make ui_watch`, and
  `make main_watch`.
- The default server configuration and Vite proxy settings are in
  `openlibrary.toml`; keep secrets and machine-specific values in ignored
  `openlibrary.private.toml`.

## Verification

Choose checks proportional to the change and report unrelated baseline failures
rather than fixing them without scope.

```sh
go test ./...             # all Go tests
go test ./web/public/...  # focused public/template tests
templ generate            # regenerate templates after .templ changes
pnpm run build            # type-check and build frontend assets
git diff --check          # whitespace errors
```

If the default Go cache is read-only, point `GOCACHE` and `GOTMPDIR` to writable
temporary directories rather than escalating permissions. Add focused tests for
nontrivial domain rules and regressions when feasible.

## Working-tree discipline

- Inspect `git status` before editing and preserve unrelated user changes.
- Do not revert, reformat, or clean up unrelated files.
- Do not commit ignored build output from `dist`, `build`, or generated templ
  files.
- Keep changes focused and follow nearby naming and organization conventions.
