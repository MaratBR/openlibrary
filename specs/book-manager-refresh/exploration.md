# Live Book Manager exploration

Inspected 2026-10-03 using Playwright against localhost:8080 after Chrome became
available. Earlier browser-launch failures are resolved. Logged in with the local
account supplied in chat; credentials are not stored in project documentation.

## Observed navigation

- Opening `/books-manager#/books` while logged out redirected to
  `/login?next=%2Fbooks-manager#/books`. Username/Password and Sign in were visible.
- Signing in returned to the books list. The sidebar label is Your books.
- The account contained one book, “tewt”, with one chapter, “Test chapter”.
- List Edit opens overview, whose second Edit opens metadata.
- General information without a tab query displayed statistics and summary.
- Chapters showed its existing row; Edit opened the details panel. Add chapter
  opened a small textbox/popover. No creation was submitted.
- The four new-book steps were explored using unsaved input; Create was not clicked.
- Trash opened confirmation and Cancel dismissed it; no trash request was sent.
- Block-editor navigation was inspected as a link; the editor was not entered.

## Confirmed UI findings

| Finding | Observation | Proposed response |
| --- | --- | --- |
| Mobile shell consumes main column | Permanent sidebar remains 220px at 360px; headings wrap into vertical fragments | Compact mobile header and full-width main body |
| Mobile overflow | At 360px: list document width 676px; wizard Tags 801px; metadata 520px; chapter panel 450px | Flexible controls/cards; full-width mobile panel |
| Mid-width list pressure | Initial 784px screenshot showed Trash extending past the card edge | Wrap actions and reduce cover/metric footprint |
| Misleading Edit | List Edit opens overview, requiring another Edit to reach fields | Open book and Edit details labels/destinations |
| Pagination noise | Single-page list shows a bare “1” above its only row | Hide pager when one page |
| Awkward metric copy | “1 chapters” and “0.1k words” displayed for a small book | Correct plurals and readable small counts |
| Incomplete analytics | Tab displays `bookManager.edit.analytics`; selecting it produces an empty body | Remove unsupported tab; fallback for old links |
| Sparse overview | Repeats book name; empty summary row; bright lime word tile dominates | Use existing metadata, explicit empty summary, neutral statistics |
| Title validation bypass | Empty Name with Next disabled advanced to Age rating on Enter | Same validation for keyboard/button navigation |
| Wizard back navigation | Completed list items are clickable; no Back action | Accessible step buttons and explicit Back |
| Review gaps | Empty Tags has no explanation; note contains “wil”; no visibility disclosure | No tags selected, corrected copy, existing-visibility note |
| No leaving warning | Navigation away from filled unsaved wizard proceeded immediately | Proposed local discard/stay prompt |
| Metadata hierarchy | One large form, Save alone, View public page in separate card | Group fields and action footer; public link in header |
| Trash/Delete mismatch | Trigger/confirm Trash but heading “Delete the book?” and description “delete” | Consistent reversible Trash copy naming book |
| Panel keyboard behavior | Escape with chapter panel open did not dismiss it | Explicit Escape/focus behavior |
| Accessible names | Back/close/rich-text toolbar buttons unnamed in snapshots; switches have no name; tabs lack aria-selected | Labels, names, selected states and panel/dialog semantics |
| Missing translations | Tag removal names show `search.removeTag` | Translate existing labels |
| Browser title | All manager routes displayed Home | Screen-specific titles |

## Evidence

Local Playwright artifacts are in `.playwright-mcp/`. Screenshots used during review:

- `page-2026-10-03T10-57-52-671Z.png`: list at initial 784px viewport.
- `page-2026-10-03T10-58-10-807Z.png`: metadata at 1280px.
- `page-2026-10-03T10-58-32-523Z.png`: chapter panel at 1280px.
- `page-2026-10-03T10-59-46-943Z.png`: Tags step at 360px.
- `page-2026-10-03T11-00-31-194Z.png`: books list at 360px.
- `page-2026-10-03T11-01-31-650Z.png`: chapter panel at 360px.

Other screenshots captured the chapter panel, creation step, mobile Tags, mobile
list, and mobile chapter panel. These generated exploration artifacts are not
part of the intended spec deliverable; findings above are the durable record.

## Untested behavior and source-only concerns

No book/chapter create, save, trash, or restore mutation was submitted. Empty,
multi-page, hidden/scheduled states, request failures, and dark theme were not
verified. Actual persistence and success/error feedback remain untested.

Source inspection separately identifies unchecked mutation responses, possible
native metadata form submission, and mismatched chapter-title limits (client 70,
server 200). Do not describe these as reproduced server defects. Console messages
included missing flash-message initialization and controlled-field warnings; no
root-cause investigation was performed as part of this UI task.
