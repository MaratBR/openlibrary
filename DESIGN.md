---
name: OpenLibrary dashboards
description: A restrained publishing production desk for library management.
colors:
  primary: "#123f31"
  primary-dark: "#a5c4b5"
  primary-foreground: "#fffdf8"
  primary-foreground-dark: "#17221d"
  surface: "#fffdf8"
  surface-dark: "#1d1f1d"
  foreground: "#2e342f"
  foreground-dark: "#d9dad5"
  secondary: "#f1eee6"
  secondary-dark: "#292b28"
  secondary-foreground: "#73756f"
  secondary-foreground-dark: "#aaaca6"
  accent: "#e5ede8"
  accent-dark: "#303a34"
  border: "#ddd7ca"
  border-dark: "#3a3d39"
  border-strong: "#c8c0b1"
  border-strong-dark: "#50544e"
typography:
  headline:
    fontFamily: "Poppins, Inter, ui-sans-serif, system-ui, sans-serif"
    fontSize: "2rem"
    fontWeight: 600
    lineHeight: 1.25
    letterSpacing: "-0.03em"
  title:
    fontFamily: "Poppins, Inter, ui-sans-serif, system-ui, sans-serif"
    fontSize: "1.15rem"
    fontWeight: 600
    lineHeight: 1.4
    letterSpacing: "-0.02em"
  body:
    fontFamily: "Poppins, Inter, ui-sans-serif, system-ui, sans-serif"
    lineHeight: 1.55
  label:
    fontFamily: "Poppins, Inter, ui-sans-serif, system-ui, sans-serif"
    fontSize: "0.83rem"
    fontWeight: 500
    lineHeight: 1.5
rounded:
  flat: "0"
  cover: "3px"
  navigation: "6px"
  control: "0.5rem"
  container: "12px"
spacing:
  compact: "4px"
  small: "8px"
  inline: "12px"
  medium: "16px"
  section: "24px"
  header: "28px"
  large: "40px"
components:
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.primary-foreground}"
    rounded: "{rounded.control}"
  navigation:
    typography: "{typography.label}"
    rounded: "{rounded.navigation}"
    padding: "11px 12px"
  record:
    rounded: "{rounded.flat}"
    padding: "20px 12px"
  container:
    rounded: "{rounded.container}"
---

# Design System: OpenLibrary dashboards

## Overview

**Creative North Star: "Publishing production desk"**

A publishing production desk: calm, compact work surfaces with persistent navigation and clearly divided records. Forest navigation, warm neutral surfaces, and restrained green states support sustained scanning in light and dark themes.

This system applies only to administration, book management, and moderation. The public library and reading interface keep their incumbent visual system. Shared theme tokens and existing input controls are dependencies, not redesigned assets.

**Key Characteristics:**

- Persistent navigation with stable state geometry.
- Flat, ruled records and cover-led book rows.
- Sans-serif dashboard hierarchy in both themes.
- Existing controls, copy, and workflows preserved.

## Colors

The frontmatter records the shared light values and their `-dark` counterparts. Source of truth: `web/frontend/src/common/style/theme.css`; dashboard-only mixtures live in `dashboard-shell.scss`.

### Primary

Forest green anchors actions and navigation; its pale sage dark-theme counterpart maintains contrast. Soft green accent surfaces identify supporting states rather than entire page backgrounds.

### Neutral

Warm paper surfaces and subdued ink form the light canvas; charcoal surfaces and pale ink form the dark canvas. Secondary surfaces support table headers and cover wells. Border tokens divide records.

The light rail mixes primary (92%) with foreground; the dark rail mixes surface (85%) with primary. Rail ink is primary foreground in light mode and foreground in dark mode. The canvas mixes surface (96%) with foreground. Rail muted ink mixes rail ink (76%) with the rail.

**The Semantic Theme Rule.** Use the existing semantic CSS variables; the dashboard must adapt to the current theme without replacing global tokens.

## Typography

Dashboard headings and text use the shared sans-serif stack from the frontmatter. Page headings use the headline role, book titles the title role, and navigation the label role. Mobile page headings reduce to (1.6rem); mobile book titles reduce to (1rem). Supporting text stays subordinate, with descriptions limited to (70ch) where implemented. Tables and record lists use tabular numerals.

## Layout

Desktop navigation occupies a sticky left rail (248px), with brand above links and utilities at the bottom. Between (761px) and (1100px), the rail narrows to (216px). At (760px) and below, navigation becomes a static top region with wrapping links.

Content uses fluid horizontal padding `clamp(24px, 4vw, 64px)` and top padding (40px). Mobile content uses (20px) horizontal and (28px) top padding. Administration content caps at (1344px); book management and moderation at (1440px). Divided page headers precede ruled directories, horizontal book records, or moderation sections and tables.

## Elevation & Depth

Dashboard cards, record lists, and action bars are flat. Tonal surfaces and one-pixel rules establish grouping. Existing shared buttons and fields retain their incumbent small shadows and focus treatment; dropdown overlays retain their shared elevation. Do not remove those control affordances to enforce flatness.

## Shapes

Working records have square corners; contained sections use the container radius. Navigation and cover wells use the navigation radius; covers use the smaller cover radius. Existing control radii remain unchanged. Nested cards lose their outer frame and become divided sections.

## Components

### Buttons

Retain shared primary, outline, ghost, secondary, and destructive variants. Primary hover darkens its main color with black (12%); active darkens it (20%) and moves down (1px). Focus uses the existing theme ring. Rail utility buttons use transparent backgrounds, muted ink, and no shadows.

### Inputs / Fields

Existing fields are outside the redesign. Continue using the shared surface, foreground, border, radius, validation, disabled, read-only, and focus behavior; do not introduce dashboard-specific field styling.

### Navigation

Links keep a minimum target height (44px), rounded shape, and unchanged dimensions across states. Hover adds a subtle inverse tint; the current primary link reverses rail and ink. Recent-record sublinks remain indented, smaller, and subordinate. Rail keyboard focus is a two-pixel ink outline with (3px) offset.

### Cards / Containers

Containers have no dashboard shadow. Nested containers become transparent, square, top-divided sections. Tables use secondary header backgrounds and tabular numerals. Chips retain shared neutral, primary, and destructive treatments.

### Working records

Administration directories and moderation records are full-width ruled rows. Book records place cover, identity, metadata, and actions in aligned columns; tablet and mobile move the footer below the text while retaining cover identity. Hover uses the shared surface; focus on linked records uses a two-pixel primary outline inset (2px).

## Do's and Don'ts

### Do:

- Do use shared semantic theme variables so both themes stay coherent.
- Do distinguish record groups with rules and spacing.
- Do keep keyboard focus visible and allow long record names to wrap.
- Do retain existing input controls and public-site styling.

### Don't:

- Don't apply dashboard heading or layout overrides outside the dashboard shell.
- Don't turn working records into oversized floating tiles.
- Don't move navigation geometry when a link becomes active.
