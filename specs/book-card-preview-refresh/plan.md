# Plan

Update the preview fragment and scoped island styles using existing theme tokens. Keep existing endpoints and hover delay. Measure the popup before positioning; use Floating UI offset, flip, shift and autoUpdate. Invalidate pending hover loads on leave and disposal. Verify template rendering and frontend build; check whitespace.

Use explicit Floating UI fallback placements. Add touch pointer long-press detection with movement cancellation and click suppression. Center the mobile preview in a fixed backdrop layer, lock document scrolling while open, and restore it on dismissal/disposal.

Adjust the shared library header styles to use responsive typography and wrapping. Match current-tab styling to the existing aria-current attribute. Verify all tabs fit narrow viewports with short and longer translated labels.
