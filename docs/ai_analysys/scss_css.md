# SCSS and CSS findings

## Overall assessment

The styling direction is stronger than the raw file count suggests. Theme values are centralized in CSS custom properties, UnoCSS maps semantic tokens, and reusable controls have component classes. The main maintainability issue is boundary enforcement: global components, page styles, utility output, fonts, and legacy rules are assembled into one broad entry without automated checks or size constraints.

## High

### CSS-01: The global stylesheet imports the entire component and page catalog

Evidence: `web/frontend/src/common/style/index.ts:1-9` imports the component barrel, Uno output, common rules, `home.scss`, and four font styles. The component barrel imports 39 files (`web/frontend/src/common/style/components/index.scss:1-39`), including page/domain-specific styles such as book reader, chapter comments, library, profile, recommendations, search, and editor rules.

The verified production build produced `dist/common.css` at 972.64 kB (595.97 kB gzip). Some of this is font and asset data, but it demonstrates that every page pays for a broad global graph.

Impact: ownership and deletion are hard to determine, unused styles accumulate, page changes can affect unrelated routes, and the common asset becomes a performance and deployment liability.

Recommendation: define three layers: foundations (reset/tokens/base), genuinely shared components, and route/island styles loaded by their owning entry. Move `home.scss` and domain-specific component files out of `common`. Keep fonts as explicit assets and load only required families/weights. Add a CSS bundle budget and per-entry size report.

## Medium

### CSS-02: There is no automated style quality gate

Evidence: package scripts contain no Stylelint, formatting check, or CSS architecture validation. SCSS filenames and naming styles already vary (`ImageUploader.scss`, `ProseMirror.scss`, `_content-rating.scss`, kebab-case files), and the component barrel order is manual.

Impact: dead selectors, invalid declarations, specificity growth, inconsistent naming, and ordering drift rely on reviewers to detect them.

Recommendation: add Stylelint with SCSS support, property ordering if the team values it, duplicate-selector checks, and a narrow rule for disallowing literal colors outside token/foundation files. Add Prettier/style checks to CI. Standardize lowercase kebab-case filenames and document whether partials use `_`.

### CSS-03: The boundary between design tokens, component CSS, Uno utilities, and inline style is implicit

Evidence: semantic theme tokens are well defined, but literal visual values still occur in component SCSS (for example modal/input shadows and loader/star colors), template inline styles, and arbitrary Uno values. Some are inherently dynamic, but no policy or lint rule distinguishes justified dynamic styles from one-off constants.

Impact: equivalent values drift, dark-theme support becomes uneven, and contributors cannot predict whether a change belongs in a token, SCSS component, utility, or markup.

Recommendation: document a short decision rule and enforce the easiest parts: semantic color literals only in `theme.css`; dynamic geometry may remain inline; reusable visual state belongs with the component; route layout uses utilities. Add missing semantic tokens only when reuse or theme behavior warrants them.

### CSS-04: Component scope relies on naming discipline alone

Evidence: shared SCSS emits global selectors and is loaded on all public pages. There are many generic nouns (`.card`, `.input`, `.list`, `.nav`, `.table`) mixed with domain selectors. Uno utilities and component styles coexist without a cascade-layer declaration.

Impact: selector collisions and override ordering become more likely as WIP features land; fixes tend toward higher specificity or `!important`.

Recommendation: introduce cascade layers such as `reset`, `tokens`, `base`, `components`, and `utilities`, with an explicit order compatible with Uno. Use a consistent component namespace/BEM convention for new global components. Prefer island-local stylesheet imports for island-specific selectors.

### CSS-05: `!important` usage should be categorized and constrained

Evidence: `!important` appears in reduced-motion/accessibility rules (appropriate), but also in buttons, chips, password input, loader, image uploader, and toast rules.

Impact: legitimate accessibility overrides are mixed with cascade workarounds, making future overrides unpredictable.

Recommendation: allowlist `!important` for accessibility/user-preference and third-party override boundaries. Refactor other cases through layer ordering, selector ownership, or state attributes. Require a short comment for remaining exceptions.

### CSS-06: Global asset resolution and browser-target maintenance are not verified cleanly

Evidence: the production build warned that `/_/embed-assets/img/headerImage.png` could not be resolved at build time and that Browserslist data was eight months old.

Impact: runtime-only asset conventions are not distinguished from accidental broken references, and generated prefixes/compatibility choices can grow stale.

Recommendation: encode runtime public URLs via a documented helper/alias or explicitly suppress only known external URLs. Add a periodic dependency/browser-data update job and a build policy that treats unknown resolution warnings as failures.

## Low

### CSS-07: The component barrel is a manual registry

Evidence: every component stylesheet must be added to `components/index.scss`, and ordering has no documented rationale.

Impact: new files can silently be unused, deleted files can leave imports, and cascade behavior depends on incidental list order.

Recommendation: as styles move to owning entries, keep the remaining shared barrel alphabetized or grouped by documented layer. A simple CI check can ensure naming/order and flag unreferenced SCSS files.

## Existing strengths worth preserving

- `theme.css` provides coherent light/dark semantic custom properties.
- `uno.config.ts` maps utilities to the same semantic tokens rather than establishing a second palette.
- Shared controls such as buttons, cards, inputs, alerts, tabs, and pagination have recognizable component homes.
- Reduced-motion rules exist in the common stylesheet.
- SCSS uses the module-era `@use` syntax rather than deprecated `@import`.
