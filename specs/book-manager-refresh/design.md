# Proposed screen design

This proposal uses the live Book Manager inspected on 2026-10-03. Your books and book editing are redesigned from the ground up as requested,
with new page composition rather than adjustments to the current table/form.
It changes presentation and existing interactions. No new business capabilities or backend
work are included. See spec.md for observable requirements.

## Priority and visual direction

1. Give Your books a library identity and the editing screen a focused workspace.
2. Make daily actions obvious, accessible, and fast to reach.
3. Build responsive layouts deliberately, with readable fields and durable hierarchy.
4. Complete copy, form navigation, and request feedback.

Retain the warm cream background, white cards, dark green actions, serif page
headings, sans-serif body text, thin borders, and existing rounded controls.
Replace the overview's neon lime statistic tile with a themed neutral tile.
Use one primary action per local task. Keep destructive actions secondary.

## Reviewable visual study

Open [the responsive design study](design-review.html) or review the
[desktop board](design-desktop.png) and [mobile board](design-mobile.png).
The study contains illustrative titles, covers, and counts, clearly labeled as
sample data. It presents two screens in one document for comparison. Controls
are visual demonstrations and do not save or connect to live book routes.

The HTML is a design artifact rather than application implementation. Its literal
colors illustrate the palette; production work must map them to existing theme
tokens. Existing artwork remains intact. No cover replacement/upload is proposed.

## Shared shell

Replace the current tall sidebar with a restrained horizontal manager masthead
on desktop and mobile: existing brand at left, Your books navigation at right.
There is one current global destination, so a full-height navigation column
spends space that the library and editing workspace can use better.

Center content within approximately 1200px. Use 32–40px desktop gutters and 16px
mobile gutters. On book screens, add a Back to books breadcrumb and a consistent
local navigation row: Overview, Details, Chapters, each pointing to its existing
route or view. Keep selected navigation states explicit. No new route is needed.
Page headings use serif typography; controls and data use the existing sans serif.
Use a screen-specific browser title.

## Your books — library grid

Replace the table with a responsive grid of book cards. Use three columns on wide
screens, two at intermediate widths, and one on mobile. A single book occupies
one normal-sized card rather than stretching to fill the grid. Preserve server
ordering and pagination. Do not add placeholder cards for missing books.

```text
OpenLibrary / Book Manager                              Your books

YOUR LIBRARY
Your books                                            [+ Add book]
A place for the stories you’re building.

┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│   cover      ··· │  │   cover      ··· │  │   cover      ··· │
│                  │  │                  │  │                  │
├──────────────────┤  ├──────────────────┤  ├──────────────────┤
│ Public · PG-13   │  │ Hidden · G       │  │ Public · PG      │
│ Book title       │  │ Book title       │  │ Book title       │
│ 12 ch · 24k words│  │ 3 ch · 6k words  │  │ 1 ch · 840 words │
│──────────────────│  │──────────────────│  │──────────────────│
│ Details Chapters│  │ Details Chapters│  │ Details Chapters│
└──────────────────┘  └──────────────────┘  └──────────────────┘
```

Give the cover a softly tinted presentation area, with an approximately 110–120px
wide portrait cover and restrained depth. The tint uses neutral/theme tokens,
not a new cover-color extraction service. Let title and artwork establish the
book's identity; avoid oversized metric chips and repeated green buttons.

Use serif book titles, compact textual status/rating, and one muted count line.
Exact counts are readable for small books; retain abbreviations for large word
totals. Titles wrap up to three lines at usual sizes; exceptionally long names
use accessible truncation rather than breaking the card/action layout. Align
card footers within each grid row without imposing equal title lengths.

Cover and title open the overview. Edit details and Chapters are explicit footer
links to existing destinations. A named actions-menu button holds Trash/Restore,
with the existing confirmation following selection. The menu is keyboard usable,
keeps destructive actions out of the main visual hierarchy, and always remains
reachable on touch; it is not hover-only. Show Trashed when applicable.

Place Add book in the page header. Hide single-page pagination; otherwise place
it below the grid. An empty library uses one calm explanatory panel with Add book,
not a dashboard of invented statistics. Use card-shaped loading placeholders
without fake titles/counts. Do not add search, sorting, or view-mode switches.

## Book overview

```text
← Back to books
Book title                          [Edit details] View public page
Overview | Details | Chapters

[cover]  Public · PG-13
         Tags: Action, Survival, ...

1 chapter        100 words        100 words per chapter

Summary
No summary yet.                       Edit details
```

Show metadata already available in the detail response. Avoid repeating the name
in a definition list immediately beneath the page title. Make absence explicit
rather than leaving a blank summary row. Statistics support the book information
rather than dominating it. The shared local navigation labels the existing General information view Overview.
Remove Analytics until it has actual supported content;
this does not propose implementing analytics. Keep tabs and their selected states
consistent with route state and accessible tab semantics.

## Add book

Retain the four current steps, with a bounded reading column and uniform spacing.
Desktop may show the step list beside the form; mobile shows Step N of 4 and a
compact progress indicator above the fields. Completed steps are buttons.

```text
← Back to books
New book
Step 2 of 4 · Age rating

Age rating
Choose the rating that applies to this book.
[G] [PG] [PG-13 selected] [R] [NC-17] [?]

[Back]                                      [Next]
```

Use permanent Name/Tags labels, not placeholder-only identification. State Tags
(optional), explain selected rating state without relying on color, and use
specific inline validation. The review lists No tags selected when empty and
explains the existing visibility/search behavior. Correct the current “wil” typo.
Use Back and Create book on the final step. Expose no new summary/upload fields.

## Book editing — details workspace

Replace the current single form card with a deliberate two-column workspace.
This is a new composition of existing fields, not merely new section headings.

```text
← Your books
THE GLASS ORCHARD
Book details                                      View public page ↗
Shape how your book appears to readers.
Overview   Details   Chapters
────────────────────────────────────────────────────────────────────
┌──────────────────────────────────┐  ┌───────────────────────────┐
│ The book                         │  │ [cover] Book identity     │
│ Title                            │  │         Existing counts   │
│ [                              ] │  │                           │
│ Tags                             │  │ Audience & visibility     │
│ [tag ×] [tag ×] Add a tag…        │  │ Age rating choices        │
│ Summary                          │  │ Adult                 [ ] │
│ [existing rich text toolbar    ] │  │ Publicly visible      [✓] │
│ [                              ] │  └───────────────────────────┘
│ [                              ] │
└──────────────────────────────────┘
────────────────────────────────────────────────────────────────────
Unsaved changes                              [Cancel] [Save changes]
```

The main panel owns title, tags, and summary, in that order. Give the summary
real writing room (approximately 180–240px minimum in production), with the
existing toolbar kept compact. Preserve every existing formatting action.
Tags stay inline with removable chips and a clearly labeled input. Do not
introduce new inputs, content limits, or summary generation.

The approximately 280–320px side panel shows the current cover as read-only
identity context and the existing age/adult/visibility controls. Use a coherent
selected rating treatment, explicit text and accessible state, rather than a
rainbow competing with the book artwork. Preserve the current rating semantics.
Explain settings using established domain copy; do not imply that visibility
alone publishes drafts or releases chapters.

Use the book title as context above Book details. View public page belongs in
the header. The shared Overview/Details/Chapters navigation lets users move
between existing screens without returning through the library.

A shared action bar below both panels shows clean/dirty/pending/success/error
state, Cancel, and Save changes. Production should keep the action bar reachable
on long forms using a sticky footer within the workspace; it must not overlay
fields, virtual keyboards, or focused controls. Errors also appear beside their
fields. All current controls remain visible; do not hide settings behind accordions.

Below approximately 900px, stack the panels into a full-width form followed by
Audience and visibility. Desktop copy/field widths remain readable; mobile rating
choices wrap and action buttons stay reachable. Long titles, dense tags, empty
summaries, and covers missing from existing data must remain usable. This study
is a layout proposal; it does not add autosave or persistence capabilities.

## Chapters and chapter details

Chapter rows show title, Visible/Hidden state, word count, and scheduled time when
provided. Preserve current server order. Add chapter remains the page action.
Offer Edit details and a separate Open editor link; the latter uses the existing
block-editor destination and does not change that editor.

Add chapter uses a labeled small form with help text, Add chapter, and Cancel.
Keep it close to the trigger on desktop and within the viewport on mobile.

The details panel is approximately 520–640px on desktop, capped by the viewport;
on mobile it fills the viewport. Use Chapter details as heading, the chapter
name as context, a named close button, and a footer grouping Save changes/Cancel.
Separate Open editor visually from metadata save. Retain the existing hidden-to-
public warning and summary controls. Ensure Escape and focus handling work.
Do not present schedules as editable within this panel.

## Trash/Restore confirmation

Use “Move ‘Book title’ to trash?” and explain that restoration is available under
the existing workflow. Buttons read Cancel and Move to trash. Avoid Delete, which
currently conflicts with Trash. Restore uses “Restore ‘Book title’?” with matching
confirm action. Include pending/error state and a meaningful dialog name.

## Limits of the inspection

The account had one book and one visible chapter. Multi-page, empty library,
hidden/scheduled chapter, restore, and successful/failed mutations were not
exercised. No data was changed. Both-theme verification remains pending; the
proposal requires existing theme tokens rather than a new theme system.
