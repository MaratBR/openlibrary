# Feature: Book Manager UI refresh

## Problem and outcome

Authors need clear, usable screens for existing book and chapter management.
Live inspection on 2026-10-03 confirmed unusable narrow layouts, ambiguous Edit
navigation, incomplete tab content, missing accessible control names, and weak
form navigation. Retain the cream surfaces, dark green actions, serif headings,
and existing component primitives while replacing the list and metadata page
compositions from the ground up. Aim for a practical, elegant library and editing
workspace rather than the current table and stacked form.

Status: UI implemented; fixture verification recorded in [verification.md](verification.md).
Implementation authorized on 2026-10-03 on `feature/sdd/bm-1`.
The proposed design, local unsaved-change prompts, and unchanged public-by-default
creation behavior are the implementation baseline. See [live findings](exploration.md),
[screen design](design.md), and the
[navigation guide](../../docs/development/book-manager-navigation.md).

## Scope

Include the existing shell, books list, four-step creation wizard, overview,
book metadata form, chapters list, Add chapter popover, chapter-details panel,
and Trash/Restore confirmation. Preserve existing routes and public visibility
semantics. Display existing DTO fields where useful.

Exclude the block editor and its content editing, autosave, drafts, publishing,
and scheduling workflows; preserve its existing entry link. Summary fields in
manager forms remain in scope. No new search/filter/sort, analytics, cover upload,
collection management, reorder/delete/bulk operations, backend endpoints, schema
changes, or domain features are proposed. Earlier backend/search proposals are
superseded by this UI scope.

## Scenarios and acceptance criteria

- AC-1 — Responsive shell: at 360px, 768px, and 1280px every in-scope view fits
  the viewport without page-level horizontal scrolling. Replace the permanent
  sidebar with a compact horizontal manager masthead, using existing brand and
  Your books navigation. Center desktop content with generous gutters; use
  16px mobile gutters. Titles wrap naturally and Add book remains reachable.
- AC-2 — Library: replace the table with a cover-led grid of three, two, or one
  columns according to available width. Cards display existing cover, linked
  title, status/rating, and readable chapter/word counts, with direct Edit details
  and Chapters links. Cover/title open overview. A named keyboard/touch-usable
  actions menu exposes existing Trash/Restore and its confirmation. Single books
  retain normal card width; long titles and missing covers preserve layout.
  Counts have correct plural copy. Preserve server ordering/pagination and hide
  the pager when only one page exists. Do not invent book status or statistics.
- AC-3 — List states: empty library offers Add book; pending and failed loads
  are distinguishable from no books. Failure offers retry. This iteration adds
  no search UI or zero-search-results state.
- AC-4 — Confirmation: use Trash consistently in the trigger, title, explanation,
  and confirmation. Name the affected book and explain the existing reversible
  behavior without implying permanent deletion. Restore uses matching language.
  Cancel changes nothing. Pending confirmation prevents repeated requests;
  failures retain the dialog and existing row state. Only confirmed success
  updates the row and reports success.
- AC-5 — Overview: use consistent Overview, Details, and Chapters navigation across
  book screens; Overview labels the existing General information view. Remove the empty
  Analytics tab; old analytics links fall back to General information. A URL
  without `t` still renders General. The header provides a labeled Back to books
  action and Edit details. Display existing cover, rating, tags, visibility,
  summary, and chapter/word counts; missing summary reads No summary yet.
  Statistics use neutral themed surfaces instead of the isolated lime gradient.
- AC-6 — Wizard: retain Name, Age rating, Tags, Review. Show Step N of 4, a
  persistent field label, concise guidance, and Back/Next actions. Completed
  steps are keyboard-operable buttons and preserve entered values; future steps
  cannot bypass validation. Enter follows the same validity rule as Next; an
  empty/whitespace-only name cannot advance. Ratings use the server-provided
  options and a visible selected state. Tags are explicitly optional.
- AC-7 — Review: show name, rating, tags or No tags selected, and a clear note
  about existing public visibility. Preserve the current explanation that an
  empty book does not appear in search until its first chapter exists, with
  corrected copy. Allow Back to correct input. Use Create book as the sole
  submission action and show a pending state that prevents repeated submission.
  Do not claim success until the existing creation flow confirms it.
- AC-8 — Editing workspace: replace the single metadata card with a main
  title/tags/summary panel and a narrower identity/audience/visibility panel,
  stacked below approximately 900px. Existing cover is read-only context.
  Rating choices have one coherent selected treatment and existing semantics.
  Book title provides context for the Book details heading; View public page is
  in the header. Overview/Details/Chapters navigation targets existing views.
  Save changes and Cancel share a reachable action bar with accurate request/
  unsaved feedback. Sticky behavior must not obscure input or keyboard focus.
  Cancel returns to overview. Existing values and rich text actions are retained;
  Enter/Save use one controlled submit path. No new fields or autosave are added.
- AC-9 — Chapters: retain server order, Add chapter, metadata editing, and the
  existing editor destination. Show chapter name, explicit Visible/Hidden status,
  word count, and scheduled time with timezone when supplied. Use Edit details
  for the panel and a distinct Open editor link to the existing destination.
  Empty chapter lists explain how to start and offer Add chapter.
- AC-10 — Add chapter: keep a lightweight form with a visible Chapter name label,
  limit/help text matching the current accepted client rule, Add chapter and
  Cancel. It fits mobile, places focus in the input, and uses identical Enter/
  button validation. Pending prevents repeat requests; failure preserves input
  and form visibility. Only confirmed success dismisses and refreshes the list.
- AC-11 — Chapter panel: title reads Chapter details with chapter name as context.
  Existing name, summary, visibility, and publication warning remain. Desktop
  uses a readable bounded panel; mobile uses the full viewport width. Save
  changes/Cancel are grouped; Open editor is visually separate. Close has an
  accessible name, Escape dismisses, focus enters the panel and returns to its
  trigger, and background interaction is contained. Failure retains input;
  only confirmed success closes and updates the chapter list.
- AC-12 — Accessible copy and controls: all visible strings and accessible names
  are translated; no raw keys such as `bookManager.edit.analytics` or
  `search.removeTag` appear. Back, close, tag removal, summary toolbar buttons,
  switches, rating options, and tabs have meaningful accessible names/states.
  Keyboard navigation and focus are visible; validation errors are associated
  with fields and status feedback is announced. Verify both existing themes.
- AC-13 — Leaving forms: local unsaved-change protection offers Stay
  or Discard when Cancel, back navigation, panel dismissal, or Open editor would
  discard changed fields. It preserves input when staying and does not appear
  on clean forms or after confirmed save. No persistent drafts/autosave are added.
- AC-14 — Compatibility: existing hash routes, default redirect, catch-all,
  authorization, current endpoints, and block-editor links keep working. Loading
  or unavailable book details offers a readable recovery state and Back to books.
  The UI must not invent data or advertise unsupported actions.

## Decisions and remaining limitations

- Use the screen layouts and copy in design.md.
- Include AC-13 local Stay/Discard prompts. Browser reload/close uses the native
  browser warning; browsers control its wording.
- Explain existing public-by-default creation in the review step.
- Chapter title validation currently differs between client (70 characters) and
  server (200). Keep current UI rules in this scope and track reconciliation
  separately if backend behavior changes are desired.
- Creation failure currently renders a separate server error page. Retaining all
  wizard input on such failure is deferred pending a separate contract decision;
  this UI proposal adds no endpoint and does not promise unsupported recovery.
