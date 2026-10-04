# Book Manager navigation

The Book Manager refresh is implemented on `feature/sdd/bm-1`. The
[feature spec](../../specs/book-manager-refresh/spec.md) describes its behavior;
[verification notes](../../specs/book-manager-refresh/verification.md) distinguish
fixture browser checks from live server checks.

## Open the manager

Start the application using the README, sign in at `/login`, and open
`/books-manager#/books`. Use the public login rather than the admin interface.
The horizontal masthead provides Your books on every screen.

## Books and chapters

- Cover and title open Overview. Each library card also has Edit details and
  Chapters links. Its named actions disclosure contains Trash or Restore.
- Trash and Restore require confirmation and name the affected book. Restore
  leaves the book hidden; change visibility in Book details when appropriate.
- Overview, Details, and Chapters navigation is shared across book screens.
  Missing or old `t=analytics` parameters show Overview.
- Book details contains title, tags, summary, audience, and visibility controls.
  Save changes submits once and shows request feedback. Cancel returns to Overview.
- Chapters lists visibility, word counts, and any scheduled time with timezone.
  Add chapter opens a small labeled form; Enter follows the button's validation.
- Edit details on a chapter opens the details panel. Escape, Close, and Cancel
  dismiss it. Focus enters the panel and returns to the trigger on dismissal.
- Open editor retains the existing block-editor destination. The block editor
  itself is outside the refresh scope.

Changed forms offer Stay or Discard before local navigation or dismissal. Browser
reload/close uses the browser's own warning. Pending saves prevent repeated
submission and dismissal. Failed requests retain input; inspect inline feedback.

## Create a book

Select Add book and complete Name, Age rating, optional Tags, and Review. Back
and completed-step buttons preserve values. The review explains existing public
visibility and that an empty book is absent from search until its first chapter
exists. Create book uses the existing native POST and redirect. Server creation
errors still render a separate error page; recovery of wizard input on that page
is outside the current contract.

## Direct routes

Replace IDs with real values from links. Query parameters for manager views
belong after the hash.

| View | URL path |
| --- | --- |
| Books | `/books-manager#/books` |
| New book | `/books-manager#/books/new` |
| Overview | `/books-manager#/books/BOOK_ID` |
| Overview (explicit) | `/books-manager#/books/BOOK_ID?t=general` |
| Chapters | `/books-manager#/books/BOOK_ID?t=chapters` |
| Book details | `/books-manager#/books/BOOK_ID/edit` |
| Public book | `/book/BOOK_ID` |
| Block editor | `/books-manager/book/BOOK_ID/chapter/CHAPTER_ID` |

The chapter-details panel does not have a URL state. Open it from its chapter row.
The API still accepts up to 200 chapter-title characters while the existing UI
rule permits up to 70. The refresh retains that client rule.
