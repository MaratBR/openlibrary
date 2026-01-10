package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/gofrs/uuid"
)

var (
	ErrTypeBookNotFound = apperror.AppErrors.NewType("book404", apperror.ErrTraitEntityNotFound)
	ErrTypeBookPrivated = apperror.AppErrors.NewType("access_denied_book_private", apperror.ErrTraitForbidden)
)

type BookDetailsDto struct {
	ID                  int64                 `json:"id,string"`
	Name                string                `json:"name"`
	AgeRating           AgeRating             `json:"ageRating"`
	IsAdult             bool                  `json:"adult"`
	Tags                []DefinedTagDto       `json:"tags"`
	Words               int                   `json:"words"`
	WordsPerChapter     int                   `json:"wordsPerChapter"`
	CreatedAt           time.Time             `json:"createdAt"`
	Collections         []BookCollectionDto   `json:"collections"`
	Author              BookDetailsAuthorDto  `json:"author"`
	Permissions         BookUserPermissions   `json:"permissions"`
	Summary             string                `json:"summary"`
	Notifications       []GenericNotification `json:"notifications,omitempty"`
	Cover               string                `json:"cover"`
	Rating              Nullable[float64]     `json:"rating"`
	Votes               int32                 `json:"votes"`
	Reviews             int32                 `json:"reviews"`
	IsPubliclyAvailable bool                  `json:"isPubliclyAvailable"`
	Slug                string                `json:"slug"`
	FirstChapterID      Nullable[int64]       `json:"firstChapterId"`
}

type BookAdultWarning struct {
	IsBookAdult bool
	HasAdultTag bool
}

func (w BookAdultWarning) ShouldShowWarning() bool {
	return w.IsBookAdult || w.HasAdultTag
}

func (d BookDetailsDto) GetAdultWarning() (warnData BookAdultWarning) {
	warnData.IsBookAdult = d.IsAdult

	for _, t := range d.Tags {
		if t.IsAdult {
			warnData.HasAdultTag = true
			break
		}
	}

	return
}

type GetBookQuery struct {
	ID          int64
	ActorUserID uuid.NullUUID
}

type GetBookChaptersQuery struct {
	ID int64
}

type ChapterListDto struct {
	ID        int64     `json:"id,string"`
	Order     int       `json:"order"`
	Name      string    `json:"name"`
	Words     int       `json:"words"`
	CreatedAt time.Time `json:"createdAt"`
	Summary   string    `json:"summary"`
}

type BookDetailsAuthorDto struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type BookUserPermissions struct {
	CanEdit bool `json:"canEdit"`
}

type BookListDto struct {
	ID              int64     `json:"id,string"`
	Name            string    `json:"name"`
	CreatedAt       time.Time `json:"createdAt"`
	AgeRating       AgeRating `json:"ageRating"`
	Words           int       `json:"words"`
	WordsPerChapter int       `json:"wordsPerChapter"`
	Chapters        int       `json:"chapters"`
	Cover           string    `json:"cover"`
	IsPinned        bool      `json:"isPinned"`
}

type GetUserBooksQuery struct {
	UserID      uuid.UUID
	ActorUserID Nullable[uuid.UUID]
	Page        uint32
	PageSize    uint32
	SearchQuery string
}

type SearchUserBooksQuery struct {
	UserID      uuid.UUID
	ActorUserID Nullable[uuid.UUID]
	Limit       int
	Offset      int
}

type SearchUserBooksResult struct {
}

type GetPinnedUserBooksQuery struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

type GetPinnedUserBooksResult struct {
	Books   []BookListDto `json:"books"`
	HasMore bool          `json:"hasMore"`
}

type ChapterDto struct {
	ID              int64                        `json:"id,string"`
	Name            string                       `json:"name"`
	Words           int32                        `json:"words"`
	Content         string                       `json:"content"`
	IsAdultOverride bool                         `json:"isAdultOverride"`
	CreatedAt       time.Time                    `json:"createdAt"`
	Order           int32                        `json:"order"`
	Summary         string                       `json:"summary"`
	NextChapter     Nullable[ChapterNextPrevDto] `json:"nextChapter"`
	PrevChapter     Nullable[ChapterNextPrevDto] `json:"prevChapter"`
	BookID          int64                        `json:"bookId"`
	CommentsCount   int64                        `json:"commentsCount"`
}

type ChapterNextPrevDto struct {
	ID    int64  `json:"id,string"`
	Name  string `json:"name"`
	Order int32  `json:"order"`
}

type BookCollectionDto struct {
	ID       int64  `json:"id,string"`
	Name     string `json:"name"`
	Position int    `json:"pos"`
	Size     int    `json:"size"`
}

type GetBookChapterQuery struct {
	BookID      int64
	ChapterID   int64
	ActorUserID uuid.NullUUID
}

type GetBookChapterResult struct {
	Chapter ChapterDto
}

type BookService interface {
	GetBookDetails(ctx context.Context, query GetBookQuery) (BookDetailsDto, error)
	GetBookChapters(ctx context.Context, query GetBookChaptersQuery) ([]ChapterListDto, error)
	GetBookChapter(ctx context.Context, query GetBookChapterQuery) (GetBookChapterResult, error)
	GetRandomBookID(ctx context.Context) (Nullable[int64], error)
	GetPinnedBooks(ctx context.Context, input GetPinnedUserBooksQuery) (GetPinnedUserBooksResult, error)
	GetBooksById(ctx context.Context, ids []int64) ([]BookListDto, error)
}
