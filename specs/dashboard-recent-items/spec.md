# Dashboard recent items

Show shortcuts to opened detail pages beneath their DashboardShell list section.

## Acceptance criteria

- Admin books, users and tags; moderation books, users and reports; and Book Manager books show the loaded item name beneath their parent section.
- Each section retains at most two items, newest first. Opening a third removes the oldest. Reopening an item moves it first without duplication.
- Detail subpages and edit pages count as the same item and link back to its main detail page.
- List, create, search, placeholder and failed detail pages do not add items.
- Shortcuts survive dashboard route navigation, independently per section, persist in localStorage across reloads using Jotai, separately for each dashboard. Invalid stored data falls back to an empty list.
- Long labels fit the sidebar, expose the full label on hover, and links support keyboard navigation and active indication.

The user's requested submenu and two-item limit are the agreed requirements; recency and separate dashboard histories are implementation defaults.
