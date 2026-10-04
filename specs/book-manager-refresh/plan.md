# UI implementation plan

## Preconditions

Implementation was authorized on 2026-10-03 on `feature/sdd/bm-1`.
Use the proposed layouts and include local unsaved-change protection.
No backend additions are included. Live inspection evidence is in exploration.md.

## Affected areas

- `BMLayout.tsx`: horizontal masthead using existing brand and Your books navigation.
- `BM.tsx`: preserve routes; usable route error recovery.
- `books/index.tsx`: library grid, card anatomy, actions menu, state presentation,
  pagination visibility, and confirmation copy.
- `books/Book.tsx` and `BookGeneral.tsx`: supported tabs and richer overview using
  existing DTO data.
- `books/BookEdit.tsx`: two-column details workspace, shared local navigation, action bar and feedback.
- `new-book/NewBookForm.tsx`: wizard layout, labels, step buttons and navigation.
- `books/BookChapters.tsx` and `ChapterSlidePanel.tsx`: list, popover, panel,
  status labels and actions.
- Shared controls only where necessary for names/focus/responsiveness; avoid
  unsolicited changes to every consumer.
- Translations: manager copy and existing missing accessible labels.

Read the nearest AGENTS.md before editing each area. Preserve the BM island
entry point, existing APIs and DTOs, theme tokens and shared controls.

## Design decisions

Use design.md and design-review.html as the new composition reference. Build a
horizontal masthead and centered main region; remove the full-height sidebar.
Use three/two/one library columns according to available width. Keep one-book
layouts bounded; validate full pages of 20 cards as well as empty and short pages.

Use an accessible book actions menu for existing Trash/Restore. Introduce no new
API calls for metrics/status. Metadata uses a flexible main panel plus a roughly
300px audience panel, stacked below approximately 900px; scope the sticky action
bar to the editing workspace. Shared Overview/Details/Chapters links map to current
routes/views, and dirty-form handling applies to these links as accepted.

The static visual study contains illustrative copy, covers and data; production
uses real DTO data, existing BookCover and all existing rich text actions. Map
study colors and type to app tokens/fonts, keep dark theme support, and translate
all copy. Verify long titles/tags and unsupported image cases. The HTML artifact
is a review aid and is not embedded into the application.

Display book and chapter status from current data; never infer publication from
word count or visibility alone. Library cards use a consistent grid on desktop and a single column on mobile. Pagination remains tied to existing results and is hidden for one page.
Audit Pagination's size argument before changing it; source inspection alone
is not sufficient to conclude that it represents the API page size.

Retain native book creation and current endpoints. Improve validation/navigation
and pending display without introducing a new create API. Book/chapter mutations
must check responses before presenting success; revalidate existing loaders.
Track any confirmed backend persistence or authorization defect separately.

Reuse summary rich text controls; their accessible toolbar names are in scope,
while block-editor internals are excluded. Dirty-form prompts, as accepted, track
local field changes only. Match UI validation guidance to actual accepted rules;
no server validation changes belong in this plan.

## Verification

| Criterion | Check |
| --- | --- |
| AC-1 | 360/768/1280px shell and every screen; document width equals viewport |
| AC-2 | Current book, long title/missing cover fixtures, mobile cards, single/multi-page pagination |
| AC-3 | Empty, slow, failed list states |
| AC-4 | Trash/restore dialogs; cancel; pending; confirmed success and failure using disposable data |
| AC-5 | Default/General/Chapters/old Analytics URLs; existing metadata and empty summary |
| AC-6 | All four steps; Enter/Next parity; empty/whitespace names; keyboard step navigation |
| AC-7 | Review with/without tags; back editing; current visibility explanation; pending creation |
| AC-8 | Initial values; Enter/save path; cancel; success/failure presentation |
| AC-9 | Empty/populated chapter list; visibility/schedule rendering and unchanged editor destination |
| AC-10 | Popover label, validation, Enter, Cancel, pending/failure behavior and mobile fit |
| AC-11 | Panel width, keyboard focus, Escape/close/cancel, save feedback and chapter refresh |
| AC-12 | Snapshot names/states, keyboard traversal, translation keys, both themes |
| AC-13 | Clean/dirty dismissal and navigation as accepted |
| AC-14 | Direct routes, anonymous redirect, detail failure, editor destination |

For frontend changes run pnpm run build and git diff --check. Use existing
frontend tests for nontrivial state transitions when available; document browser
checks otherwise. Run focused browser checks against fixture responses when the local application is unavailable; distinguish these from live persistence/auth checks.
Actual exploration did not submit any create/save/trash mutation; those outcomes
remain unverified and require disposable data during implementation verification.

## Implementation choices

Use native modal dialogs within the island for chapter details and confirmations.
They provide focus containment, inert background, Escape handling, and focus
restoration without changing shared Modal/SlidePanel consumers. React Router
blocking protects hash navigation; local Stay/Discard prompts protect dismissals
and editor links; beforeunload protects reload/close.

Pagination size is the number of page links, not the API page size. Keep requests
at 20 books and display up to five page links for narrow screens.

See [verification.md](verification.md) for implementation checks and remaining live checks.
The action bar stays in place during clicks; focused fields scroll above it,
including after visual-viewport resize for a virtual keyboard.
