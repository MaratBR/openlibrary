# Application package boundaries

Keep a feature's contracts, validation, persistence adapters, tests, and Fx
registration together when they can depend on other packages without importing
the root `app` package. Choose boundaries by ownership and dependencies rather
than by filename or file count.

## Extracted features

| Package | Owns | Reason for isolation |
| --- | --- | --- |
| `internal/app/upload` | Upload configuration, MinIO client/buckets, public URLs, upload URIs | Storage and URI behavior have no dependency on books or users. |
| `internal/app/tags` | Tag contracts, categories, validation, queries, predefined tag import | Tag management and import share one vocabulary and persistence adapter, independent of book workflows. |
| `internal/app/reader` | Reader preferences, font catalog, validation, preference persistence | Reader settings have their own rules and storage; they do not depend on authoring or book font policy. |
| `internal/app/userdata` | Per-user JSON data, key whitelist, size limit, persistence | JSON data storage has a distinct contract and validation policy, independent of account authentication. |

Each feature exports an `FXModule`, included by `internal/app/fx.go`. Reader
preferences and user data need `store.DBTX`, not a transaction-capable database.
Existing `app` constructors retain their signatures for current callers.

`*_compat.go` files contain type/constant aliases, references to existing error
values, and forwarding constructors. Existing handlers and services can continue
using `app`; new code within an extracted feature uses its package directly.
Keep business logic and validation in the owning package, not in these bridges.

## Shared dependencies

- `internal/app/value` owns `Nullable[T]` and its JSON behavior. The root
  `app.Nullable[T]` is a type alias, so values and service signatures remain
  interchangeable. Existing `app.Null`, `app.Value`, and `app.NullableFromPtr`
  forward to the same implementation.
- `internal/app/dal` owns the shared UUID and timestamp database conversions.
  Root conversion helpers forward to those implementations. Other helpers remain
  at their current scope until another feature needs them.
- `internal/app/idgen` owns the existing process-wide numeric ID generator.
  Both tag creation and `app.GenID` use the same counter; do not create a second
  generator when extracting a service.

These support packages exist to remove dependencies on root `app`, not as
destinations for unrelated utilities.

## Code retained in root `app`

The remaining files were reviewed as the following dependency groups. Their
existing coupling makes a larger extraction require additional design choices.
Keep them together until those choices have a concrete maintenance benefit.

| Group and files | Why retained |
| --- | --- |
| Books: `book*.go`, `age_rating.go`, `adult_content.go`, `draft_publishing_service.go` | Reading, authoring, drafts, visibility, covers, statistics, and reindexing share book types, helpers, and services. Splitting these by filename would scatter one workflow. |
| Discovery: `search*.go`, `opensearch.go` | Search uses book presentation, ratings, tag aggregation, request content preferences, user avatars, and the root cache flag. The empty backend placeholder provides no independent boundary. |
| Personal library: `collections*.go`, `reading_list*.go` | These reuse book covers, ratings, tags, and content filtering. Extracting them first would require a shared book presentation contract. |
| Accounts: `auth*.go`, `signup*.go`, `session*.go`, `user.go`, `user_impl.go`, `password_hashing.go`, `ip_location.go` | Authentication, account creation, sessions, roles, banned-user errors, and password hashing form a connected cluster. User settings also participate in request content filtering. |
| Moderation: `moderation*.go`, `reports*.go`, `login_history*.go` | Authorization, user roles, book actions, report decisions, login locations, and audit records share contracts. Report persistence also uses moderation detail types, so separating reports alone would introduce a cycle. |
| Discussion: `comments*.go`, `reviews*.go` | Comments and reviews currently share account presentation; reviews also use user services and book recalculation. A separate package needs a deliberate dependency contract rather than a copy of root helpers. |
| Content filtering: `censored_tags*.go` | Request-scoped filtering connects account settings, book/search results, and collections. It is not tag CRUD or predefined tag import. |
| Composition and support: `fx.go`, `cron.go`, `app_version.go`, `global_feature_flags.go`, `db_abstraction.go`, `id.go`, `id_generator.go`, `utils*.go`, `slug.go`, `otel.go`, `notifications.go`, `tags_aggregator.go`, `content_test.go` | Wiring, cross-feature helpers, tiny shared values, and integration tests have no useful new feature boundary. The tag aggregator assembles book/collection results and stays with those consumers. |

Tests follow the rules and implementations they exercise. Reader, upload URI,
and user-data tests live in their feature packages; cross-feature tests remain
in root `app`.

## Extending the layout

Before extracting another group, identify its owned behavior and the contracts
it needs from its neighbors. Avoid an import back to root `app`, copying shared
rules, or introducing a broad shared package solely to make files movable.
Keep Fx composition in `internal/app/fx.go` and verify service wiring after
changing registrations.

Run `go test ./internal/app/...`, then `go test ./...` to check existing callers.
When moving types marked for `go2tsdef`, verify the generated TypeScript shapes
remain compatible; their source-file comments will change.
