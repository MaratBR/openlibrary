# Server entry point

This guide applies to `cmd/server`. Follow the repository guide as well.

- Keep application startup, root middleware, and server-level configuration in
  this directory; business rules belong in `internal/app`.
- `olproto/search.pb.go` is generated from `olproto/search.proto`. Change the
  proto and regenerate it instead of editing the generated Go file.
