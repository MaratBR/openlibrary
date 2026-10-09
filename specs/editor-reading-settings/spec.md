# Feature: Editor reading settings

## Problem and outcome
Authors need to change the reading preferences inherited by the chapter editor.

## Scope
Reuse the reader settings and shared account/device persistence. These controls change presentation, not saved chapter formatting.

## Scenarios and acceptance criteria
- AC-1: A muted icon button fixed at the top right of the editor canvas opens settings. The editor offers font size, serif/sans/dyslexic font, background/surface page color, and website/light/dark theme controls.
- AC-2: Preferences initialize from the account when available, otherwise device cookies, and changes persist through the existing reader preference API and cookies.
- AC-3: Changes immediately affect the editor canvas without modifying draft content or explicitly formatted text.
- AC-4: The panel closes using its close button, backdrop, or Escape and fits narrow viewports.

## Open questions
None.
