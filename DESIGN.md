---
name: OpenLibrary
description: A publishing production desk, reading gallery, and book detail literary dossier.
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
  book-display:
    fontFamily: "Literata, Iowan Old Style, Palatino Linotype, Georgia, serif"
    fontSize: "clamp(2.5rem, 4.2vw, 3.75rem)"
    fontWeight: 500
    lineHeight: 1.14
    letterSpacing: "-0.035em"
  book-summary:
    fontFamily: "Poppins, Inter, ui-sans-serif, system-ui, sans-serif"
    fontSize: "1.0625rem"
    lineHeight: 1.8
  library-display:
    fontFamily: "Literata, Iowan Old Style, Palatino Linotype, Georgia, serif"
    fontSize: "clamp(2.25rem, 4vw, 3.5rem)"
    fontWeight: 500
    lineHeight: 1.15
    letterSpacing: "-0.03em"
  library-section:
    fontFamily: "Literata, Iowan Old Style, Palatino Linotype, Georgia, serif"
    fontSize: "1.6rem"
    fontWeight: 500
  library-book-title:
    fontFamily: "Poppins, Inter, ui-sans-serif, system-ui, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 500
    lineHeight: 1.5
  library-tab:
    fontFamily: "Poppins, Inter, ui-sans-serif, system-ui, sans-serif"
    fontSize: "0.95rem"
    fontWeight: 500
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
  book-cover: "2px 5px 5px 2px"
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
  book-column: "48px"
components:
  book-masthead:
    typography: "{typography.book-display}"
    textColor: "{colors.foreground}"
  book-cover:
    backgroundColor: "{colors.surface}"
    rounded: "{rounded.book-cover}"
    width: "200px"
  book-chapter:
    rounded: "{rounded.flat}"
    padding: "18px 8px"
  book-summary:
    typography: "{typography.book-summary}"
  library-header:
    backgroundColor: "{colors.surface}"
    typography: "{typography.library-display}"
  library-tab:
    textColor: "{colors.secondary-foreground}"
    typography: "{typography.library-tab}"
    padding: "8px 0 16px"
  library-tab-current:
    textColor: "{colors.primary}"
    typography: "{typography.library-tab}"
    padding: "8px 0 16px"
  library-cover:
    backgroundColor: "{colors.secondary}"
    rounded: "{rounded.cover}"
  library-book-title:
    textColor: "{colors.foreground}"
    typography: "{typography.library-book-title}"
  library-create:
    textColor: "{colors.primary}"
    padding: "20px 0"
  library-collection:
    typography: "{typography.library-section}"
    rounded: "{rounded.flat}"
    padding: "0 0 32px"
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

# Design System: OpenLibrary

## Overview

**Creative North Star: "Publishing production desk"**

A publishing production desk: calm, compact work surfaces with persistent navigation and clearly divided records. Forest navigation, warm neutral surfaces, and restrained green states support sustained scanning in light and dark themes.

The publishing production desk applies to administration, book management, and moderation. The scoped **Reading gallery** applies to `/library`, `/library/archive`, and `/library/collections`: serif headings, real covers, readable titles, and flat sections help readers recognize and open saved books. The scoped **Literary dossier** applies to public book detail routes (`/book/:slug`): a cover shelf and serif masthead lead into flat metadata and sequential chapters. Other public pages and the reading interface retain their incumbent styling. Shared theme tokens and existing input controls are dependencies, not redesigned assets.

**Key Characteristics:**

- Persistent navigation with stable state geometry.
- Flat, ruled records and cover-led book rows.
- Sans-serif dashboard hierarchy in both themes.
- Existing controls, copy, and workflows preserved.
- Reader gallery with serif headings and fully wrapping book titles.
- Book detail dossier with a cover shelf, direct reading action, and sequential chapter records.

## Colors

The frontmatter records the shared light values and their `-dark` counterparts. Source of truth: `web/frontend/src/common/style/theme.css`; dashboard-only mixtures live in `dashboard-shell.scss`.

### Primary

Forest green anchors actions and navigation; its pale sage dark-theme counterpart maintains contrast. Soft green accent surfaces identify supporting states rather than entire page backgrounds.

### Neutral

Warm paper surfaces and subdued ink form the light canvas; charcoal surfaces and pale ink form the dark canvas. Secondary surfaces support table headers and cover wells. Border tokens divide records.

The light rail mixes primary (92%) with foreground; the dark rail mixes surface (85%) with primary. Rail ink is primary foreground in light mode and foreground in dark mode. The canvas mixes surface (96%) with foreground. Rail muted ink mixes rail ink (76%) with the rail.

**The Semantic Theme Rule.** Use the existing semantic CSS variables; the dashboard must adapt to the current theme without replacing global tokens.

The reading gallery uses the same semantic palette in both themes: surface for its header, secondary for cover wells, primary for active tabs and title hover, and border for section rules. It introduces no new color primitives.

The literary dossier retains these semantic colors, with muted ink mixed from foreground (75%) and background. It introduces no new palette.

## Typography

Dashboard headings and text use the shared sans-serif stack from the frontmatter. Page headings use the headline role, book titles the title role, and navigation the label role. Mobile page headings reduce to (1.6rem); mobile book titles reduce to (1rem). Supporting text stays subordinate, with descriptions limited to (70ch) where implemented. Tables and record lists use tabular numerals.

Library display and section roles use the existing Literata stack. Book titles and tabs use the shared sans-serif stack; mobile tabs reduce to (0.85rem). Titles wrap completely rather than clamp.

Book detail uses the scoped book-display role for the masthead and book-summary for readable prose. At (767px) and below, the masthead uses `clamp(2rem, 7.4vw, 3rem)` and summary text reduces to (0.9375rem). Section headings use Literata (1.5rem); chapter/review tabs use (1.375rem), reducing to (1.25rem) on mobile. These intentional sizes do not alter dashboard or library typography.

## Layout

Desktop navigation occupies a sticky left rail (248px), with brand above links and utilities at the bottom. Between (761px) and (1100px), the rail narrows to (216px). At (760px) and below, navigation becomes a static top region with wrapping links.

Content uses fluid horizontal padding `clamp(24px, 4vw, 64px)` and top padding (40px). Mobile content uses (20px) horizontal and (28px) top padding. Administration content caps at (1344px); book management and moderation at (1440px). Divided page headers precede ruled directories, horizontal book records, or moderation sections and tables.

The reading gallery retains the public `OlContainer`. Header top padding is `clamp(32px, 5vw, 64px)`; body top padding is `clamp(28px, 4vw, 48px)` with (64px) below and (56px) between reading-status sections. Books use an auto-fill grid with a (150px) minimum column and (32px × 24px) gaps. At (480px) and below, two equal columns use (28px × 20px) gaps. Collections remain flat groups separated by (40px), with wrapping (100px × 150px) cover strips.

Book detail caps its grid at (1120px), with a sticky (320px) cover/action rail, (48px) column gap, and (32px) row gap. The cover is (240px × 360px), reducing to (208px) wide at intermediate widths. Book facts sit below the reading actions, so the synopsis aligns with the top of the cover. At (1000px) and below, the rail becomes (248px) with a (36px) column gap. At (767px) and below, one content column places a compact cover beside the actions; below (360px), the cover reduces further. Metadata stays a flat two-column grid; summary measure caps at (68ch). Chapters remain a single sequential column.

## Elevation & Depth

Dashboard cards, record lists, and action bars are flat. Tonal surfaces and one-pixel rules establish grouping. Existing shared buttons and fields retain their incumbent small shadows and focus treatment; dropdown overlays retain their shared elevation. Do not remove those control affordances to enforce flatness.

Library covers alone carry resting shadows (`0 4px 10px rgb(0 0 0 / 12%)`), deepening on hover (`0 8px 18px rgb(0 0 0 / 20%)`). Sections and collection groups remain flat. Cover shadow transitions use (180ms ease-out); tab transitions use (140ms ease). Reduced motion removes both transitions.

Book detail confines resting elevation to its cover: `4px 8px 12px rgb(0 0 0 / 12%), 12px 20px 30px rgb(0 0 0 / 14%)`, with an inset spine shade `inset 5px 0 8px rgb(0 0 0 / 18%)`. Metadata, summary, and chapter records stay flat. These scoped shadows are intentional book material cues.

## Shapes

Working records have square corners; contained sections use the container radius. Navigation and cover wells use the navigation radius; covers use the smaller cover radius. Existing control radii remain unchanged. Nested cards lose their outer frame and become divided sections.

Library covers use a (2:3) aspect ratio and the cover radius. The collection-create icon has a circular (44px) accent well. Gallery collection groups have square corners and bottom rules.

Book detail uses the asymmetric book-cover radius to suggest a spine; this intentional scoped shape does not replace the gallery cover or shared control radii.

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

### Reading gallery

Underlined tabs expose the current page with `aria-current="page"`; their minimum height is (48px). Book links combine covers and complete titles; hover deepens cover shadow and underlines the primary-colored title. Keyboard focus outlines the whole book with the shared ring, offset (6px). Collection headers underline on hover, and their links and cover links retain visible ring focus. The full-width create action reveals the existing labeled collection form and focuses its input. Existing preview and censorship wrappers remain attached to books.

**The Gallery Scope Rule.** Apply gallery collection geometry only to library collections; profile collections retain their incumbent compact treatment.

### Literary dossier

The masthead places author and compact rating/statistics on a wrapping byline. A sticky cover shelf holds a full-width Start reading link only when a first chapter exists; it opens that chapter directly. Existing save/status actions remain separate and retain permissions and state behavior. The compact phone shelf places cover and actions beside each other. Use existing cover assets, including the text fallback when absent.

Metadata, tags, summary, external links, chapters, and reviews retain their content. Chapter records have a minimum height (88px), sequential numbering, supporting metadata, and an arrow. Selected chapter/review tabs use primary ink and an underline; Arrow keys and Home/End move focus and selection. Linked controls use a two-pixel shared ring with (4px) offset. Chapter hover uses secondary; reading-now uses accent. Reduced motion disables chapter background transitions.

**The Dossier Scope Rule.** Apply book masthead, cover, and chapter geometry only to public book detail and its table of contents; preserve global navigation, theme, search results, and the reading interface.

## Do's and Don'ts

### Do:

- Do use shared semantic theme variables so both themes stay coherent.
- Do distinguish record groups with rules and spacing.
- Do keep keyboard focus visible and allow long record names to wrap.
- Do retain existing input controls and styling outside the approved surfaces.
- Do preserve real covers, complete book titles, and visible gallery keyboard focus.

### Don't:

- Don't apply dashboard heading or layout overrides outside the dashboard shell.
- Don't turn working records into oversized floating tiles.
- Don't move navigation geometry when a link becomes active.
- Don't apply gallery collection overrides to profile collections.
