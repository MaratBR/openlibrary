# go2tsdef

`go2tsdef` scans Go files and emits TypeScript definitions for type declarations
whose doc comment contains `// go2tsdef:generate`.

```yaml
src:
  - ./internal/**/*.go
  - ./cmd/**/*.go
dest:
  default: ./web/frontend/src/generated/types.ts
  admin: ./web/frontend/src/generated/admin-types.ts
insert_before: |
  import type { Brand } from './brand';

  type ExternalID = Brand<string, 'ExternalID'>;
types:
  time.Time: string
  github.com/MaratBR/openlibrary/internal/app:Int64String: string
```

Paths are relative to the config file. A source entry may be a directory, a file,
or a glob supporting `**`. Test files are ignored. `insert_before` accepts arbitrary
TypeScript and places it before the generated definitions. Run it with:

```sh
go run ./cmd/go2tsdef -config ./go2tsdef.yaml
```

Every generated declaration has a JSDoc comment identifying its original Go
source file. The source path is relative to the configuration file and uses `/`
separators so generated output is stable across operating systems.

Types without a destination argument are written to `dest.default`. Add a
destination name after the generate directive to route a type elsewhere:

```go
// go2tsdef:generate admin
type AdminDTO struct {
	ID int64 `json:"id,string"`
}
```

Override the entire generated definition with `go2tsdef:override_type`:

```go
// go2tsdef:generate
// go2tsdef:override_type "draft" | "published"
type Status string
```

This emits `export type Status = "draft" | "published";`.

Struct field names follow `json` tags, including `omitempty`. Override an individual
field's generated type with a `go2tsdef` struct tag:

```go
// go2tsdef:generate
type MyDTO struct {
	SomeType time.Time `go2tsdef:"number"`
}
```
