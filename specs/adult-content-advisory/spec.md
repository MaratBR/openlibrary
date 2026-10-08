# Feature: Consistent adult-content and warning preferences

Status: classification service and editor contract agreed on 2026-10-08.
Censored-tag behavior is implemented according to [the newer agreed slice](censored-tags.md),
including direct-page blurring and explicit-tag search exceptions. The broader
adult-setting/interstitial rollout below remains draft.

## Agreed implementation slice

- Adult classification is owned by an injectable `AdultContentService` in the
  application layer, allowing later replacement without changing callers.
- The default policy classifies NC-17 or any assigned adult tag as adult.
  R and unrated books are not adult unless an adult tag is assigned.
- Public and manager details use the same service for their derived `adult`
  field. Direct-page warnings consume that result without repeating tag rules.
- The editor offers age rating only. It does not submit a writable adult field
  or report an adult-switch mismatch after a successful save.
- Legacy non-null `isAdult` submissions receive HTTP 400 before saving anything,
  with an instruction to omit the field. Omitted/null values are accepted.
- Discovery filtering, random selection, preference migration, and interstitials
  remain subsequent work; this slice does not change those mechanisms.

## Problem and outcome

[The source audit](../../ADULT_CONTENT_ADVISORY.md) finds persisted moderation
preferences without enforcement, inconsistent random-book filtering, a browser-wide
warning bypass, chapter URLs without warnings, and an Adult editor switch whose
submitted value is ignored. Readers should control adults-only content and
individual warnings through one predictable policy across public entry points.
Authors should have one authoritative audience classification.

## Scope

Implement account preferences, anonymous defaults, selectable warning descriptors,
consistent SQL/OpenSearch exclusions, direct-book and chapter interstitials,
intentional handling of saved books, and migration of existing settings. Fix the
book-editor classification contract. Include equivalent public APIs, covers,
summaries, previews, counts, and caches. Apply the policy to existing discovery
surfaces; future recommendations or feeds must consume the same policy.

The first release uses existing rating labels and warning tags. A replacement
rating enum, per-descriptor blur/warn actions, recommendation algorithms, legal
age verification, and a general moderation-portal redesign are outside scope.
Preferences never grant access to banned, private, or otherwise unavailable works.

## Proposed behavior

The classification boundary is agreed above; the remaining defaults are proposals:

- Adults-only means NC-17 or an assigned descriptor requiring adults-only.
  R means mature; non-adult warnings do not raise the audience classification.
- Anonymous visitors hide adults-only content and have no selected warnings.
  Signed-in account preferences take precedence over browser preferences.
- Matching any selected warning excludes a book from discovery. Adults-only
  exclusion and warning exclusion operate independently.
- Saved/library items retain a placeholder explaining the preference match.
  Covers, summaries, and chapter content remain obscured until a deliberate view.
- Direct URLs show reasons and offer Go back and View this book once. An exception
  covers that book and its chapters in the current session; it does not change
  preferences, reveal discovery results, or bypass availability restrictions.

## Scenarios and acceptance criteria

- AC-1: An author edits one age classification. The editor exposes no independently
  writable adult boolean. Changing ratings saves successfully and the returned
  derived classification matches the saved rating. Legacy adult-field submissions
  receive documented behavior rather than silently appearing to save a switch.
- AC-2: Readers can save Hide adults-only and a searchable selection of warning
  tags. Reloading restores the selection. Empty selection is valid. Invalid modes,
  unknown IDs, and ineligible descriptors produce a useful validation error;
  failed saves preserve the last persisted settings and entered choices.
- AC-3: Every rating, including unrated, follows the agreed boundary. Suicide or
  self-harm alone does not make a book adults-only. An adults-only descriptor
  raises the effective classification and exposes a moderation conflict when
  the authored rating is lower. All matching reasons are retained.
- AC-4: Search, random selection, home/discovery lists, author lists, public
  collections, and equivalent APIs exclude matching books. Excluding multiple
  descriptors means matching any one is sufficient. Counts and pagination use
  the same constraints; no eligible random book returns the existing empty-state
  behavior rather than relaxing preferences. Manual exclusions remain effective.
- AC-5: A blocked direct book, chapter, or preview request returns an interstitial
  with translated matching reasons. The response does not embed blocked covers,
  summaries, chapter text, or equivalent API payloads. Go back has a safe fallback
  when there is no usable previous page. Spoiler-marked reasons require a deliberate
  reveal so the warning itself does not unexpectedly reveal plot details.
- AC-6: View this book once allows only the selected book and its chapters within
  the agreed lifetime. Other blocked books still show an interstitial. The action
  does not update account settings or discovery constraints. Exceptions expire
  and cannot cross accounts after logout/sign-in.
- AC-7: Saved books retain their identity through a minimal placeholder and show
  why their content is hidden. Reader history and membership are preserved;
  deliberate viewing follows AC-5 and AC-6.
- AC-8: Anonymous defaults and account precedence behave as agreed. Preference
  changes take effect on the next request without revealing another reader's
  results through shared caches. A preference-loading failure never silently
  grants broader visibility; the response offers a retry or service error.
- AC-9: Existing preferences migrate according to an explicit mapping, with
  ambiguous entries reported rather than silently discarded. The legacy
  `view_adult` cookie no longer bypasses the new policy. No migration silently
  opts a reader into broader content visibility.
- AC-10: Renames preserve selected descriptors; merges resolve to canonical IDs;
  deletion follows an agreed, visible recovery path. Full int64 IDs round-trip
  through browser APIs without precision loss. Tag classification and relationship
  edits update policy results and affected indexes/caches consistently.
- AC-11: Author and moderator management access remains available under existing
  authorization. Management privileges do not change public discovery behavior.
  Public exceptions never override bans, private visibility, or chapter permissions.
- AC-12: Settings and interstitials work with keyboard navigation, visible focus,
  both themes, and mobile widths. Direct-link protection works without JavaScript;
  accessible form submissions can perform deliberate viewing.

## Open questions requiring agreement

1. Resolved: NC-17 plus assigned adult tags, behind an injectable service.
   Unrated books are not adult unless an adult tag is assigned.
2. Adopt descriptor-driven classification escalation? Narrow existing `is_adult`
   flags in place or migrate them to a warning-only `requires_adults_only` field?
   How should adult genre/topic tags and synonyms migrate?
3. Adopt anonymous defaults above? Is local anonymous customization required in
   this release, and how long should its preference storage persist?
4. Define “once” precisely: browser session or shorter expiry, and whether tabs
   share exceptions. Confirm book-plus-chapters scope and account isolation.
5. Approve saved-item placeholders and what minimal identity they may expose.
   Confirm spoiler-reason disclosure and the safe Go back destination.
6. Map `show_adult_content`, `none`/`hide`/`censor`, and existing string tag entries
   without losing reader intent. Decide recovery for unresolved/deleted tags.
7. Define legacy API compatibility for `isAdult`/`adult` and existing clients.
   Confirm the rollout consistency strategy while indexes are rebuilt.

## Verification

Use focused policy tests for AC-3, AC-6, AC-8, and AC-11; settings/editor API
regressions for AC-1 and AC-2; SQL/OpenSearch integration checks for AC-4 and
AC-10; HTTP response inspection for AC-5 and AC-7; migration fixtures for AC-9;
and browser checks for AC-2, AC-6, and AC-12. Record results and unavailable
checks during implementation. The classification service and editor contract are implemented; broader criteria
remain pending. See verification.md for checks and limitations.
