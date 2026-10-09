# Tasks

- [x] Record requirements and implementation plan.
- [x] Implement toolbar and selection state (criteria 1–4).
- [x] Add translated labels and responsive toolbar scrolling (criterion 5).
- [x] Verify live formatting and selection through Playwright; inspect desktop and phone screenshots.
- [x] Run Vite production asset build and git diff --check (passed).
- [ ] Full pnpm run build: blocked by existing errors in block-editor/state.ts (missing ./contracts) and islands/header-user-menu/UserMenu.tsx (unused useEffect).

Live checks exercised bold, underline, font family/size, right alignment, heading 4, both list types, returning to paragraph, active states, and absence of bubble menus. Unsaved test edits were discarded by reloading. At 390px the toolbar scrolls to every control; the existing header/sidebar layout crowds out the content area.

Follow-up verified: toolbar bounds match the center pane while sidebars start alongside it; all three dropdowns use shared Select controls and preserve selection. Abel option renders at 20px in Abel and document.fonts confirms it is loaded. Font/size/heading changes passed in the live editor; test edits were discarded. Vite build and diff checks passed; full build has the same baseline errors.

Sidebar follow-up: removed toolbar duplicates from the sidebar only; the shared widget catalog remains available to slash commands. Live sidebar contains Code, Code block, Blockquote, Divider. Vite and diff checks pass; full build retains the same two baseline errors.
