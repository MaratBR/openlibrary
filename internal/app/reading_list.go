package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type ReadingListStatus store.ReadingListStatus

var (
	ReadingListErrors            = apperror.AppErrors.NewSubNamespace("reading_list")
	ReadingListInvalidTransition = ReadingListErrors.NewType("invalid_transition")
	ErrReadingListEmptyBook      = ReadingListInvalidTransition.NewSubtype("empty_book").New("cannot update reading list for empty book")
	ReadingListChapterNotFound   = ReadingListErrors.NewType("chapter404")
)

const (
	ReadingListStatusDnf        ReadingListStatus = ReadingListStatus(store.ReadingListStatusDnf)
	ReadingListStatusPaused     ReadingListStatus = ReadingListStatus(store.ReadingListStatusPaused)
	ReadingListStatusRead       ReadingListStatus = ReadingListStatus(store.ReadingListStatusRead)
	ReadingListStatusReading    ReadingListStatus = ReadingListStatus(store.ReadingListStatusReading)
	ReadingListStatusWantToRead ReadingListStatus = ReadingListStatus(store.ReadingListStatusWantToRead)
)

type BookReadingListDto struct {
	Status        ReadingListStatus     `json:"status"`
	ChapterID     Nullable[Int64String] `json:"chapterId"`
	ChapterName   string                `json:"chapterName"`
	ChapterOrder  int32                 `json:"chapterOrder"`
	LastUpdatedAt time.Time             `json:"lastUpdatedAt"`
}

type BookLibraryDto struct {
	ID          int64                                       `json:"id,string"`
	Name        string                                      `json:"name"`
	Cover       string                                      `json:"cover"`
	AgeRating   AgeRating                                   `json:"ageRating"`
	LastChapter Nullable[BookReadingListItemLastChapterDto] `json:"lastChapter"`
}

type BookReadingListItemLastChapterDto struct {
	ID    int64  `json:"id,string"`
	Name  string `json:"name"`
	Order int32  `json:"order"`
}

type GetReadingListItemsQuery struct {
	UserID uuid.UUID
	Limit  uint32
	Status ReadingListStatus
}

type MarkChapterCommand struct {
	ChapterID int64
	UserID    uuid.UUID
}

type ReadingListService interface {
	MarksAsWantToRead(ctx context.Context, userID uuid.UUID, bookID int64) error
	MarkAsDnf(ctx context.Context, userID uuid.UUID, bookID int64) error
	MarkAsRead(ctx context.Context, userID uuid.UUID, bookID int64) error
	MarkAsPaused(ctx context.Context, userID uuid.UUID, bookID int64) error
	MarkAsReading(ctx context.Context, userID uuid.UUID, bookID int64) error
	MarkAsReadingWithChapterID(ctx context.Context, userID uuid.UUID, bookID int64, chapterID int64) error
	GetStatus(ctx context.Context, userID uuid.UUID, bookID int64) (Nullable[BookReadingListDto], error)
	GetReadingListBooks(ctx context.Context, query GetReadingListItemsQuery) ([]BookLibraryDto, error)
	MarkChapterRead(ctx context.Context, command MarkChapterCommand) error
}
