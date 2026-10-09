# Plan

Use the existing editor atom to render a toolbar in the center pane of the parent document, above a flexing iframe wrapper. Use shared Radix Select controls, suppress their automatic trigger refocus after choosing formatting, and preload all offered fonts through FontsLoader. Subscribe to transaction state, add the existing Tiptap font-size extension, and correct heading/list/underline state reporting. Reuse favorite fonts, font loading, and the full font browser. Keep button mousedown from moving selection; Select controls retain the editor selection and restore focus on change. Use scoped theme-token styles with horizontal scrolling on narrow screens.

Verify with the frontend build and live Playwright checks for formatting, selection, active states, and viewport overflow. No API or data migration.
