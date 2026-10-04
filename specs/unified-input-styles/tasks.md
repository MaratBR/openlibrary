# Implementation tasks

- [x] Inventory shared controls, variants, states, and feature overrides (AC-1–5).
- [x] Consolidate shared field tokens and default/state styling (AC-1, AC-2).
- [x] Align selects, password/tag fields, and rich-text frames (AC-1, AC-2, AC-6).
- [x] Move interactive rating/choice appearance into shared primitives (AC-3, AC-4).
- [x] Remove manager/dialog appearance overrides and migrate affected consumers (AC-4, AC-6).
- [x] Document necessary reusable variants and styling conventions (AC-5).
- [ ] Verify state/theme matrix, external consumers, and form interactions (AC-1–3, AC-6).
- [x] Audit compiled selectors, run implementation checks, and record evidence (AC-4).

## Implementation evidence

- Inventory: shared fields/selects and dropdowns in `input.scss`; password inputs
  already use `.input`; React and Alpine editors use `.OlSimpleEditor`; ratings,
  switches, and native checkboxes have shared styles. Book Manager scoped focus,
  editor sizing, and rating appearance to `.BM`/`.BM-dialog`. Public search/code,
  admin tag input, and search range utilities supplied competing field appearance.
- Shared field/state rules now cover selects, tag frames, and rich-text frames.
  Manager appearance overrides were removed; summary height uses the reusable
  `OlSimpleEditor--long` variant. Rating selection is shared and static badges keep
  content colors. See `docs/development/form-controls.md` for variants.
- `pnpm run build`: blocked by existing TS2307 in `block-editor/state.ts`
  (`./contracts` missing) and TS6133 in `header-user-menu/UserMenu.tsx`
  (unused `useEffect`). Neither file was changed by this feature.
- `pnpm exec vite build --mode production`: passed, including SCSS/UnoCSS output.
- `templ generate`: passed; warns that installed generator v0.3.960 is older than
  the module's v0.3.1020. Generated files remain ignored.
- `GOCACHE=/tmp/openlibrary-go-cache GOTMPDIR=/tmp go test ./web/public/...`:
  passed. `git diff --check`: passed.
- Source and compiled `dist/common.css` selector audits found no `.BM` or
  `.BM-dialog` descendants styling shared controls or focus.
- Playwright checked a temporary fixture using compiled production CSS with
  manager, dialog, and public wrappers, in light/dark at 360/768/1280px: no
  horizontal overflow; text/select typography 16px/24px, minimum height 44px,
  padding 12px, radius 8px. Hover, brass focus ring, destructive invalid focus,
  disabled opacity, and read-only surfaces were checked. Keyboard focus showed
  solid outlines on rating, switch, and checkbox in both themes. Native typing,
  selection, contenteditable editing, and checked values worked.
- Remaining verification: real React/Radix menus, tag addition/removal and remote
  autocomplete, rich-text toolbar actions, submission/validation and dirty-form
  behavior on live manager/public/admin routes. The fixture verifies styles and
  native behavior, not application workflows; the verification task stays open.

### Visual refinement

- Adopted the stronger Book Manager-inspired focus outline as the shared design:
  2px primary outline, 3px gap, and soft halo. Invalid focus uses destructive
  colors. Fields and interactive rating choices have subtle resting depth;
  selected ratings use a matching primary border and editor toolbars use the
  secondary surface. Joined range fields raise the focused start field so its
  outline remains visible above its neighbor.
- Rebuilt Vite assets and repeated compiled selector audit and `git diff --check`:
  passed. Full `pnpm run build` still reports the same two unrelated TS errors.
- Repeated Playwright fixture checks in both themes at 360/768/1280px: text,
  native selects, tag frames, editor frames, invalid fields, and keyboard-focused
  ratings/switches/checkboxes all show the 2px outline with 3px gap; no horizontal
  overflow. Live application workflow verification remains pending as above.
