# Shared form controls

Use `.input` for text, password, and native fields; `.textarea` or `textarea.input`
for multiline fields; and `.Select` for native or Radix selects. Their shared
appearance lives in `web/frontend/src/common/style/components/input.scss`, using
the site surface, foreground, border, primary focus, and destructive tokens. Tag
fields (`.input-dropdown`) and rich-text frames (`.OlSimpleEditor`) use the same
outer treatment and show focus when an inner control receives focus.

Set `aria-invalid="true"` on the field or its inner input for validation styling.
Native disabled attributes and `data-disabled` style disabled controls. Read-only
inputs and noneditable rich-text content use a subdued surface while preserving
native selection behavior. Keep labels, hints, and errors associated with their
controls; styles do not supply accessible names or validation behavior.

Page wrappers may set width, spacing, and placement. Do not override shared field
colors, borders, focus, typography, or radii through page/feature ancestors or
utility classes. Use an explicit shared variant for a functional difference:

| Variant | Intended use |
| --- | --- |
| `.Select--sm` | Compact selects in toolbars, such as comment sorting; 36px minimum height |
| `.input--search` | Prominent search entry, with larger padding and pill frame |
| `.input--code` | Verification code entry, with large monospaced text and character spacing |
| `.input--joined-start`, `.input--joined-end` | Adjacent fields such as numeric range endpoints; logical inner corners are square |
| `.OlSimpleEditor--long` | Longer summaries; 200px minimum content height |
| `.OlSimpleEditor--review` | Review composition; 300px minimum content height |
| `.input-transparent` | Inline title editing, where the surrounding heading supplies typography and the frame stays transparent |

Textareas may set a local minimum content height. Rich-text toolbars wrap at narrow
widths; the editor content receives the shared outer focus ring.

Interactive age ratings use a keyboard-focusable `.age-rating-input` followed by
`.age-rating`. The adjacent badge shows selection and focus. Standalone
`.age-rating` badges retain their content-rating colors. Switches retain native
checkbox inputs and pair them with `.Switch-slider`; native `.checkbox` controls
share their primary accent and primary keyboard focus treatment.

Focus uses the shared `control-focus` mixin: a 2px primary outline with a 3px
gap and soft halo. Invalid fields use the destructive token for the same treatment.
Field frames have a subtle resting shadow; editor toolbars use a subdued surface.
