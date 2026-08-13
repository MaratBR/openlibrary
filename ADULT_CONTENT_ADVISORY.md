# Adult content advisory

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
- A derived book `IsAdult` value, which is true for both `R` and `NC-17`.
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

`AgeRating.IsAdult()` returns true for both `R` and `NC-17`.

The database stores only `books.age_rating`; there is no independent adult
boolean on books. `BookDetailsDto.IsAdult` is derived from the age rating in the
application service. Exposing both `AgeRating` and `IsAdult` makes the model
appear to have two independent classifications even though it does not.

There is also a semantic mismatch in treating `R` as 18+. Existing translation
copy describes it as content for adults and mature teenagers, whereas `NC-17`
is explicitly adults-only. Therefore `R` and `NC-17` should not automatically
be treated as the same user preference boundary.

### Adult tags and warnings

Every defined tag has independent `is_adult` and `is_spoiler` booleans. A book
shows an adult warning when either:

- Its derived `IsAdult` value is true; or
- At least one assigned tag has `is_adult = true`.

The repository already seeds warning tags including:

- Major character death
- Graphical depiction of violence
- Suicide
- Self-harm
- Underage

These are currently not marked adult, which is appropriate: they describe
content or possible triggers rather than necessarily imposing an age
restriction. However, readers cannot yet save preferences that hide them.

The tag administration copy also notes a consistency problem: changing a tag
to adult does not automatically update books that were indexed before that
change.

### Book-page warning

The only active adult-content handling is a generic overlay on a direct book
page. It is shown for R/NC-17 books or books with an adult tag.

The overlay:

- Says only "Potential adult content" and does not identify the reason.
- Offers `Proceed` but no explicit `Go back` action.
- Sets the browser cookie `view_adult=1` after proceeding.
- Suppresses all future adult warnings on that device, not only for that book.

The cookie is not scoped to an account, content type, reason, or book. A single
proceed action therefore permanently collapses all adult-content distinctions
for that browser until the cookie is removed.

The account-level `ShowAdultContent` check in `web/public/adult_flag.go` is
commented out, so the saved preference does not control this overlay.

### Account preferences

The user table and application DTO already contain:

- `show_adult_content boolean`
- `censored_tags text[]`
- `censored_tags_mode`, with `none`, `hide`, and `censor`

The settings endpoint can load and save these values. The moderation settings
page also exposes the adult-content switch and censor mode.

However:

- The censored-tag selector is disabled with a TODO placeholder.
- Saved censored tags are not applied to content retrieval.
- `show_adult_content` is not applied to search, lists, or the book warning.
- There is no central content-visibility policy.

Thus these are persisted settings rather than implemented filtering features.

### Search and other discovery surfaces

Search supports explicit include/exclude tag IDs from the search request, but
does not incorporate the reader's account preferences. Neither the SQL search
filter nor the OpenSearch request has an age-rating constraint derived from the
user's settings.

Consequently, adult or user-censored content can still appear in search and
other lists. Opening a result only adds the generic warning overlay.

Random-book selection is a notable exception: its SQL query always excludes
both R and NC-17 books. It does this regardless of user preference and does not
exclude adult-tagged books. This creates inconsistent behavior across entry
points.

## Feature status

| Feature | Status | Behavior |
| --- | --- | --- |
| Book age rating | Implemented | Stored on every book using the six-value enum. |
| Book `IsAdult` | Derived | True for R and NC-17; not an independent database flag. |
| Adult tags | Implemented | Any tag can be marked adult. |
| Warning/trigger tags | Partially implemented | Some warnings are seeded and displayed as tags. |
| Direct-page warning | Implemented with limitations | One generic, device-wide dismissible overlay. |
| Show-adult account setting | Stored but unused | Does not affect discovery or direct pages. |
| Censored-tag preferences | Skeleton only | Storage exists, selector is disabled, policy is unused. |
| Manual search exclusion | Implemented | Callers can explicitly exclude tag IDs. |
| Hide all 18+ content | Not implemented | No consistent account-aware filtering exists. |
| Hide specific topics/triggers | Not implemented | Warning tags cannot be selected as preferences. |

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

1. Define the exact adults-only boundary and rename derived `IsAdult` usages to
   `IsAdultsOnly`.
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

- Every age-rating boundary, especially R versus adults-only.
- A non-adult warning such as suicide matching an individual preference.
- An adults-only descriptor raising the effective classification.
- Multiple matching reasons on one book.
- Search results and counts excluding hidden books consistently.
- Direct links returning an interstitial rather than silently exposing content.
- `View once` not changing global or persistent preferences.
- Anonymous defaults and signed-in preference precedence.
- Author and moderator access without weakening public discovery filtering.
- Tag rename/synonym behavior when preferences use stable IDs.

## Conclusion

The main problem is not a lack of fields; it is that age rating, adult tagging,
warnings, preferences, and discovery filtering do not form one coherent policy.

Use age rating only for audience suitability, use warning descriptors for
specific content and triggers, and derive adults-only status rather than asking
authors or moderators to maintain multiple overlapping flags. A single
`Hide adults-only` option plus selected hidden warnings satisfies the immediate
requirements while leaving room for more nuanced blur/warn behavior later.
