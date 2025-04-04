package app

import (
	"context"
	"io"
	"time"

	"github.com/gofrs/uuid"
)

var (
	ErrTypeBookSanitizationFailed = AppErrors.NewType("content_sanitization_failed")
	ErrTypeChaptersReorder        = AppErrors.NewType("chapters_reorder")
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

type GetUserBooksQuery struct {
	UserID uuid.UUID
	Limit  int
	Offset int
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
	Books []ManagerAuthorBookDto
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

type UpdateBookChapterCommand struct {
	ID              int64
	Name            string
	Content         string
	IsAdultOverride bool
	Summary         string
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

type BookManagerService interface {
	CreateBook(ctx context.Context, input CreateBookCommand) (int64, error)
	UpdateBook(ctx context.Context, input UpdateBookCommand) error
	UploadBookCover(ctx context.Context, input UploadBookCoverCommand) (UploadBookCoverResult, error)
	UpdateBookChaptersOrder(ctx context.Context, input UpdateBookChaptersOrders) error

	GetBook(ctx context.Context, query ManagerGetBookQuery) (ManagerGetBookResult, error)
	GetUserBooks(ctx context.Context, input GetUserBooksQuery) (GetUserBooksResult, error)
	CreateBookChapter(ctx context.Context, input CreateBookChapterCommand) (CreateBookChapterResult, error)
	UpdateBookChapter(ctx context.Context, input UpdateBookChapterCommand) error
	ReorderChapters(ctx context.Context, input ReorderChaptersCommand) error
	GetBookChapters(ctx context.Context, query ManagerGetBookChaptersQuery) (ManagerGetBookChapterResult, error)
	GetChapter(ctx context.Context, query ManagerGetChapterQuery) (ManagerGetChapterResult, error)
}
