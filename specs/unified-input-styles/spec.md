# Feature: Unified input styles

## Problem and outcome

Form controls should have one coherent design across the application. Book
Manager currently customizes shared controls through its container styling;
shared `.input`, `.Select`, and compound fields also have differing treatments.
The same control and state should look consistent wherever it appears.

Status: implemented; live application interaction verification remains pending.

## Scope

Unify the shared design of text inputs, native/custom selects, textareas,
password inputs, tag/autocomplete fields, rich-text field frames, and shared
choice controls used by Book Manager (rating choices and switches/checkboxes).
Cover React islands and server-rendered/Alpine forms using these primitives.
Preserve specialized transparent reader/editor controls where their purpose
requires a documented explicit variant. Page/form layout remains local.
No validation, persistence, form-flow, or editor functionality changes are included.

## Scenarios and acceptance criteria

- AC-1: Equivalent controls in manager forms/dialogs and public forms use the
  same default typography, field height, padding, border, radius, surface, and
  placeholder treatment. Compound controls share the same outer field treatment;
  textareas and rich-text fields retain appropriate content height.
- AC-2: Default, hover, focus, invalid, disabled, and read-only states follow a
  shared design in both themes. Keyboard focus and validation remain clearly
  visible; labels, hints, errors, accessible names, and native semantics survive.
  Read-only fields remain readable and selectable where currently supported.
- AC-3: Shared choice controls have consistent checked/selected and focus states
  across consumers. Display-only age-rating badges remain distinguishable from
  interactive choices; visual unification does not change their meaning.
- AC-4: Field appearance is defined by shared control classes/tokens or explicit
  reusable variants, never an ancestor identifying a page or feature. Remove
  equivalent overrides in manager dialogs too. Avoid selectors such as
  `.BM .input`, `.BM-dialog .Select`, or feature-scoped shared choice styling,
  including selectors produced by SCSS nesting. Local wrappers may control
  width, placement, and spacing without changing the field's visual state.
- AC-5: Any necessary size or specialized editor variant is explicitly named,
  reusable, and documented with its intended use. Equivalent controls use the
  same variant regardless of page; a manager-only appearance variant is not a
  substitute for removing container overrides.
- AC-6: Existing form behavior is preserved: entered values, validation,
  submission, password reveal, selection, tag removal/autocomplete, rich-text
  actions, and dirty-form handling. Menus/popovers and fields fit at 360px,
  768px, and 1280px without clipped text or hidden focus.

## Open questions

- None. Use site tokens with the stronger Book Manager-inspired focus treatment:
  a primary outline, separated from the field by a small gap, with a soft halo.
  Apply the same modern treatment to shared fields and interactive choices.
