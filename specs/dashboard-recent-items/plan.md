# Implementation plan

Use a shared DashboardNavItem component for all three DashboardShell layouts. Typed route handles provide labels from successful loader data. Use Jotai atomWithStorage with Effect Schema validation of localStorage data in a separate state module, keyed by dashboard and section, for a two-item deduplicated list in each navigation group, linking detail subpages to the canonical entity route. Style nested lists with the shared shell's theme tokens and truncate long names.

No API or translation changes are needed. Persistence is browser-local. Verify recency transitions with focused assertions, run the frontend build and whitespace check, and report unrelated baseline failures.
