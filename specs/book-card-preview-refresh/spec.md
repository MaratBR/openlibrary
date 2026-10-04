# Book card preview refresh

Refresh the existing hover preview requested by the user.

Acceptance criteria:
- Short or empty summaries and short titles retain a balanced, consistent-width card.
- Cover (or decorative fallback), title and author establish a clear hierarchy. Long text is bounded; tags follow content without overlapping it.
- Previews fit the viewport and reposition on scroll and resize.
- Leaving an anchor or switching books cannot reveal stale responses. Failed requests remain hidden.

- Desktop placement prefers right, left, top, then bottom, choosing the first side that fits.
- Touch long press (500 ms) opens a centered preview with a blurred backdrop; regular taps still navigate and movement cancels the gesture. Release after long press does not navigate. Tap the backdrop or press Escape to dismiss.

- Library navigation (My library, Archive, Collections) fits narrow screens, wrapping when necessary without page overflow, with a visible current tab.
