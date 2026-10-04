# Feature: Book Manager site navigation and loading skeletons

## Problem and outcome

Authors need a clear way to return to the main site and useful visual feedback
while manager data loads. The current masthead only links within the manager;
initial route loading and subsequent navigation show a text loading message.

This is a follow-up to [Book Manager refresh](../book-manager-refresh/spec.md).
Status: implemented; browser acceptance checks and full build verification remain pending.

## Scope

Include the manager masthead, initial data loads, navigation between books,
pagination, overview, details, and chapters. Preserve existing routes, API
contracts, errors, empty states, mutation feedback, and unsaved-change protection.
The block editor and public-page loading are outside this feature.

## Scenarios and acceptance criteria

- AC-1: Every manager screen, including loading and route error states, provides
  a visible, translated Back to main site link to `/`. It uses normal document
  navigation rather than a manager hash route. Your books remains a separate link.
- AC-2: Leaving a dirty form through the site link offers the existing Stay or
  Discard choice. Stay retains the form and values; Discard navigates to `/`.
  Clean forms navigate directly. Keyboard and modified-click behavior follow
  the existing guarded external-navigation convention.
- AC-3: Initial library loading shows cover/card skeletons; overview loading
  shows cover, heading, and metadata placeholders; details loading shows the
  form layout; loading a book with the chapters tab selected shows list-shaped
  placeholders. Placeholders match the destination's responsive composition.
- AC-4: Navigation to another book, page, or unloaded view shows the destination
  skeleton until its data resolves. Old book data cannot appear under the new
  destination or accept edits. Revalidation of already loaded data retains
  content and form values with a nonblocking busy indication; saves, creation,
  trash, and restore retain their existing pending controls and feedback.
- AC-5: Successful loading replaces placeholders with real data. Empty results
  show the existing empty state only after loading; failure replaces the skeleton
  with the existing recovery state and retry action. Retry can show skeletons
  again. Fast loads leave no stuck placeholders or artificial delay.
- AC-6: Skeletons contain no fake data or focusable controls and are hidden from
  assistive technology. The affected region exposes busy state and one translated
  loading announcement. Completion clears busy state. Reduced-motion preference
  disables animated shimmer/pulse.
- AC-7: At 360px, 768px, and 1280px, in both themes, the masthead links remain
  reachable, skeletons fit without horizontal overflow, and loaded content uses
  the same major layout dimensions to limit movement.

## Open questions

- None. Main site means `/`; the manager logo retains its current destination.
