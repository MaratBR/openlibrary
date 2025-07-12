package app

import (
	"context"
	"io"
	"time"

	"github.com/gofrs/uuid"
	"github.com/joomcode/errorx"
)

var (
	ErrTypeBookSanitizationFailed = AppErrors.NewType("content_sanitization_failed")
	ErrTypeChaptersReorder        = AppErrors.NewType("chapters_reorder")
	ErrDraftNotFound              = AppErrors.NewType("draft_not_found", errorx.NotFound(), ErrTraitEntityNotFound).New("draft not found")
	ErrTypeChapterDoesNotExist    = AppErrors.NewType("chapter_not_found", ErrTraitEntityNotFound)
)

type CreateBookCommand struct {
	Name              string
	UserID            uuid.UUID
	AgeRating         AgeRating
	Tags              []int64
	Summary           string
	IsPubliclyVisible bool
}

type UpdateBookCommand struct {
	BookID            int64
	UserID            uuid.UUID
	Name              string
	AgeRating         AgeRating
	Tags              []int64
	Summary           string
	IsPubliclyVisible bool
}

type ManagerGetBookQuery struct {
	ActorUserID uuid.UUID
	BookID      int64
}

type ManagerBookDetailsDto struct {
	ID                int64                `json:"id,string"`
	Name              string               `json:"name"`
	AgeRating         AgeRating            `json:"ageRating"`
	IsAdult           bool                 `json:"adult"`
	Tags              []DefinedTagDto      `json:"tags"`
	Words             int                  `json:"words"`
	WordsPerChapter   int                  `json:"wordsPerChapter"`
	CreatedAt         time.Time            `json:"createdAt"`
	Collections       []BookCollectionDto  `json:"collections"`
	Chapters          []BookChapterDto     `json:"chapters"`
	Author            BookDetailsAuthorDto `json:"author"`
	Summary           string               `json:"summary"`
	IsPubliclyVisible bool                 `json:"isPubliclyVisible"`
	IsBanned          bool                 `json:"isBanned"`
	Cover             string               `json:"cover"`
}

type ManagerGetBookResult struct {
	Book ManagerBookDetailsDto
}

type ManagerAuthorBookDto struct {
	ID                int64               `json:"id,string"`
	Name              string              `json:"name"`
	CreatedAt         time.Time           `json:"createdAt"`
	AgeRating         AgeRating           `json:"ageRating"`
	Tags              []DefinedTagDto     `json:"tags"`
	Words             int                 `json:"words"`
	WordsPerChapter   int                 `json:"wordsPerChapter"`
	Chapters          int                 `json:"chapters"`
	Collections       []BookCollectionDto `json:"collections"`
	IsPubliclyVisible bool                `json:"isPubliclyVisible"`
	IsBanned          bool                `json:"isBanned"`
	Summary           string              `json:"summary"`
	Cover             string              `json:"cover"`
}

type GetUserBooksResult struct {
	Books      []ManagerAuthorBookDto
	TotalPages uint32
	PageSize   uint32
	Page       uint32
}

type CreateBookChapterCommand struct {
	BookID          int64
	Name            string
	Content         string
	IsAdultOverride bool
	Summary         string
}

type CreateBookChapterResult struct {
	ID int64
}

type ReorderChaptersCommand struct {
	UserID     uuid.UUID
	BookID     int64
	ChapterIDs []int64
}

type ManagerBookChapterDto struct {
	ID              int64     `json:"id,string"`
	Name            string    `json:"name"`
	CreatedAt       time.Time `json:"createdAt"`
	Words           int       `json:"words"`
	Summary         string    `json:"summary"`
	Order           int32     `json:"order"`
	IsAdultOverride bool      `json:"isAdultOverride"`
}

type ManagerGetBookChaptersQuery struct {
	BookID int64
	UserID uuid.UUID
}

type ManagerGetBookChapterResult struct {
	Chapters []ManagerBookChapterDto
}

type ManagerGetChapterQuery struct {
	BookID    int64
	ChapterID int64
	UserID    uuid.UUID
}

type ManagerBookChapterDetailsDto struct {
	ID                int64     `json:"id,string"`
	Name              string    `json:"name"`
	CreatedAt         time.Time `json:"createdAt"`
	Words             int       `json:"words"`
	Summary           string    `json:"summary"`
	Order             int32     `json:"order"`
	IsAdultOverride   bool      `json:"isAdultOverride"`
	Content           string    `json:"content"`
	IsPubliclyVisible bool      `json:"isPubliclyVisible"`
}

type ManagerGetChapterResult struct {
	Chapter ManagerBookChapterDetailsDto
}

type UploadBookCoverCommand struct {
	UserID uuid.UUID
	BookID int64
	File   io.Reader
}

type UploadBookCoverResult struct {
	URL string
}

type UpdateBookChaptersOrders struct {
	BookID     int64
	ChapterIDs []int64
}

type GetDraftQuery struct {
	UserID    uuid.UUID
	DraftID   int64
	ChapterID int64
	BookID    int64
}

type DraftDto struct {
	ID          int64     `json:"id,string"`
	ChapterName string    `json:"chapterName"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ChapterID   int64     `json:"chapterId,string"`
	CreatedBy   struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	} `json:"createdBy"`
	Book struct {
		ID   int64  `json:"id,string"`
		Name string `json:"name"`
	} `json:"book"`
}

type UpdateDraftCommand struct {
	Content         string
	Summary         string
	Name            string
	IsAdultOverride bool
	DraftID         int64
	ChapterID       int64
	BookID          int64
	UserID          uuid.UUID
}

type UpdateDraftContentCommand struct {
	Content   string
	DraftID   int64
	ChapterID int64
	BookID    int64
	UserID    uuid.UUID
}

type DeleteDraftCommand struct {
	DraftID int64
	UserID  uuid.UUID
}

type PublishDraftCommand struct {
	DraftID int64
	UserID  uuid.UUID
}

type GetLatestDraftQuery struct {
	ChapterID int64
	UserID    uuid.UUID
}

type CreateDraftCommand struct {
	ChapterID int64
	UserID    uuid.UUID
}

type GetBookDraftsQuery struct {
	Page int
}

type BookManagerService interface {
	GetUserBooks(ctx context.Context, input GetUserBooksQuery) (GetUserBooksResult, error)
	GetBook(ctx context.Context, query ManagerGetBookQuery) (ManagerGetBookResult, error)
	CreateBook(ctx context.Context, input CreateBookCommand) (int64, error)
	UpdateBook(ctx context.Context, input UpdateBookCommand) error
	UploadBookCover(ctx context.Context, input UploadBookCoverCommand) (UploadBookCoverResult, error)

	UpdateBookChaptersOrder(ctx context.Context, input UpdateBookChaptersOrders) error
	CreateBookChapter(ctx context.Context, input CreateBookChapterCommand) (CreateBookChapterResult, error)
	ReorderChapters(ctx context.Context, input ReorderChaptersCommand) error
	GetBookChapters(ctx context.Context, query ManagerGetBookChaptersQuery) (ManagerGetBookChapterResult, error)
	GetChapter(ctx context.Context, query ManagerGetChapterQuery) (ManagerGetChapterResult, error)

	GetDraft(ctx context.Context, query GetDraftQuery) (DraftDto, error)
	// GetBookDrafts(ctx context.Context, query GetBookDraftsQuery)
	UpdateDraft(ctx context.Context, cmd UpdateDraftCommand) error
	UpdateDraftContent(ctx context.Context, cmd UpdateDraftContentCommand) error
	DeleteDraft(ctx context.Context, cmd DeleteDraftCommand) error
	PublishDraft(ctx context.Context, cmd PublishDraftCommand) error
	CreateDraft(ctx context.Context, cmd CreateDraftCommand) (int64, error)
	GetLatestDraft(ctx context.Context, cmd GetLatestDraftQuery) (Nullable[int64], error)
}
