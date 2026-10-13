# Plan

Share the insertion-widget filter between sidebar and slash providers. Register an atomic Tiptap bookPreview node serialized as ol-widget with name=book-preview and JSON data containing a decimal book ID. Add an inline configuration form and reuse the existing server preview template and endpoint for new editor insertions.

Expand stored references into validated containers, then render their book cards in web/public using the request’s actor and existing BookService access rules. Cache repeated references only within the request. Include rendered fragments in editor island data, so saved widgets mount immediately with no preview fetch. The node view uses the supplied ProseMirror view rather than editor.view during construction. Remove client reader hydration; readers receive complete HTML, including unavailable states and links. Do not recursively expand references within preview summaries.

Share existing preview card styles through common styles. Keep widgets compact with a themed tint, small cover, limited summary/tags, and one row of actions. Memoize editor options in EditorIframe so sidebar rerenders preserve unsaved changes. Add translations, focused markup/reference/server-rendering regressions, and check frontend build plus scoped Go tests. Verify saved editor startup and reader output with JavaScript disabled using Playwright response fixtures.
