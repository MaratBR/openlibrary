---
version: 1
slug: "account-settings"
primary_target: "web/public/templates/account_settings.templ"
related_targets: ["web/frontend/src/common/style/components/account-settings.scss", "/account/settings/about", "/account/settings/privacy", "/account/settings/moderation", "/account/settings/customization"]
---

# Account settings

## Direction contract

An ordinary extension of the existing theme: warm paper, a neutral secondary
sidebar, serif headings, and ruled fields make personal settings distinct and
readable. The user approved the layout, rejected the first flat treatment as
boring, then rejected forest navigation as too green. Neutral navigation is the
final direction; these decisions apply only to account settings.

The settings masthead uses the existing title font at (3.5rem), with a (2rem)
panel heading. A thin brass panel-header rule and the selected link's small brass
icon provide restrained accents. The selected link uses a surface background,
strong border, and foreground ink; the sidebar uses secondary. All colors remain
semantic theme variables in light and dark modes. Shared inputs and switches keep
their existing styling and behavior.

## Layout and behavior

Desktop uses a sticky (232px) rail, fluid panel, and (48px) gap. At (1000px) and
below, the rail becomes (208px), the gap becomes (32px), and field contents stack.
At (640px) and below, the page becomes one column and navigation becomes a
two-column grid; masthead and panel headings reduce to (2.25rem) and (1.6rem).
Panels and sidebar retain (12px) corners; fields use thin dividing rules.

About, Privacy, Content moderation, and Customization retain their routes,
controls, tag search, theme selection, and Alpine autosave PATCH logic. Labels
and descriptions connect to controls through Alpine `x-id`; current navigation
uses `aria-current="page"`, save feedback uses a polite live region, and links
retain visible keyboard focus. Reduced motion disables navigation transitions.

## Evidence and drift

Implementation sources are the targets above. Review captures cover all four
tabs at (1440px) and (390px), in both themes:
`.impeccable/review/settings/{about,privacy,moderation,customization}-{1440,390}-{light,dark}.png`.
The implementation review reports a passing frontend build and Go package tests,
plus successful autosave responses (HTTP 200) on every tab with unchanged payloads.
Finish review returned Ship with no material findings across all (16) captures;
whitespace checks passed.
An existing global flash-initializer error was observed outside this scope.

This is a scoped expression of the shared neutral palette and serif type, with
brass already supplied by the theme. Its neutral navigation intentionally differs
from the dashboard's forest rail. Global `DESIGN.md` and
`.impeccable/design.json` remain unchanged; do not promote these settings-specific
layout or navigation choices into global rules.
