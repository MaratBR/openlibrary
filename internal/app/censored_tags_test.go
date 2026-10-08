package app

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"uuid"
)

func TestTagPreferencesModesAndExplicitSearch(t *testing.T) {
	for _, mode := range []CensorMode{"none", "hide", "censor"} {
		t.Run(string(mode), func(t *testing.T) {
			prefs := &ContentPreferences{Mode: mode, TagIDs: []int64{1, 2, 3}, Families: map[int64]int64{1: 1, 2: 1, 3: 3}}
			ctx := WithContentPreferences(context.Background(), prefs)
			ApplyBookTagPreferences(ctx, 10, []int64{2})
			ApplyBookTagPreferences(ctx, 11, []int64{4})
			require.Equal(t, mode != "none", IsBookCensored(ctx, 10))
			require.False(t, IsBookCensored(ctx, 11))
			require.Equal(t, mode != "none", IsTagCensored(ctx, 2))
			if mode == "hide" {
				require.Equal(t, []int64{1, 2, 3}, HiddenSearchTagIDs(ctx, nil))
				require.Equal(t, []int64{3}, HiddenSearchTagIDs(ctx, []int64{2}))
				require.Empty(t, HiddenSearchTagIDs(ctx, []int64{1, 3}))
			} else {
				require.Empty(t, HiddenSearchTagIDs(ctx, nil))
			}
			// Requests cannot inherit another reader's matches or reveal choices.
			require.False(t, IsBookCensored(context.Background(), 10))
		})
	}
}

type censoredTagDB struct {
	store.DBTX
	tags  []store.DefinedTag
	saved *store.User_UpdateModerationSettingsParams
}

type censoredTagRows struct {
	pgx.Rows
	rows     [][]any
	position int
}

func (r *censoredTagRows) Next() bool { r.position++; return r.position <= len(r.rows) }
func (r *censoredTagRows) Close()     {}
func (r *censoredTagRows) Err() error { return nil }
func (r *censoredTagRows) Scan(dest ...any) error {
	for i, value := range r.rows[r.position-1] {
		reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(value))
	}
	return nil
}
func (db *censoredTagDB) Query(_ context.Context, query string, args ...any) (pgx.Rows, error) {
	rows := &censoredTagRows{}
	for _, tag := range db.tags {
		matches := false
		canonical := tag.ID
		if tag.SynonymOf.Valid {
			canonical = tag.SynonymOf.Int64
		}
		switch {
		case strings.HasPrefix(query, "-- name: GetTagsByIds"):
			for _, id := range args[0].([]int64) {
				matches = matches || id == tag.ID
			}
		case strings.HasPrefix(query, "-- name: GetTagsByName"):
			for _, name := range args[0].([]string) {
				matches = matches || name == tag.Name
			}
		case strings.HasPrefix(query, "-- name: GetCensoredTagFamilyIDs"):
			for _, id := range args[0].([]int64) {
				matches = matches || id == canonical
			}
			if matches {
				rows.rows = append(rows.rows, []any{tag.ID, canonical})
			}
			continue
		default:
			return nil, errors.New("unexpected query")
		}
		if matches {
			rows.rows = append(rows.rows, []any{tag.ID, tag.Name, tag.Description, tag.IsSpoiler, tag.IsAdult, tag.CreatedAt, tag.TagType, tag.SynonymOf, tag.IsDefault, tag.LowercasedName})
		}
	}
	return rows, nil
}
func (db *censoredTagDB) Exec(_ context.Context, _ string, args ...any) (pgconn.CommandTag, error) {
	db.saved = &store.User_UpdateModerationSettingsParams{
		ID: args[0].(pgtype.UUID), ShowAdultContent: args[1].(bool),
		CensoredTags: args[2].([]string), CensoredTagsMode: args[3].(store.CensorMode),
	}
	return pgconn.NewCommandTag("UPDATE 1"), nil
}

func TestCensoredTagsSaveCanonicalIDsAndPreserveSettingsOnError(t *testing.T) {
	const largeID int64 = 9223372036854775806
	db := &censoredTagDB{tags: []store.DefinedTag{
		{ID: largeID, Name: "Violence"},
		{ID: 7, Name: "Violent", SynonymOf: pgtype.Int8{Int64: largeID, Valid: true}},
	}}
	service := &userService{queries: store.New(db)}
	settings := UserModerationSettings{CensoredTagsMode: "hide", CensoredTags: []string{"Violent", "7", strconv.FormatInt(largeID, 10)}}
	require.NoError(t, service.UpdateUserModerationSettings(context.Background(), uuid.NewV4(), settings))
	require.Equal(t, []string{strconv.FormatInt(largeID, 10)}, db.saved.CensoredTags)
	prefs, err := ResolveContentPreferences(context.Background(), store.New(db), &settings)
	require.NoError(t, err)
	require.ElementsMatch(t, []int64{largeID, 7}, prefs.TagIDs)
	previous := db.saved
	for _, values := range [][]string{{"missing"}, {""}, {"9223372036854775808"}, make([]string, 101)} {
		settings.CensoredTags = values
		require.ErrorIs(t, service.UpdateUserModerationSettings(context.Background(), uuid.NewV4(), settings), ErrInvalidCensoredTags)
		require.Same(t, previous, db.saved)
	}
	settings.CensoredTagsMode = "invalid"
	require.ErrorIs(t, service.UpdateUserModerationSettings(context.Background(), uuid.NewV4(), settings), ErrInvalidCensoredTags)
	settings.CensoredTagsMode, settings.CensoredTags = "none", nil
	require.NoError(t, service.UpdateUserModerationSettings(context.Background(), uuid.NewV4(), settings))
	require.Equal(t, []string{}, db.saved.CensoredTags)
}
