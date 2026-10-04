# Implementation plan

## Affected areas

- `web/frontend/src/common/style/components/input.scss`: shared field/state rules.
- Related shared password, rating, switch/checkbox, and rich-text styles.
- `web/frontend/src/components`: FormControl, TagsInput, select/password and
  rich-text consumers as identified by the inventory.
- `web/frontend/src/islands/bookmanager/manager.scss` and manager forms/dialogs:
  remove visual control overrides and adopt shared classes/variants.
- Public/admin templates and Alpine forms using these primitives, as required.
- `docs/development`: document control/variant conventions during implementation.

## Design and data changes

Inventory control classes, states, inline styles, utility overrides, and nested
SCSS across consumers before editing; do not assume `.BM .input` literally exists.
Current manager styling scopes focus, rating selection, and editor rules under
`.BM`/`.BM-dialog`; shared select defaults also differ from text inputs. Classify
local layout separately from field appearance and move reusable control rules
into the appropriate shared component stylesheet.

Use the current shared input treatment and theme tokens as the baseline. Align
native/custom selects and compound field frames while preserving their internal
layout and interaction semantics. Scope rating-choice rules to a shared
interactive primitive so static badges retain display styling. Keep intentional
content sizing as documented shared editor variants where needed.

Adopt explicit reusable variants only for a demonstrated functional need. Audit
both React and templ consumers for regressions; this feature owns shared control
appearance needed by the other two proposals. No API, DTO, schema, or migration
changes are required. Read nearest AGENTS.md before implementation in each area.

## Verification

| Criterion | Verification |
| --- | --- |
| AC-1 | Side-by-side equivalent controls in manager, dialogs, public/admin forms |
| AC-2 | State matrix in both themes; keyboard, invalid, disabled/read-only checks |
| AC-3 | Interactive choices versus display badges; checked and keyboard states |
| AC-4 | Source and compiled selector audit including nested SCSS and utility overrides |
| AC-5 | Review shared variant definitions, documentation, and consumers |
| AC-6 | Representative form interactions and compound menus at 360/768/1280px |

During implementation run `pnpm run build` and `git diff --check`; regenerate
and run focused Go checks if templ files change. Record baseline failures and
browser evidence. Verification must cover consumers outside the manager.
