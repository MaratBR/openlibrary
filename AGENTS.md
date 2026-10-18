# OpenLibrary contributor guide

- Before editing, inspect `git status`; preserve unrelated changes.
- Read applicable nested `AGENTS.md` files along the path to changed files.
- Use specs only for end-to-end feature implementation or when the user
  explicitly requests spec-based development. In those cases, follow
  [Spec-driven development](docs/development/spec-driven-development.md).
  Other work does not require a spec.
- Keep changes focused and follow nearby conventions. Update development docs
  when workflows or conventions change.
- Never edit generated sources directly. Edit `.templ` or SQL query sources,
  then run `templ generate` or `make db_sqlc`, respectively. Do not commit
  ignored output from `dist`, `build`, or generated templ files.
- Run checks proportional to the change; report unrelated baseline failures.
  Add focused tests for nontrivial domain rules and regressions when feasible.
- Consult [Contributor workflow](docs/development/CONTRIBUTING.md) when needed
  for change placement, local setup, or verification commands.
- Keep communication brief: only necessary updates of 1–2 sentences; avoid
  long explanations and large code blocks.
