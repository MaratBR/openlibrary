# Tags autocomplete usability

Improve the existing multi-tag picker used in search and book management.

Acceptance criteria:
- Selected tags wrap cleanly and have accessible removal buttons.
- The suggestions stay closed until nonblank search text is entered. Clearing
  the input or selecting a tag closes the list.
- Search suggestions exclude selected tags and display adult markers.
- Arrow keys navigate suggestions; Enter selects without submitting the form;
  Escape dismisses suggestions; tabbing away closes the menu.
- Selecting a tag clears the query and retains input focus.
- Loading, failure, and empty results have translated feedback.
- The input exposes combobox/listbox semantics and active-option state.
