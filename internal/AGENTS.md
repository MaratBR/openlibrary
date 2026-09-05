# Internal Go code

This guide applies to `internal`. Follow the repository guide as well.

## Boundaries

- Put interfaces, DTOs, validation, state transitions, business rules, and
  application errors in `app`.
- Keep HTTP parsing, authentication context, and responses in `web/public`;
  application services must not depend on HTTP handlers.
- Persistence operations originate in `store/query.*.sql`; application services
  call generated `store.Queries` methods rather than embedding SQL elsewhere.
- Register application services in `internal/app/fx.go`.

## Domain values

- Use `app.Nullable[T]` for optional domain or JSON values when surrounding code
  does. It serializes as its contained value or `null`, never as a
  `Value`/`Valid` object.
- Add focused tests for nontrivial validation, state transitions, and regression
  fixes when practical.
