# Dashboard sidebars

The admin and book manager workspaces should use a standard left sidebar consistent with the site design, as requested by the user. Design details are delegated to implementation judgment.

## Acceptance criteria

- Desktop workspaces have a persistent left sidebar containing branding, route navigation, and site-return actions.
- Admin retains all existing sections, theme switching, and logout. Book Manager offers Your books and Add book; book detail routes remain under Your books.
- Navigation clearly identifies the current section and is keyboard accessible.
- Small screens retain visible, wrapping navigation above the content without page overflow.
- Both themes use existing theme tokens and site controls. Existing loading, error, mutation, and unsaved-form behavior remains intact.

## Visual refinement

Use a slim sidebar, a rounded inset content surface on desktop, compact sans-serif dashboard headings, subtle active navigation surfaces, and consistent utility icons and labels. Retain the site’s theme palette and book-title typography.

The admin frame fills its workspace column; its main content and footer are centered and limited to 1344px, keeping wide-screen content readable.
