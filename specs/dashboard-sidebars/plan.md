# Plan

Use a shared SCSS dashboard frame loaded through the common stylesheet entry. Reuse existing router links, translated labels, logos, controls, and theme tokens. Keep desktop navigation sticky in its own scrollable column; use an always-visible compact navigation layout below 761px. Preserve each workspace’s route content and loading boundaries.

Verify with the frontend production build, existing admin and book manager browser fixtures, responsive geometry checks, and whitespace checks. No backend or template changes are required.

The modern visual pass uses a 232px navigation column, an inset rounded desktop content panel, compact sans-serif workspace headings, subtle card shadows, and labeled utility actions. Mobile retains an edge-to-edge content surface. Book titles retain the site’s serif typography.
