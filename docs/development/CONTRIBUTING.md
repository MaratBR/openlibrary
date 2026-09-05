# Contributor workflow

OpenLibrary separates HTTP delivery, application rules, persistence, and
frontend behavior. The scoped `AGENTS.md` files provide the authoritative,
directory-specific guidance for changes in each area.

## Change placement

| Change | Home |
| --- | --- |
| Routes, request parsing, authentication context, rendered/API responses | `web/public` |
| Validation, policies, state transitions, and service interfaces | `internal/app` |
| SQL queries, migrations, and generated database access | `internal/store` |
| Server startup, root middleware, and runtime configuration | `cmd/server` |
| Browser behavior, islands, frontend API clients, and styling | `web/frontend` |
| User-visible text | `translations` |

Controllers should stay thin: parse input, acquire request context, invoke an
application service, and use the shared response helpers. Application services
own domain rules and use SQLC-generated `store.Queries` for persistence.

## Generated sources

Never edit generated sources directly:

- Change `.templ` files, then run `templ generate`.
- Change `internal/store/query.*.sql`, then run `make db_sqlc`.
- Change `cmd/server/olproto/search.proto`, then regenerate its protobuf output.
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
temporary paths for the command. Preserve unrelated changes and report baseline
failures rather than fixing them incidentally.
