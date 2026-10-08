# Contributor workflow

OpenLibrary separates HTTP delivery, application rules, persistence, and
frontend behavior. The scoped `AGENTS.md` files provide the authoritative,
directory-specific guidance for changes in each area.

The application uses Go 1.24, Chi, Fx, PostgreSQL with pgx and SQLC,
server-rendered templ pages, and a TypeScript frontend with Alpine.js, React,
SCSS, UnoCSS, Vite, and pnpm.

## Change placement

| Change | Home |
| --- | --- |
| Routes, request parsing, authentication context, rendered/API responses | `web/public` |
| Admin controllers and templates | `web/admin` |
| Validation, policies, state transitions, and service interfaces | `internal/app` |
| SQL queries, migrations, and generated database access | `internal/store` |
| Server startup, root middleware, and runtime configuration | `cmd/server` |
| Browser behavior, islands, frontend API clients, and styling | `web/frontend` |
| User-visible text | `translations` |

Controllers should stay thin: parse input, acquire request context, invoke an
application service, and use the shared response helpers. Application services
own domain rules and use SQLC-generated `store.Queries` for persistence.

For a new feature or substantial behavior change, start with the
[spec-driven development workflow](spec-driven-development.md). Small fixes
can go straight to implementation when the expected behavior is already clear.

## Local development

See [README.md](../../README.md) for prerequisites and libvips setup. Typical
local startup uses `docker compose up -d`, `make migrate_db`, `make ui_watch`,
and `make main_watch`.

Default server configuration and Vite proxy settings are in
[`openlibrary.toml`](../../openlibrary.toml). Keep secrets and machine-specific
values in ignored `openlibrary.private.toml`.

## Generated sources

Never edit generated sources directly:

- Change `.templ` files, then run `templ generate`.
- Change `internal/store/query.*.sql`, then run `make db_sqlc`.
- `make codegen` runs SQLC and templ generation together. Generated templ files
  are ignored and should not be committed.

## Local checks

Run the narrowest relevant check first, then broaden it when useful:

```sh
go test ./web/public/...
go test ./...
pnpm run build
git diff --check
```

If Go cannot write its default cache, set `GOCACHE` and `GOTMPDIR` to writable
temporary paths for the command rather than escalating permissions. Preserve
unrelated changes and report baseline failures rather than fixing them
incidentally.

Diagnostic endpoints, server freshness checks, and Playwright verification are
documented separately in [Debugging and visual checks](debugging.md).
