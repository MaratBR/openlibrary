package app

import (
	"context"
	"fmt"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type readingListService struct {
	db            store.DBTX
	uploadService *UploadService
}

// GetReadingListBooks implements ReadingListService.
func (s *readingListService) GetReadingListBooks(ctx context.Context, query GetReadingListItemsQuery) ([]BookLibraryDto, error) {
	queries := store.New(s.db)
	rows, err := queries.GetUserLibrary(ctx, store.GetUserLibraryParams{
		UserID: uuidDomainToDb(query.UserID),
		Status: store.ReadingListStatus(query.Status),
		Limit:  int32(query.Limit),
	})
	if err != nil {
		return nil, wrapUnexpectedDBError(err)
	}
	return mapSlice(rows, func(r store.GetUserLibraryRow) BookLibraryDto {
		var lastChapter Nullable[BookReadingListItemLastChapterDto]
		if r.ChapterID.Valid {
			lastChapter = Value(BookReadingListItemLastChapterDto{
				ID:    r.ChapterID.Int64,
				Name:  r.ChapterName.String,
				Order: r.ChapterOrder.Int32,
			})
		}

		return BookLibraryDto{
			ID:          r.ID,
			Name:        r.Name,
			Cover:       getBookCoverURL(s.uploadService, r.ID, r.HasCover),
			AgeRating:   ageRatingFromDbValue(r.AgeRating),
			LastChapter: lastChapter,
		}
	}), nil
}

// GetStatus implements ReadingListService.
func (r *readingListService) GetStatus(ctx context.Context, userID uuid.UUID, bookID int64) (Nullable[BookReadingListDto], error) {
	queries := store.New(r.db)

	status, err := queries.GetBookReadingListState(ctx, store.GetBookReadingListStateParams{
		UserID: uuidDomainToDb(userID),
		BookID: bookID,
	})
	if err != nil {
		if err == store.ErrNoRows {
			return Null[BookReadingListDto](), nil
		}
		return Null[BookReadingListDto](), wrapUnexpectedDBError(err)
	}

	return Value(BookReadingListDto{
		Status:        ReadingListStatus(status.Status),
		ChapterID:     int64ToNullable(status.LastAccessedChapterID),
		ChapterOrder:  status.ChapterOrder.Int32,
		LastUpdatedAt: timeDbToDomain(status.LastUpdatedAt),
	}), nil
}

// MarksAsWantToRead implements ReadingListService.
func (r *readingListService) MarksAsWantToRead(ctx context.Context, userID uuid.UUID, bookID int64) error {
	return r.setStatus(ctx, userID, bookID, store.ReadingListStatusWantToRead)
}

// MarkAsDnf implements ReadingListService.
func (r *readingListService) MarkAsDnf(ctx context.Context, userID uuid.UUID, bookID int64) error {
	return r.setStatus(ctx, userID, bookID, store.ReadingListStatusDnf)
}

// MarkAsPaused implements ReadingListService.
func (r *readingListService) MarkAsPaused(ctx context.Context, userID uuid.UUID, bookID int64) error {
	return r.setStatus(ctx, userID, bookID, store.ReadingListStatusPaused)
}

// MarkAsRead implements ReadingListService.
func (r *readingListService) MarkAsRead(ctx context.Context, userID uuid.UUID, bookID int64) error {
	return r.setStatus(ctx, userID, bookID, store.ReadingListStatusRead)
}

// MarkAsReadingWithChapterID implements ReadingListService.
func (r *readingListService) MarkAsReadingWithChapterID(ctx context.Context, userID uuid.UUID, bookID int64, chapterID int64) error {
	return r.setStatusAndChapter(ctx, userID, bookID, store.ReadingListStatusReading, chapterID)
}

// MarkAsReading implements ReadingListService.
func (r *readingListService) MarkAsReading(ctx context.Context, userID uuid.UUID, bookID int64) error {
	return r.setStatus(ctx, userID, bookID, store.ReadingListStatusReading)
}

func (r *readingListService) setStatus(
	ctx context.Context,
	userID uuid.UUID,
	bookID int64,
	status store.ReadingListStatus) error {
	queries := store.New(r.db)
	_, err := queries.SetBookReadingListStatus(ctx, store.SetBookReadingListStatusParams{
		UserID: uuidDomainToDb(userID),
		BookID: bookID,
		Status: status,
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil
}

func (r *readingListService) setStatusAndChapter(
	ctx context.Context,
	userID uuid.UUID,
	bookID int64,
	status store.ReadingListStatus,
	chapterID int64) error {
	queries := store.New(r.db)
	_, err := queries.SetBookReadingListStatusAndChapter(ctx, store.SetBookReadingListStatusAndChapterParams{
		UserID:                uuidDomainToDb(userID),
		BookID:                bookID,
		Status:                status,
		LastAccessedChapterID: pgtype.Int8{Valid: true, Int64: chapterID},
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil
}

func (r *readingListService) MarkChapterRead(ctx context.Context, command MarkChapterCommand) error {
	queries := store.New(r.db)
	// TODO: account for public availability
	bookID, err := queries.GetChapterBookID(ctx, command.ChapterID)
	if err != nil {
		if err == store.ErrNoRows {
			return ReadingListChapterNotFound.New(fmt.Sprintf("chapter %d not found", command.ChapterID), command.ChapterID)
		} else {
			return wrapUnexpectedDBError(err)
		}
	}
	err = queries.SetBookReadingListChapter(
		ctx,
		store.SetBookReadingListChapterParams{
			BookID:                bookID,
			LastAccessedChapterID: pgtype.Int8{Valid: true, Int64: command.ChapterID},
			UserID:                uuidDomainToDb(command.UserID),
		},
	)

	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil
}

func NewReadingListService(db store.DBTX, uploadService *UploadService) ReadingListService {
	return &readingListService{db: db, uploadService: uploadService}
}
