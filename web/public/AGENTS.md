# Public web handlers

This guide applies to `web/public`. Follow the repository guide as well.

## Handler responsibilities

- Controllers parse HTTP input, obtain authentication context, call application
  services, and render templates or API responses.
- Use `handler_response_api.go` and `handler_response_html.go` rather than
  duplicating response or error conventions.
- Use `internal/olhttp` helpers for URL, path, and query parsing.
- Register new public controllers in `handler.go` and constructors in
  `FXModule`.

## Templates

- Edit `.templ` sources only; never edit ignored `*_templ.go` output.
- After a template change, run `templ generate` and compile or test the affected
  Go package.
- Keep useful server-rendered, no-JavaScript content where practical. Serialize
  server data with template helpers; never construct JSON by string
  concatenation.
