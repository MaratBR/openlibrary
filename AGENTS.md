# OpenLibrary contributor guide

This file applies to the whole repository. Add a more specific `AGENTS.md` in a subtree only when that area needs genuinely different rules.

## Project shape

OpenLibrary is a Go web application with server-rendered templ pages and a TypeScript frontend. The main technologies are:

- Go 1.24, Chi, and Fx for the HTTP server and dependency injection.
- PostgreSQL through pgx and SQLC.
- templ for server-rendered HTML.
- Alpine.js for behavior attached directly to server-rendered markup.
- React for richer client-side components and islands.
- SCSS component classes plus UnoCSS utilities for styling.
- Vite and pnpm for frontend builds.

Important directories:

- `cmd/server`: application entry point, configuration, root middleware, and server startup.
- `internal/app`: domain DTOs, service interfaces, business logic, and application errors.
- `internal/store`: SQLC query sources, migrations, and generated persistence code.
- `web/public`: public Chi controllers, HTML handlers, API handlers, and route registration.
- `web/public/templates`: public `.templ` source files and template helpers.
- `web/admin`: admin controllers and templates.
- `web/frontend/src`: TypeScript, Alpine components, React islands/components, API clients, and styles.
- `translations`: TOML translation catalogs.

## Architecture and change placement

Keep HTTP, business, and persistence concerns separated:

1. Controllers in `web/public` parse HTTP input, obtain authentication context, call application services, and render templates or API responses.
2. Interfaces, DTOs, state transitions, validation, and business rules belong in `internal/app`.
3. Database operations originate in `internal/store/query.*.sql`; application services use the generated `store.Queries` methods.
4. Register new public controllers in `web/public/handler.go` and their constructors in `web/public.FXModule`. Register application services in `internal/app/fx.go`.

Use the existing response helpers in `web/public/handler_response_api.go` and `handler_response_html.go`. Use helpers from `internal/olhttp` for URL/path/query parsing. Do not duplicate authentication or error-response conventions in individual handlers.

Use `app.Nullable[T]` for optional domain/JSON values where surrounding code already uses it. Its JSON representation is the contained value or `null`, not an object with `Value` and `Valid` fields.

## Templates and frontend behavior

- Frontend code under `web/frontend` follows a feature-oriented structure with two main variants: `features` and `islands`.
- `web/frontend/src/features` divides code by the application's functional features.
- `web/frontend/src/islands` contains UI elements rendered as islands in the application. Ideally, each island represents a feature, while not every feature is an island; in practice, the two structures currently overlap.
- This organization is a recent convention and is not yet consistently enforced across the existing frontend. Follow it for new work and when naturally touching related code, without broad unrelated reorganizations.
- Put APIs shared by multiple features in `web/frontend/src/features/api`. Put an API used by only one feature in that feature's `api` subfolder.
- Edit `.templ` sources, never ignored `*_templ.go` output.
- After changing a template, run `templ generate` and compile/test the affected Go package.
- Keep simple, page-local interactivity in Alpine `x-data` near the markup. Reuse/register an Alpine component under `web/frontend/src/alpinejs` when behavior is substantial or shared.
- Use a React island for stateful, reusable, or complex client UI. Existing islands live under `web/frontend/src/islands`; reusable React controls live under `web/frontend/src/components`.
- `web/frontend/src/public.api` is deprecated compatibility code for browser globals used by server-rendered templates. Do not add new APIs there; move shared APIs to `features/api`, move highly specific APIs to their owning feature's `api` subfolder, and migrate existing globals when touching their callers. Keep Zod response schemas aligned with the complete Go DTO fields that callers need; Zod object parsing strips undeclared fields.
- Frontend imports may use the `@/` alias for `web/frontend/src`.
- Server-provided Alpine data is commonly serialized to JSON in a template helper. Do not construct JSON by string concatenation.
- Use existing i18n helpers such as `i18nExtractKeys`/`i18nExtractKeysByPrefix`. Add user-facing copy to `translations/en.toml`; do not hard-code visible English when the surrounding component is translated.
- Preserve no-JavaScript server-rendered content where practical. Dynamic text should have a meaningful server-rendered fallback.

## Styling

The styling system intentionally combines shared SCSS components and UnoCSS:

- Look in `web/frontend/src/common/style/components` before creating new controls. Prefer established classes such as `.btn`, `.btn--primary`, `.btn--outline`, `.btn--ghost`, `.BtnGroup`, `.card`, and `.card--elevated` for consistent visual behavior.
- Use UnoCSS utilities for local layout, spacing, responsive behavior, positioning, and small one-off presentation details.
- Add a component-scoped SCSS class when a coherent visual pattern is reused or utility-heavy markup becomes difficult to read. Avoid adding a generic global class for a single local element.
- Reuse theme tokens (`primary`, `secondary`, `surface`, `foreground`, `border`, etc.) instead of literal colors. Use `secondary-foreground` for subdued text, `secondary` for subdued surfaces, and `foreground/5` for neutral interaction feedback. Their CSS variables and Uno mappings are defined in `theme.css` and `uno.config.ts`.
- Shared component SCSS is imported from `web/frontend/src/common/style/components/index.scss`.
- Uno scans templ, TS/TSX, CSS, and SCSS files, but dynamically constructed utility class names may not be detected. Use complete static class names or add a justified safelist entry.

## Database and generated code

- Edit `internal/store/query.*.sql`, then run `make db_sqlc` (or `sqlc -f internal/store/sqlc.yaml generate`).
- Put SQLC queries used exclusively by CLI commands, and not by the application layer, in a command-scoped `query.cli_<scope>.sql` file. Prefix their query names with `CLI_` followed by the same command scope; for example, use `query.cli_util_minifycontent.sql` and `CLI_Util_MinifyContent_Something`.
- Do not hand-edit files marked `Code generated by sqlc`, including `internal/store/models.go`, `query.*.sql.go`, `db.go`, and `copyfrom.go`.
- Database schema changes require a migration in `internal/store/migrations`. Create one with `make migration N=<name>`.
- `cmd/server/olproto/search.pb.go` is generated from `search.proto`; edit the proto and regenerate rather than patching generated protobuf code.
- `make codegen` runs both SQLC and templ generation. It may rewrite many generated templ files, but those outputs are ignored and should not be committed.

## Local development

Prerequisites and libvips setup are documented in `README.md`. Common commands:

```sh
docker compose up -d
make migrate_db
make ui_watch
make main_watch
```

The default server is configured by `openlibrary.toml` and listens on port 8080. Vite development proxy settings are in the same file. Keep secrets and machine-specific overrides in ignored `openlibrary.private.toml`, not the committed default config.

## Verification

Choose checks proportional to the change, and report unrelated baseline failures rather than fixing them without scope.

```sh
# All Go tests
go test ./...

# Focused public/template tests
go test ./web/public/...

# Regenerate templates
templ generate

# Type-check and build frontend assets
pnpm run build

# Check whitespace errors
git diff --check
```

If the execution environment has a read-only default Go cache, point `GOCACHE` and `GOTMPDIR` at a writable temporary directory rather than escalating permissions.

There are currently only a small number of Go unit tests. Add focused tests for nontrivial domain rules and regressions when feasible. HTTP or template-only changes should at minimum compile the affected packages after template generation.

## Working-tree discipline

- The worktree may contain user changes. Inspect `git status` before editing and preserve unrelated modifications.
- Do not revert, reformat, or “clean up” unrelated files.
- Avoid committing build output from `dist`, `build`, or generated templ files; these paths are ignored.
- Keep changes focused. Follow existing naming and file organization near the feature rather than introducing a parallel convention.
