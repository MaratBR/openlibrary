package app

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/MaratBR/openlibrary/internal/store"
)

var ErrInvalidCensoredTags = errors.New("invalid censored tags")

// ContentPreferences is request-scoped. Management calls have no preferences.
// Matches are recorded for server-rendered blur controls, including fragments.
type ContentPreferences struct {
	Mode     CensorMode
	TagIDs   []int64
	Families map[int64]int64
	mu       sync.RWMutex
	censored map[int64]bool
}

type contentPreferencesKey struct{}

func WithContentPreferences(ctx context.Context, prefs *ContentPreferences) context.Context {
	return context.WithValue(ctx, contentPreferencesKey{}, prefs)
}

func contentPreferences(ctx context.Context) *ContentPreferences {
	prefs, _ := ctx.Value(contentPreferencesKey{}).(*ContentPreferences)
	return prefs
}

func HiddenTagIDs(ctx context.Context) []int64 {
	if prefs := contentPreferences(ctx); prefs != nil && prefs.Mode == CensorMode(store.CensorModeHide) {
		return append([]int64{}, prefs.TagIDs...)
	}
	return []int64{}
}

// ApplyBookTagPreferences matches any selected tag (including expanded synonyms).
func ApplyBookTagPreferences(ctx context.Context, bookID int64, tags []int64) {
	prefs := contentPreferences(ctx)
	if prefs == nil || prefs.Mode == CensorMode(store.CensorModeNone) {
		return
	}
	for _, tag := range tags {
		for _, selected := range prefs.TagIDs {
			if tag != selected {
				continue
			}
			if prefs.Mode == CensorMode(store.CensorModeCensor) || prefs.Mode == CensorMode(store.CensorModeHide) {
				prefs.mu.Lock()
				if prefs.censored == nil {
					prefs.censored = make(map[int64]bool)
				}
				prefs.censored[bookID] = true
				prefs.mu.Unlock()
			}
			return
		}
	}
	return
}

func IsBookCensored(ctx context.Context, bookID int64) bool {
	prefs := contentPreferences(ctx)
	if prefs == nil {
		return false
	}
	prefs.mu.RLock()
	defer prefs.mu.RUnlock()
	return prefs.censored[bookID]
}

// ResolveCensoredTags supports legacy names while storing stable decimal IDs.
// Unknown/deleted selections are reported so they can be removed explicitly.
func ResolveCensoredTags(ctx context.Context, queries *store.Queries, values []string) ([]int64, error) {
	if len(values) == 0 {
		return []int64{}, nil
	}
	if len(values) > 100 {
		return nil, fmt.Errorf("%w: at most 100 tags", ErrInvalidCensoredTags)
	}
	ids := make([]int64, 0, len(values))
	names := make([]string, 0)
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			return nil, fmt.Errorf("%w: empty tag", ErrInvalidCensoredTags)
		}
		id, err := strconv.ParseInt(value, 10, 64)
		if err == nil && id > 0 {
			ids = append(ids, id)
		} else {
			names = append(names, value)
		}
	}
	byID, err := queries.GetTagsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	byName, err := queries.GetTagsByName(ctx, names)
	if err != nil {
		return nil, err
	}
	resolved := map[string]int64{}
	for _, tag := range append(byID, byName...) {
		canonical := tag.ID
		if tag.SynonymOf.Valid {
			canonical = tag.SynonymOf.Int64
		}
		resolved[strconv.FormatInt(tag.ID, 10)] = canonical
		resolved[tag.Name] = canonical
	}
	result := make([]int64, 0, len(values))
	seen := map[int64]bool{}
	for _, value := range values {
		id, ok := resolved[strings.TrimSpace(value)]
		if !ok {
			return nil, fmt.Errorf("%w: unknown tag %q", ErrInvalidCensoredTags, value)
		}
		if !seen[id] {
			result = append(result, id)
			seen[id] = true
		}
	}
	return result, nil
}

func ResolveContentPreferences(ctx context.Context, queries *store.Queries, settings *UserModerationSettings) (*ContentPreferences, error) {
	switch settings.CensoredTagsMode {
	case "none", "hide", "censor":
	default:
		return nil, fmt.Errorf("%w: invalid mode", ErrInvalidCensoredTags)
	}
	prefs := &ContentPreferences{Mode: settings.CensoredTagsMode, TagIDs: []int64{}}
	if settings.CensoredTagsMode == CensorMode(store.CensorModeNone) || len(settings.CensoredTags) == 0 {
		return prefs, nil
	}
	ids, err := ResolveCensoredTags(ctx, queries, settings.CensoredTags)
	if err != nil {
		return nil, err
	}
	families, err := queries.GetCensoredTagFamilyIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	prefs.Families = make(map[int64]int64, len(families))
	for _, tag := range families {
		prefs.TagIDs = append(prefs.TagIDs, tag.ID)
		prefs.Families[tag.ID] = tag.CanonicalID
	}
	return prefs, nil
}

// An explicit include is an intentional request; those matches remain blurred.
func HiddenSearchTagIDs(ctx context.Context, included []int64) []int64 {
	prefs := contentPreferences(ctx)
	if prefs == nil {
		return []int64{}
	}
	revealedFamilies := map[int64]bool{}
	for _, id := range included {
		if family, ok := prefs.Families[id]; ok {
			revealedFamilies[family] = true
		}
	}
	result := make([]int64, 0)
	for _, id := range HiddenTagIDs(ctx) {
		if !revealedFamilies[prefs.Families[id]] {
			result = append(result, id)
		}
	}
	return result
}

func IsTagCensored(ctx context.Context, id int64) bool {
	prefs := contentPreferences(ctx)
	if prefs == nil || prefs.Mode == "none" {
		return false
	}
	for _, selected := range prefs.TagIDs {
		if selected == id {
			return true
		}
	}
	return false
}
