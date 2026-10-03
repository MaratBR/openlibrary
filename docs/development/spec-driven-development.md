# Spec-driven development

Use a spec for new features and substantial changes to existing behavior. A
small fix with an obvious expected result does not need one. The spec captures
what users and other callers should observe; implementation details belong in
the plan.

## Start a feature

From the repository root, run:

```sh
./scripts/new-spec.sh book-search-filters
```

This creates `specs/book-search-filters/` with `spec.md`, `plan.md`, and
`tasks.md`. Use a short, lowercase, hyphenated name. Keep the three files with
the implementation in the same change so a reviewer can compare the intended
behavior with the code.

## Work through the files

1. Write `spec.md` first. Describe the problem, scope, user scenarios, and
   acceptance criteria as observable outcomes. Cover empty, invalid, and error
   cases where they matter. Record open questions explicitly; resolve any that
   affect behavior before coding.
2. Write `plan.md` against the agreed spec. Identify the affected project
   areas, important design choices, data or API changes, rollout or migration
   needs, and how the acceptance criteria will be verified. Follow the scoped
   `AGENTS.md` files for the directories being changed.
3. Break the plan into small, ordered checkboxes in `tasks.md`. Link tasks to
   acceptance criteria so the implementation can be checked against the spec.
4. Implement and verify each slice. Update the spec when a requirement changes,
   then update the plan and tasks to match. Mark tasks complete only when their
   code and relevant checks are done.

At review time, every acceptance criterion should have an implementation and a
verification path. Report any check that could not run, and leave unfinished
tasks unchecked. The repository's normal checks and generated-source rules are
listed in [CONTRIBUTING.md](CONTRIBUTING.md).
