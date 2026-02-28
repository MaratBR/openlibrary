package app

import (
	"context"
	"io"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/analytics"
	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/gofrs/uuid"
	"github.com/joomcode/errorx"
)

var (
	ErrTypeBookSanitizationFailed = apperror.AppErrors.NewType("content_sanitization_failed")
	ErrTypeChaptersReorder        = apperror.AppErrors.NewType("chapters_reorder")
	ErrDraftNotFound              = apperror.AppErrors.NewType("draft_not_found", errorx.NotFound(), apperror.ErrTraitEntityNotFound).New("draft not found")
	ErrTypeChapterDoesNotExist    = apperror.AppErrors.NewType("chapter_not_found", apperror.ErrTraitEntityNotFound)
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
	ID                int64                   `json:"id,string"`
	Name              string                  `json:"name"`
	AgeRating         AgeRating               `json:"ageRating"`
	IsAdult           bool                    `json:"adult"`
	Tags              []DefinedTagDto         `json:"tags"`
	Words             int                     `json:"words"`
	WordsPerChapter   int                     `json:"wordsPerChapter"`
	CreatedAt         time.Time               `json:"createdAt"`
	Collections       []BookCollectionDto     `json:"collections"`
	Chapters          []ManagerBookChapterDto `json:"chapters"`
	Author            BookDetailsAuthorDto    `json:"author"`
	Summary           string                  `json:"summary"`
	IsPubliclyVisible bool                    `json:"isPubliclyVisible"`
	IsBanned          bool                    `json:"isBanned"`
	Cover             BookCover               `json:"cover"`
}

type ManagerGetBookQuery_Result struct {
	Book ManagerBookDetailsDto
}

type ManagerBookDto_Stats struct {
	Views   analytics.Views `json:"views"`
	Reviews int32           `json:"reviews"`
	Ratings int32           `json:"ratings"`
}

type ManagerBookDto struct {
	ID                int64                `json:"id,string"`
	Slug              string               `json:"slug"`
	Name              string               `json:"name"`
	CreatedAt         time.Time            `json:"createdAt"`
	AgeRating         AgeRating            `json:"ageRating"`
	Tags              []DefinedTagDto      `json:"tags"`
	Words             int                  `json:"words"`
	WordsPerChapter   int                  `json:"wordsPerChapter"`
	Chapters          int                  `json:"chapters"`
	Collections       []BookCollectionDto  `json:"collections"`
	IsPubliclyVisible bool                 `json:"isPubliclyVisible"`
	IsBanned          bool                 `json:"isBanned"`
	IsTrashed         bool                 `json:"isTrashed"`
	Summary           string               `json:"summary"`
	Cover             BookCover            `json:"cover"`
	Stats             ManagerBookDto_Stats `json:"stats"`
}

type ManagerGetUserBooksQuery_Result struct {
	Books      []ManagerBookDto
	TotalPages uint32
	PageSize   uint32
	Page       uint32
}

type CreateBookChapterCommand struct {
	BookID            int64
	Name              string
	Content           string
	IsAdultOverride   bool
	Summary           string
	UserID            uuid.UUID
	IsPubliclyVisible bool
}

type CreateBookChapterCommand_Result struct {
	ID int64
}

type ReorderChaptersCommand struct {
	UserID     uuid.UUID
	BookID     int64
	ChapterIDs []int64
}

type ManagerBookChapterDto struct {
	ID                int64                 `json:"id,string"`
	Name              string                `json:"name"`
	CreatedAt         time.Time             `json:"createdAt"`
	Words             int                   `json:"words"`
	Summary           string                `json:"summary"`
	Order             int32                 `json:"order"`
	IsAdultOverride   bool                  `json:"isAdultOverride"`
	IsPubliclyVisible bool                  `json:"isPubliclyVisible"`
	DraftID           Nullable[Int64String] `json:"draftId"`
}

type ManagerGetBookChaptersQuery struct {
	BookID int64
	UserID uuid.UUID
}

type ManagerGetBookChaptersQuery_Result struct {
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

type ManagerGetChapterQuery_Result struct {
	Chapter ManagerBookChapterDetailsDto
}

type UploadBookCoverCommand struct {
	UserID uuid.UUID
	BookID int64
	File   io.Reader
}

type UploadBookCoverCommand_Result struct {
	URL BookCover
}

type UpdateBookChapterOrdersCommand struct {
	BookID        int64
	Modifications []ChapterOrderModification
}

type ChapterOrderModification struct {
	ChapterID        int64
	NewPositionIndex int
}

type UpdateBookChapterOrdersCommand_Result struct {
	ModifiedPositions map[int64]int
}

type GetDraftQuery struct {
	UserID    uuid.UUID
	DraftID   int64
	ChapterID int64
	BookID    int64
}

type DraftDto struct {
	ID          int64               `json:"id,string"`
	ChapterName string              `json:"chapterName"`
	Content     string              `json:"content"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   Nullable[time.Time] `json:"updatedAt"`
	CreatedBy   struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	} `json:"createdBy"`
	Book struct {
		ID   int64  `json:"id,string"`
		Name string `json:"name"`
	} `json:"book"`
	Chapter struct {
		ID               int64     `json:"id,string"`
		ContentUpdatedAt time.Time `json:"contentUpdatedAt"`
	} `json:"chapter"`
	IsChapterPubliclyAvailable bool `json:"isChapterPubliclyAvailable"`
}

type UpdateDraftChapterNameCommand struct {
	ChapterName string
	ChapterID   int64
	BookID      int64
	DraftID     int64
	UserID      uuid.UUID
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
	DraftID    int64
	UserID     uuid.UUID
	MakePublic bool
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

type TrashBookCommand struct {
	BookID      int64
	ActorUserID uuid.UUID
}

type UntrashBookCommand struct {
	BookID      int64
	ActorUserID uuid.UUID
}

type BookManagerService interface {
	GetUserBooks(ctx context.Context, input ManagerGetUserBooksQuery) (ManagerGetUserBooksQuery_Result, error)
	GetBook(ctx context.Context, query ManagerGetBookQuery) (ManagerGetBookQuery_Result, error)
	CreateBook(ctx context.Context, input CreateBookCommand) (int64, error)
	UpdateBook(ctx context.Context, input UpdateBookCommand) error
	UploadBookCover(ctx context.Context, input UploadBookCoverCommand) (UploadBookCoverCommand_Result, error)
	TrashBook(ctx context.Context, input TrashBookCommand) error
	UntrashBook(ctx context.Context, input UntrashBookCommand) error

	UpdateBookChaptersOrder(ctx context.Context, input UpdateBookChapterOrdersCommand) (UpdateBookChapterOrdersCommand_Result, error)
	CreateBookChapter(ctx context.Context, input CreateBookChapterCommand) (CreateBookChapterCommand_Result, error)
	ReorderChapters(ctx context.Context, input ReorderChaptersCommand) error
	GetBookChapters(ctx context.Context, query ManagerGetBookChaptersQuery) (ManagerGetBookChaptersQuery_Result, error)
	GetChapter(ctx context.Context, query ManagerGetChapterQuery) (ManagerGetChapterQuery_Result, error)

	GetDraft(ctx context.Context, query GetDraftQuery) (DraftDto, error)
	UpdateDraft(ctx context.Context, cmd UpdateDraftCommand) error
	UpdateDraftChapterName(ctx context.Context, cmd UpdateDraftChapterNameCommand) error
	UpdateDraftContent(ctx context.Context, cmd UpdateDraftContentCommand) error
	DeleteDraft(ctx context.Context, cmd DeleteDraftCommand) error
	PublishDraft(ctx context.Context, cmd PublishDraftCommand) error
	CreateDraft(ctx context.Context, cmd CreateDraftCommand) (int64, error)
	GetLatestDraft(ctx context.Context, cmd GetLatestDraftQuery) (Nullable[int64], error)
}
