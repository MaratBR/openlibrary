# Reader settings
User-requested requirements: redesign settings, improve theme selection, vary content width, and obtain font choices from a backend service with a static initial catalog.

## Acceptance criteria
- Theme selection uses accessible preview cards for website, light, and dark.
- Content width adjusts immediately from 48–100 characters in increments of two, defaults to 72, and stays within the viewport.
- Width persists in guest cookies and signed-in account preferences, including editor preview.
- Backend font catalog supplies stable IDs, translated label keys, and CSS font stacks; rendered choices and validation use this catalog.
- Existing saved font, size, page-color, and theme preferences remain compatible.
- Controls fit mobile and desktop, with visible keyboard focus.

## Open questions
None. Existing bundled fonts form the initial catalog.
