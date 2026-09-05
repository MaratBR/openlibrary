# Persistence and SQLC

This guide applies to `internal/store`. Follow `internal/AGENTS.md` and the
repository guide as well.

- Edit `query.*.sql`, then run `make db_sqlc` (or
  `sqlc -f internal/store/sqlc.yaml generate`).
- Do not hand-edit SQLC output, including `models.go`, `query.*.sql.go`, `db.go`,
  and `copyfrom.go`.
- Put queries used only by CLI commands in `query.cli_<scope>.sql`; prefix names
  with `CLI_<Scope>_`, such as `CLI_Util_MinifyContent_Something`.
- Add schema changes through a migration in `migrations`, created with
  `make migration N=<name>`.
