# Implementation plan

1. Add authenticated JSON endpoints under /admin/api using existing user, tag
   and debug services. Use explicit transport DTOs and validated form bodies;
   preserve root CSRF middleware. Retain legacy POST handlers for compatibility.
2. Serve a single translated React island at /admin; redirect legacy GET routes.
   Keep the implementation under features/admin and use thin island/asset entry
   wrappers. Remove obsolete SSR pages and their unused controllers/islands.
   Register with the existing island loader and use React Router hash routing,
   route loaders, cancellable reads and explicit mutation state.
3. Build shared masthead, navigation, cards, forms, loading/error states and dirty
   form protection. Compose home, users/edit, books, tags/detail/edit and debug.
4. Add focused server regression tests and browser fixture checks covering routes,
   filters, forms, failure states, confirmations and responsive themes. Run Go
   tests, templ generation, frontend build and whitespace checks. Record any
   environment/baseline limitations honestly.
