# Admin SPA refresh

Authorized by the user's requests on 2026-10-04; work on feature/admin-1.

Replace the existing admin pages with a hash-routed React island following
Book Manager, with cream surfaces, green actions, serif headings and redesigned
admin cards. Preserve existing domain operations and server authorization.

## Acceptance criteria

- AC-1: /admin opens a useful home screen; users, user editing, books lookup,
  tags, tag details/editing and debug actions navigate without document reloads.
  Existing admin page URLs redirect to their equivalent hash routes. Missing
  settings functionality is not advertised; unknown routes show recovery.
- AC-2: A compact masthead and responsive navigation replace the fixed sidebar.
  At 360, 768 and 1280px, pages have no document overflow; wide tables scroll
  within their cards. Both existing themes use project tokens. Admin cards have
  consistent rounded borders, padding, section headers and action groups.
- AC-3: User search, role filters and pagination retain query state. Editing
  preserves custom/empty gender, biography and role; password reset is optional,
  cancellable and only saved on explicit submission. Tag search, adult/parent
  filters, pagination, synonym selection, category and flags remain functional.
  Unchecked flags and empty descriptions can be submitted.
- AC-4: Loading, empty, failed, saving and success states are distinct. Failed
  mutations preserve input. Duplicate submissions are disabled. Leaving dirty
  forms asks before discarding. Back/forward and refresh preserve route state.
- AC-5: All admin API reads/mutations require the existing admin authorization
  and mutations retain CSRF protection. Invalid IDs/payloads fail without
  calling services. Client data is validated; IDs avoid JS integer precision loss.
- AC-6: Labels, focus, keyboard controls and translated visible copy are present.
  No fake delete/create/settings actions or invented statistics. Books retain
  their existing ID lookup; debug reindex only runs after explicit confirmation.

No database changes, new settings, book search, tag creation/deletion, or new
moderation operations are included.
