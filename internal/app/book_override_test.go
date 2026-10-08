package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"testing"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5"
	"uuid"
)

type overrideRoleDB struct {
	store.DBTX
	role  store.UserRole
	calls int
}

func (db *overrideRoleDB) QueryRow(context.Context, string, ...interface{}) pgx.Row {
	db.calls++
	return overrideRoleRow{db.role}
}

type overrideRoleRow struct{ role store.UserRole }

func (r overrideRoleRow) Scan(dest ...interface{}) error {
	*dest[6].(*store.UserRole) = r.role
	return nil
}
func TestBookOverrideRequiresCurrentAdminRole(t *testing.T) {
	for _, role := range []store.UserRole{store.UserRoleAdmin, store.UserRoleSystem, store.UserRoleModerator, store.UserRoleUser} {
		db := &overrideRoleDB{role: role}
		service := &bookService{queries: store.New(db)}
		allowed, err := service.authorizeBookOverride(context.Background(), Value(uuid.NewV4()), true)
		want := role == store.UserRoleAdmin || role == store.UserRoleSystem
		if allowed != want || (err == nil) != want || db.calls != 1 {
			t.Fatalf("role %s: allowed=%v err=%v calls=%d", role, allowed, err, db.calls)
		}
	}
	service := &bookService{}
	if allowed, err := service.authorizeBookOverride(context.Background(), Null[uuid.UUID](), true); allowed || err == nil {
		t.Fatal("anonymous override allowed")
	}
	if allowed, err := service.authorizeBookOverride(context.Background(), Null[uuid.UUID](), false); allowed || err != nil {
		t.Fatal("ordinary reads should not require admin")
	}
}
func TestBookVisibilityReasons(t *testing.T) {
	if reasons := bookVisibilityReasons(true, false, false, false, false); len(reasons) != 0 {
		t.Fatal(reasons)
	}
	want := []string{"private", "banned", "shadowBanned", "deleted", "removed"}
	if got := bookVisibilityReasons(false, true, true, true, true); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v", got)
	}
}

type overrideReadDB struct {
	store.DBTX
	book            store.Book_GetRow
	role            store.UserRole
	hiddenChapter   bool
	chapterOverride bool
}
type overrideScanRow struct{ scan func(...interface{}) error }

func (r overrideScanRow) Scan(dest ...interface{}) error { return r.scan(dest...) }
func (db *overrideReadDB) QueryRow(_ context.Context, sql string, args ...interface{}) pgx.Row {
	switch {
	case strings.HasPrefix(sql, "-- name: User_Get "):
		return overrideRoleRow{db.role}
	case strings.HasPrefix(sql, "-- name: Book_GetFirstChapterID "):
		return overrideScanRow{func(...interface{}) error { return store.ErrNoRows }}
	case strings.HasPrefix(sql, "-- name: GetBookChapterWithDetails "):
		db.chapterOverride = args[2].(bool)
		return overrideScanRow{func(dest ...interface{}) error {
			*dest[0].(*int64) = 2
			*dest[2].(*int64) = 1
			*dest[3].(*string) = "chapter content"
			*dest[11].(*bool) = !db.hiddenChapter
			return nil
		}}
	default:
		return overrideScanRow{func(dest ...interface{}) error {
			*dest[0].(*int64) = 1
			*dest[1].(*string) = "Book"
			*dest[4].(*pgtype.UUID) = uuidDomainToDb(uuid.MustParse("416f17a0-1740-4e61-99e0-e9e5309c8030"))
			*dest[6].(*store.AgeRating) = store.AgeRatingPG
			*dest[7].(*bool) = db.book.IsPubliclyVisible
			*dest[8].(*bool) = db.book.IsBanned
			*dest[9].(*bool) = db.book.IsTrashed
			*dest[20].(*bool) = db.book.IsPermRemoved
			*dest[21].(*bool) = db.book.IsShadowBanned
			return nil
		}}
	}
}

type overrideEmptyRows struct{ pgx.Rows }

func (overrideEmptyRows) Close()     {}
func (overrideEmptyRows) Next() bool { return false }
func (overrideEmptyRows) Err() error { return nil }
func (db *overrideReadDB) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	return overrideEmptyRows{}, nil
}

type overrideTags struct{ TagsService }

func (overrideTags) GetTagsByIds(context.Context, []int64) ([]DefinedTagDto, error) { return nil, nil }
func TestBookRestrictionsAndOverride(t *testing.T) {
	for _, book := range []store.Book_GetRow{
		{IsPubliclyVisible: false},
		{IsPubliclyVisible: true, IsBanned: true},
		{IsPubliclyVisible: true, IsShadowBanned: true},
		{IsPubliclyVisible: true, IsTrashed: true},
		{IsPubliclyVisible: true, IsPermRemoved: true},
	} {
		db := &overrideReadDB{book: book, role: store.UserRoleAdmin}
		service := &bookService{queries: store.New(db), tagsService: overrideTags{}, adultContentService: NewAdultContentService(), log: zap.NewNop().Sugar()}
		actor := Value(uuid.NewV4())
		if _, err := service.GetBookDetails(context.Background(), GetBookQuery{ID: 1, ActorUserID: actor}); err == nil {
			t.Fatalf("restricted book allowed without override: %+v", book)
		}
		dto, err := service.GetBookDetails(context.Background(), GetBookQuery{ID: 1, ActorUserID: actor, AdminOverride: true})
		if err != nil || !dto.AdminOverride || len(dto.VisibilityReasons) == 0 {
			t.Fatalf("override failed: %+v %v", dto, err)
		}
		db.role = store.UserRoleModerator
		if _, err = service.GetBookDetails(context.Background(), GetBookQuery{ID: 1, ActorUserID: actor, AdminOverride: true}); err == nil {
			t.Fatal("moderator bypassed restriction")
		}
	}
	db := &overrideReadDB{book: store.Book_GetRow{IsPubliclyVisible: true}, role: store.UserRoleAdmin}
	service := &bookService{queries: store.New(db), tagsService: overrideTags{}, adultContentService: NewAdultContentService(), log: zap.NewNop().Sugar()}
	dto, err := service.GetBookDetails(context.Background(), GetBookQuery{ID: 1, ActorUserID: Value(uuid.NewV4()), AdminOverride: true})
	if err != nil || len(dto.VisibilityReasons) > 0 {
		t.Fatalf("public book changed: %+v %v", dto, err)
	}
}
func TestHiddenChapterRequiresOverride(t *testing.T) {
	db := &overrideReadDB{book: store.Book_GetRow{IsPubliclyVisible: true}, role: store.UserRoleAdmin, hiddenChapter: true}
	service := &bookService{queries: store.New(db), tagsService: overrideTags{}, adultContentService: NewAdultContentService(), log: zap.NewNop().Sugar()}
	actor := Value(uuid.NewV4())
	query := GetBookChapterQuery{BookID: 1, ChapterID: 2, ActorUserID: actor}
	if _, err := service.GetBookChapter(context.Background(), query); err == nil {
		t.Fatal("hidden chapter allowed")
	}
	query.AdminOverride = true
	result, err := service.GetBookChapter(context.Background(), query)
	if err != nil || !reflect.DeepEqual(result.VisibilityReasons, []string{"hiddenChapter"}) || !db.chapterOverride {
		t.Fatalf("hidden chapter override failed: %+v %v", result, err)
	}
	db.book.IsBanned = true
	result, err = service.GetBookChapter(context.Background(), query)
	if err != nil || !reflect.DeepEqual(result.VisibilityReasons, []string{"banned", "hiddenChapter"}) {
		t.Fatalf("lost book reasons: %+v %v", result, err)
	}
	db.role = store.UserRoleUser
	if _, err := service.GetBookChapter(context.Background(), query); err == nil {
		t.Fatal("ordinary user overrode hidden chapter")
	}
}
