# Implementation plan

## Affected areas

- `web/frontend/src/islands/bookmanager/BMLayout.tsx`: site link and busy region.
- `BM.tsx`: initial route fallbacks and errors; `books/Book.tsx`,
  `books/index.tsx`, and `books/BookEdit.tsx`: destination loading layouts.
- `ui.tsx`: reuse the existing guard for navigation out of dirty forms.
- Shared `.skeleton` in `common/style/components/util.scss`: inspect before reuse;
  add accessible, motion-aware manager compositions as needed.
- `translations/en.toml`: link and status copy.

## Design and data changes

Read scoped frontend and manager AGENTS.md before implementation. Keep the site
link in the masthead with the existing logo and Your books. Ensure the shell
survives initial pending loads and route errors. Use a real anchor and the
existing unsaved-change flow for document navigation.

Select skeleton composition from the destination route and tab, including direct
hash entry. Distinguish initial/navigation loading from background revalidation;
route loaders already supply chapter data, so no chapter endpoint is required.
Reuse theme tokens and the shared skeleton primitive, reserving realistic cover,
card, and field geometry. Do not introduce timers that delay resolved content.
No backend, schema, migration, or route changes are required.

## Verification

| Criterion | Verification |
| --- | --- |
| AC-1 | Site link on all routes, pending/error views; correct document URL |
| AC-2 | Clean/dirty forms; Stay/Discard; keyboard and modified clicks |
| AC-3 | Throttled direct entry to list, overview, details, and chapters |
| AC-4 | Slow pagination/book navigation and background revalidation; retained edits |
| AC-5 | Success, empty, failure, retry, fast resolution |
| AC-6 | Accessibility tree, keyboard traversal, busy announcements, reduced motion |
| AC-7 | Both themes at 360/768/1280px; overflow and layout movement |

During implementation run `pnpm run build` and `git diff --check`; record baseline
failures separately. Use browser fixtures for pending/error cases and distinguish
fixture evidence from live persistence checks.
