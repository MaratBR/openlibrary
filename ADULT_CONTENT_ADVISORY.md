# Adult content advisory

Reviewed on 2026-10-04; classification/editor implementation updated 2026-10-08.
NC-17 plus assigned adult tags is now the agreed default behind an injectable
application service. Broader filtering and preference recommendations remain
proposals. The review is based on source inspection and automated tests, not a
running deployment or production data.

## Goal

Give readers control over content they do not want to encounter while keeping
the model and settings reasonably flexible and simple.

Readers should be able to:

1. Hide all adults-only (18+) content.
2. Hide content containing selected topics, such as nudity, violence, or mature
   themes.
3. Hide individual trigger warnings, such as suicide or self-harm, without
   treating every such topic as adults-only.

## Executive summary

OpenLibrary already has most of the underlying data structures, but content
filtering is not wired together into a working user-facing policy.

Three overlapping concepts currently exist:

- A book age rating (`?`, `G`, `PG`, `PG-13`, `R`, or `NC-17`).
- A derived book `IsAdult` value from the classification service (NC-17 or adult tags).
- An independent `is_adult` flag on any defined tag.

This is confusing because an age classification, a content topic, and an
adults-only restriction answer different questions. For example, suicide can
be an important trigger without making a work adults-only. Likewise, an adult
work should not depend on whether somebody happened to mark one of its tags as
adult.

The recommended model has two orthogonal dimensions:

1. One authoritative audience/age classification for the book.
2. Structured content descriptors or warnings for topics present in the book.

The user's first version of preferences can then remain small:

- Hide adults-only content.
- Hide books matching selected content descriptors.

## Current implementation

### Age ratings

`internal/app/age_rating.go` defines the following ratings:

- `?`
- `G`
- `PG`
- `PG-13`
- `R`
- `NC-17`

`AdultContentService.IsAdult(rating, tags)` owns the classification policy.
Its default implementation returns true for NC-17 or any assigned tag marked
adult. R and unrated books are not adult without an adult tag. Fx injects the
service into public and manager book services; a replacement provider can
customize the rule later. There is no independent adult boolean stored on books.

### Book-editor classification contract

The editor exposes age rating and derives adult status on the server. The
independent Adult switch and its false save-error check have been removed.
The update request no longer submits `isAdult`; legacy non-null submissions
receive HTTP 400 before saving, instructing callers to omit it. The response
retains the read-only `adult` property, including tag-driven classification.

### Adult tags and warnings

Every defined tag has independent `is_adult` and `is_spoiler` booleans. A book
shows an adult warning when either:

- Its service-derived `IsAdult` value is true (NC-17 or an adult tag).

The repository defines warning-tag seeds in
[`oldata/tags/warnings.hjson`](oldata/tags/warnings.hjson), including:

- Major character death
- Graphical depiction of violence
- Suicide
- Self-harm
- Underage

These seed definitions are not marked adult; deployed tags may have been
changed by administrators. They describe content or possible triggers rather
than necessarily imposing an age restriction. The seed comment for Underage
explicitly leaves its classification uncertain, so its meaning needs review. Readers can now select these or other existing tags in moderation settings.

The tag administration copy claims existing books are not automatically marked
adult in search when a tag changes. Treat that as UI copy, not evidence of an
implemented search restriction: the current index stores the book rating and
cached parent tag IDs, without a separate adult-tag classification. Direct-page
warnings inspect loaded tag flags. A future indexed descriptor policy will need
explicit invalidation when tags or their relationships change.

### Book-page warning

The only active adult-content handling is a generic overlay on a direct book
page. It is shown for books classified adult by the service: NC-17 or adult-tagged.

The overlay:

- Says only "Potential adult content" and does not identify the reason.
- Offers `Proceed` but no explicit `Go back` action.
- Sets the browser cookie `view_adult=1` after proceeding.
- Suppresses subsequent book-page adult warnings in that browser for the site,
  not only for that book.

The cookie is not scoped to an account, content type, reason, or book. A single
proceed action therefore collapses all adult-content distinctions for that
browser until the cookie expires or is removed. The shared
[`setCookie`](web/frontend/src/common/cookies.ts) helper defaults to a 20-year
expiry and path `/`; browser retention may be shorter.

The account-level `ShowAdultContent` check in `web/public/adult_flag.go` is
commented out, so the saved preference does not control this overlay.

The overlay is a visual warning, not a server-side interstitial: the normal
book page, including its cover and summary, is still rendered beneath it.
[`routes_chapter.go`](web/public/routes_chapter.go) renders the chapter reader
without calling the adult-warning check, so a direct chapter URL does not
require dismissing this warning.

### Account preferences

The user table and application DTO already contain:

- `show_adult_content boolean`
- `censored_tags text[]`
- `censored_tags_mode`, with `none`, `hide`, and `censor`

The settings endpoint can load and save these values and validates the censor
mode. `CensoredTags` is currently `[]string`; there is no implemented contract
resolving those strings into canonical descriptor IDs. The moderation settings
page also exposes the adult-content switch and censor mode.

However:

- The censored-tag selector now searches existing tags and stores decimal IDs.
- Saved censored tags now drive public filtering and blur/reveal controls.
- `show_adult_content` is not applied to search, lists, or the book warning.
- Request-scoped tag preferences centralize matching and administrator bypasses.

Tag filtering is implemented; the separate show-adult setting remains unused.

### Search and other discovery surfaces

Search supports explicit include/exclude tag IDs and incorporates censored-tag
preferences before pagination/counts. Explicitly including a banned tag reveals
its family in blurred form, while other exclusions remain effective. Banned tags
are omitted from the public selector except in authorized administrator mode. Neither the SQL search
filter nor the OpenSearch request has an age-rating constraint derived from the
user's settings.

Adult content can still appear unless censored-tag preferences independently
exclude it. Direct-page adult warnings remain the existing generic overlay.

Random-book selection is a notable exception: its SQL query always excludes
both R and NC-17 books. It does this regardless of user preference and does not
exclude adult-tagged books. This creates inconsistent behavior across entry
points.

## Feature status

| Feature | Status | Behavior |
| --- | --- | --- |
| Book age rating | Implemented | Stored on every book using the six-value enum. |
| Book `IsAdult` | Derived through service | NC-17 or an assigned adult tag; replaceable policy. |
| Book-editor Adult switch | Removed | Age rating is authoritative; legacy non-null `isAdult` receives HTTP 400 before saving. |
| Adult tags | Implemented | Any tag can be marked adult. |
| Warning/trigger tags | Partially implemented | Some warnings are seeded and displayed as tags. |
| Direct-page warning | Implemented with limitations | Generic browser-wide dismissible overlay; underlying content is rendered. |
| Direct chapter warning | Not implemented | Chapter reader does not use the book-page adult-warning check. |
| Show-adult account setting | Stored but unused | Does not affect discovery or direct pages. |
| Censored-tag preferences | Implemented | Searchable selector; none/hide/blur behavior across public book surfaces, with authorized admin bypasses. |
| Manual search exclusion | Implemented | Callers can explicitly exclude tag IDs. |
| Hide all 18+ content | Not implemented | No consistent account-aware filtering exists. |
| Hide specific topics/triggers | Implemented | Hide removes matching books from discovery; explicit banned-tag searches and direct pages remain blurred. |

## Source map

| Concern | Sources |
| --- | --- |
| Rating and derived adult status | [adult_content.go](internal/app/adult_content.go), [age_rating.go](internal/app/age_rating.go), [book_impl.go](internal/app/book_impl.go), [book.go](internal/app/book.go) |
| Editor update contract | [BookEdit.tsx](web/frontend/src/islands/bookmanager/books/BookEdit.tsx), [routes_bm_api.go](web/public/routes_bm_api.go), [book_manager_impl.go](internal/app/book_manager_impl.go) |
| Warning and bypass | [routes_book.go](web/public/routes_book.go), [adult_flag.go](web/public/adult_flag.go), [book.templ](web/public/templates/book.templ), [cookies.ts](web/frontend/src/common/cookies.ts) |
| Stored preferences and settings UI | [query.user.sql](internal/store/query.user.sql), [user.go](internal/app/user.go), [settings.go](web/public/account/settings.go), [account_settings.templ](web/public/templates/account_settings.templ) |
| Search and index | [search_impl.go](internal/app/search_impl.go), [book_search.go](internal/store/book_search.go), [search.go](internal/elasticstore/search.go), [book_search_reindex_impl.go](internal/app/book_search_reindex_impl.go) |
| Random selection | [query.book.sql](internal/store/query.book.sql) (`GetRandomPublicBookIDs`) |

`SiteConfig.AdultWebsite` also exists and defaults to false in
[openlibrary.toml](openlibrary.toml), but has no policy consumer in the current
Go code. It does not enable account filtering or a site-wide age gate.

## Decisions needed before implementation

- Resolved: NC-17 plus assigned adult tags is the default classification,
  abstracted behind `AdultContentService` for later customization.
- Define handling of unrated books and descriptors that contradict a rating.
- Agree on anonymous defaults, saved-item behavior, and the lifetime and scope
  of a per-book exception, including chapter access.
- Decide how existing `show_adult_content`, `censored_tags`, censor modes, and
  `view_adult` cookies migrate. Do not silently discard saved preferences or
  carry the global bypass into the new policy.
- Define descriptor eligibility, synonym canonicalization, and behavior when a
  selected tag is merged or deleted. Transport int64 IDs as decimal strings to
  JavaScript clients to avoid precision loss.

Record agreed observable requirements and the implementation plan in
[`specs/adult-content-advisory/spec.md`](specs/adult-content-advisory/spec.md)
(the accompanying plan and tasks remain provisional) following
[spec-driven development](docs/development/spec-driven-development.md) before
implementing substantial behavior changes.

## Recommended domain model

### Separate audience classification from content descriptors

A book should have one age/audience classification and zero or more content
descriptors:

```text
Book
  age_rating: everyone | teen | mature | adults_only | unknown
  content_descriptors: [nudity, graphic_violence, suicide, ...]
```

This separates:

- **Who the work is suitable for**, represented by the age rating.
- **What the work contains**, represented by descriptors/warnings.

Do not retain a separately authored book-level adult boolean. If API clients
need a convenience value, expose a clearly derived `isAdultsOnly` property.

### Age-rating boundary

The existing US film-rating labels are not an especially natural or global fit
for written works. A product-owned scale would be clearer:

- Everyone
- Teen
- Mature
- Adults only
- Unrated

If changing the enum is too disruptive initially, keep it for compatibility
but define only `NC-17` as unambiguously adults-only. Treat `R` as mature unless
the product explicitly decides that all R-rated books must be hidden by the
18+ preference.

This policy decision should be made once and encoded centrally rather than
inferred independently by handlers.

### Replace vague adult tags

The existing warning tags can serve as content descriptors in the first
iteration. The generic `is_adult` flag should be deprecated or narrowed to an
unambiguous meaning such as `requires_adults_only`.

For example:

```text
defined_tags
  tag_type = warning
  requires_adults_only = true | false
```

Alternatively, future classification may use a small sensitivity level:

```text
sensitivity_level = general | mature | adults_only
```

Warning tags are eligible for individual user filtering. A descriptor marked
`requires_adults_only` also raises the book's effective classification to
adults-only. If that conflicts with a lower author-selected age rating, the
book should be treated safely as adults-only and flagged for moderation.

## Recommended user experience

Start with two controls:

```text
Adult content
  [x] Hide adults-only content

Content I don't want to see
  [ ] Sexual content
  [ ] Nudity
  [ ] Graphic violence
  [ ] Suicide
  [ ] Self-harm
  [ ] Other warning tags...
```

The corresponding initial settings model can be:

```go
type ContentPreferences struct {
	HideAdultsOnly      bool
	HiddenDescriptorIDs []int64
}
```

Use stable tag IDs rather than names. Names can be translated, renamed, or
merged through synonyms.

For simplicity, do not expose a global `none`/`hide`/`censor` mode in the first
complete version. A predictable default policy is sufficient:

- Remove matching books from search, recommendations, feeds, random books, and
  other discovery surfaces.
- Preserve deliberately saved/library items if silently removing them would
  confuse users, but obscure their metadata and explain why.
- On a direct URL, show an interstitial listing the exact matching reasons.
- Provide `Go back` and `View this book once` actions.
- Do not turn `View once` into a device-wide permanent bypass.

A later version can support per-descriptor actions such as hide, blur, or warn
if real user needs justify the additional complexity.

### Anonymous readers

Anonymous visitors need a conservative default and can store preferences
locally. Account preferences should take precedence after sign-in. The precise
default is a product decision, but `HideAdultsOnly = true` is consistent with
the stated goal.

The interface should distinguish "hidden by your preferences" from content
that is unavailable due to site policy or law. User preferences are not an
access-control mechanism.

## Central visibility policy

Filtering should be decided in `internal/app`, not reimplemented by individual
HTTP handlers. A central policy could accept the book classification, content
descriptors, and resolved user preferences, then return an outcome such as:

```go
type ContentVisibilityDecision struct {
	Visibility ContentVisibility // visible, hidden, interstitial
	Reasons    []ContentReason
}
```

The same resolved constraints must be applied to:

- Search and search counts/facets
- Home and discovery feeds
- Recommendations
- Random-book selection
- Public lists and collections
- Reading lists and saved items, with an intentional saved-item policy
- Direct book pages
- Chapter access and previews
- API responses used by equivalent frontend views

Filtering only after fetching results is insufficient: it creates short pages,
incorrect counts, and possible content leaks through summaries or covers. Age
and descriptor constraints should be pushed into SQL/OpenSearch where possible,
while the application policy remains the source of truth.

## Suggested implementation order

1. Classification service and editor contract implemented. Agree on remaining
   filtering, exception, and preference migration decisions.
2. Introduce a central content-visibility policy with focused unit tests.
3. Enable a searchable warning-tag selector in account settings and store tag
   IDs rather than names.
4. Apply resolved age and descriptor exclusions to OpenSearch and SQL-backed
   discovery surfaces.
5. Make random-book behavior use the same policy instead of unconditional
   R/NC-17 exclusion.
6. Replace the global cookie overlay with a reason-specific direct-link
   interstitial and a per-book/session `View once` exception.
7. Apply the policy to covers, summaries, chapters, collections, reading lists,
   and APIs so alternate paths cannot leak hidden content.
8. Add moderation validation for conflicts between ratings and descriptors.

## Tests to add

At minimum, cover:

- Editor saves across the R/NC-17 boundary and ignored-switch mismatches,
  including whether other fields were persisted before the error.
- Every age-rating boundary, especially R versus adults-only and unrated books.
- A non-adult warning such as suicide matching an individual preference.
- An adults-only descriptor raising the effective classification.
- Multiple matching reasons on one book.
- Search results and counts excluding hidden books consistently.
- Direct book and chapter links returning an interstitial without embedding
  blocked content, covers, or summaries in the response.
- `View once` not changing global or persistent preferences.
- Anonymous defaults and signed-in preference precedence.
- Author and moderator access without weakening public discovery filtering.
- Tag rename, merge, deletion, and synonym behavior when preferences use stable
  IDs, including int64 ID round trips through the frontend.
- Tag classification changes invalidating affected search documents and caches.
- Migration of existing preference values and the long-lived global cookie.

## Conclusion

The main problem is not a lack of fields; it is that age rating, adult tagging,
warnings, preferences, and discovery filtering do not form one coherent policy.

Use age rating only for audience suitability, use warning descriptors for
specific content and triggers, and derive adults-only status rather than asking
authors or moderators to maintain multiple overlapping flags. A single
`Hide adults-only` option plus selected hidden warnings satisfies the immediate
requirements while leaving room for more nuanced blur/warn behavior later.

## Censored-tag implementation update (2026-10-08)

See [agreed behavior and verification](specs/adult-content-advisory/censored-tags.md).
`none` leaves public content unchanged. `hide` excludes matches from search,
random selection, home books, profiles, collections and library queries before
pagination where applicable. `hide` and `censor` blur direct book/chapter content;
`censor` also retains blurred discovery results. Revealing is scoped to the
current document and book, including fragments, and does not change settings.
Synonym families match together; legacy names resolve to stable string IDs on save.
Admin and management routes bypass preferences. Public `admin.link=1` and
`admin.override=1` bypass them only for authenticated admin/system users.
