# Plan

Update book.templ and book.scss with a full-width Literata masthead, a cover shelf,
action rail, flat metadata, and a sequential ruled chapter index. Preserve Alpine
state and existing endpoint contracts. Add the first-chapter link using existing
localized text and publicBookURL. Remove the inappropriate TOC listbox role while
retaining the list and tabpanel semantics. No data, API, or migration changes.

Verify public Go tests, frontend build, desktop/mobile/theme browser captures,
chapter URLs, tabs, and overflow. Record unrelated baseline failures separately.
