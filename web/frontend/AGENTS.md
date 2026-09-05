# Frontend guide

This guide applies to `web/frontend`. Follow the repository guide as well.

## Code organization

- Organize new frontend work by feature. Use `src/features` for functional
  features and `src/islands` for server-rendered React island entry points; do
  not broadly reorganize existing code just to fit the convention.
- Put APIs shared by multiple features in `src/features/api`; keep
  feature-specific APIs under the owning feature. `src/public.api` is deprecated:
  do not add APIs there, and migrate touched callers to the appropriate API.
- Use the `@/` alias for `src` imports when it improves readability.
- Keep simple, page-local behavior in Alpine `x-data` by the markup. Put shared
  or substantial Alpine behavior in `src/alpinejs`. Use React islands for
  stateful, reusable, or complex UI; shared React controls belong in
  `src/components`.
- Keep Zod schemas aligned with every Go DTO field that callers consume; Zod
  object parsing strips undeclared fields.
- Follow `docs/development/EFFECT.md` for Effect services, layers, React program
  execution, logging, and Schema.

## Copy and styling

- Use `i18nExtractKeys` or `i18nExtractKeysByPrefix` as appropriate. Add
  user-facing English to `translations/en.toml`; do not hard-code visible English
  in translated components.
- Inspect `src/common/style/components` before adding controls. Prefer existing
  `.btn`, `.BtnGroup`, `.card`, and related variants.
- Use UnoCSS for local layout, spacing, responsiveness, and one-off details.
  Add scoped SCSS only for a coherent reused pattern or to simplify utility-heavy
  markup.
- Use theme tokens instead of literal colors. Prefer `secondary-foreground` for
  subdued text, `secondary` for subdued surfaces, and `foreground/5` for neutral
  interaction feedback.
- Use complete static Uno utility names; dynamically constructed classes may not
  be scanned. Add a justified safelist entry when necessary.

## Verification

- Run `pnpm run build` after TypeScript, React, Alpine, API, or style changes.
