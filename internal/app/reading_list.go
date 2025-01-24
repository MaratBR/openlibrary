package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type ReadingListStatus store.ReadingListStatus

var (
	ReadingListErrors            = AppErrors.NewSubNamespace("reading_list")
	ReadingListInvalidTransition = ReadingListErrors.NewType("invalid_transition")
	ErrReadingListEmptyBook      = ReadingListInvalidTransition.NewSubtype("empty_book").New("cannot update reading list for empty book")
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
	LastUpdatedAt time.Time             `json:"lastUpdatedAt"`
}

type ReadingListService interface {
	MarksAsWantToRead(ctx context.Context, userID uuid.UUID, bookID int64) error
	MarkAsDnf(ctx context.Context, userID uuid.UUID, bookID int64) error
	MarkAsRead(ctx context.Context, userID uuid.UUID, bookID int64) error
	MarkAsPaused(ctx context.Context, userID uuid.UUID, bookID int64) error
	MarkAsReading(ctx context.Context, userID uuid.UUID, bookID int64) error
	MarkAsReadingWithChapterID(ctx context.Context, userID uuid.UUID, bookID int64, chapterID int64) error

	GetStatus(ctx context.Context, userID uuid.UUID, bookID int64) (Nullable[BookReadingListDto], error)
}
