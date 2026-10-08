# OpenLibrary contributor guide

This file applies to the whole repository. Read the nearest nested `AGENTS.md`
before changing files in a scoped area; it adds to this guide.

- Read [Contributor workflow](docs/development/CONTRIBUTING.md) for project
  structure, local setup, generated sources, and verification commands.
- For new features and substantial behavior changes, follow
  [Spec-driven development](docs/development/spec-driven-development.md): agree
  on observable requirements before implementation, then record the plan and
  track work in `specs/<feature>/`.
- Keep development documentation in `docs/development` current when workflows
  or conventions change.
- Inspect `git status` before editing. Preserve unrelated changes; do not
  revert, reformat, or clean up unrelated files.
- Keep changes focused and follow nearby conventions. Do not commit ignored
  output from `dist`, `build`, or generated templ files.
- Run checks proportional to the change and report unrelated baseline failures.
  Add focused tests for nontrivial domain rules and regressions when feasible.

## Communication — all sessions

- Keep chatter minimal; avoid long explanations.
- For each step, give only a 1–2 sentence explanation outside the code when
  necessary.
- Present code in small, skinny blocks; avoid giant walls of code.
- Follow these communication rules consistently throughout every session.
