package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type bookService struct {
	queries            *store.Queries
	tagsService        TagsService
	uploadService      *UploadService
	readingListService ReadingListService
	reviewService      ReviewsService
}

// GetUserBooks implements BookService.
func (s *bookService) GetUserBooks(ctx context.Context, input GetUserBooksQuery) (GetUserPinnedBooksResult, error) {
	rows, err := s.queries.GetUserBooks(ctx, store.GetUserBooksParams{
		AuthorUserID: uuidDomainToDb(input.UserID),
		Offset:       int32(input.Offset),
		Limit:        int32(input.Limit + 1),
	})
	if err != nil {
		return GetUserPinnedBooksResult{}, wrapUnexpectedDBError(err)
	}

	hasMore := len(rows) == input.Limit+1
	books := make([]PinnedBookDto, 0, min(len(rows), input.Limit))
	for i := 0; i < len(rows); i++ {
		books = append(books, PinnedBookDto{
			ID:        rows[i].ID,
			Name:      rows[i].Name,
			CreatedAt: rows[i].CreatedAt.Time,
			AgeRating: ageRatingFromDbValue(rows[i].AgeRating),
			Words:     int(rows[i].Words),
			WordsPerChapter: getWordsPerChapter(
				int(rows[i].Words),
				int(rows[i].Chapters)),
			Favorites: rows[i].Favorites,
			Chapters:  int(rows[i].Chapters),
			Cover:     getBookCoverURL(s.uploadService, rows[i].ID, rows[i].HasCover),
			IsPinned:  rows[i].IsPinned,
		})
	}

	return GetUserPinnedBooksResult{
		Books:   books,
		HasMore: hasMore,
	}, nil
}

func NewBookService(
	db store.DBTX,
	tagsService TagsService,
	uploadService *UploadService,
	readingListService ReadingListService,
	reviewService ReviewsService,
) BookService {
	return &bookService{
		queries:            store.New(db),
		tagsService:        tagsService,
		uploadService:      uploadService,
		readingListService: readingListService,
		reviewService:      reviewService,
	}
}

type BookCollectionDto struct {
	ID       int64  `json:"id,string"`
	Name     string `json:"name"`
	Position int    `json:"pos"`
	Size     int    `json:"size"`
}

func getWordsPerChapter(words, chapters int) int {
	if chapters == 0 {
		return 0
	}

	return words / chapters
}

func (s *bookService) GetBook(ctx context.Context, query GetBookQuery) (BookDetailsDto, error) {
	book, err := s.queries.GetBook(ctx, query.ID)
	if err != nil {
		return BookDetailsDto{}, err
	}

	userPermissionState := getUserBookPermissionsState(getUserBookPermissionsStateRequest{
		IsPubliclyVisible: book.IsPubliclyVisible,
		UserID:            query.ActorUserID,
		BookAuthorID:      uuidDbToDomain(book.AuthorUserID),
	})
	if !userPermissionState.CanView {
		return BookDetailsDto{}, ErrGenericForbidden
	}

	ageRating := ageRatingFromDbValue(book.AgeRating)
	authorID := uuidDbToDomain(book.AuthorUserID)
	tags, err := s.tagsService.GetTagsByIds(ctx, book.TagIds)
	if err != nil {
		return BookDetailsDto{}, err
	}

	bookDto := BookDetailsDto{
		ID:              book.ID,
		Name:            book.Name,
		AgeRating:       ageRating,
		IsAdult:         ageRating.IsAdult(),
		Tags:            tags,
		Summary:         book.Summary,
		Words:           int(book.Words),
		WordsPerChapter: getWordsPerChapter(int(book.Words), int(book.Chapters)),
		CreatedAt:       book.CreatedAt.Time,
		Collections:     []BookCollectionDto{},
		Favorites:       book.Favorites,
		Author: BookDetailsAuthorDto{
			ID:   authorID,
			Name: book.AuthorName,
		},
		Permissions: BookUserPermissions{CanEdit: query.ActorUserID.Valid && authorID == query.ActorUserID.UUID},
		Cover:       getBookCoverURL(s.uploadService, book.ID, book.HasCover),
		Rating:      float64ToNullable(book.Rating),
		Reviews:     book.TotalReviews,
		Votes:       book.TotalRatings,
	}

	if userPermissionState.IsOwner {
		if !book.IsPubliclyVisible {
			bookDto.Notifications = append(bookDto.Notifications, GenericNotification{
				ID:   "book:owner:not_publicly_visible",
				Text: fmt.Sprintf("This book is not publicly visible, you can change that in the settings if you want. [Click here](/manager/book/%d?tab=info&from=book-page-notification) to edit book settings.<br />Only you can see this message.", book.ID),
			})
		}

		if book.IsBanned {
			bookDto.Notifications = append(bookDto.Notifications, GenericNotification{
				ID:   "book:owner:banned",
				Text: fmt.Sprintf("This book has been banned by our moderation team, please [click here](/manager/book/banned/%d?from=book-page-notification) to find out more about this", book.ID),
			})
		}
	}

	if query.ActorUserID.Valid {
		isFavorite, err := s.queries.IsFavoritedBy(ctx, store.IsFavoritedByParams{
			BookID: query.ID,
			UserID: uuidDomainToDb(query.ActorUserID.UUID),
		})
		if err != nil {
			if err == store.ErrNoRows {
				bookDto.IsFavorite = false
			} else {
				slog.Error("failed to get isFavorite value for a book and a user", "bookID", book.ID, "userID", query.ActorUserID.UUID, "error", err)
			}
		} else {
			bookDto.IsFavorite = isFavorite
		}
	}

	{
		collections, err := s.queries.GetBookCollections(ctx, query.ID)
		if err != nil {
			return BookDetailsDto{}, err
		}
		bookDto.Collections = mapSlice(collections, func(collection store.GetBookCollectionsRow) BookCollectionDto {
			return BookCollectionDto{
				ID:       collection.ID,
				Name:     collection.Name,
				Position: int(collection.Position),
				Size:     int(collection.Size),
			}
		})
	}

	return bookDto, nil
}

// GetBookChapters implements BookService.
func (s *bookService) GetBookChapters(ctx context.Context, query GetBookChaptersQuery) ([]BookChapterDto, error) {
	chapters, err := s.queries.GetBookChapters(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	chapterDtos := mapSlice(chapters, func(chapter store.GetBookChaptersRow) BookChapterDto {
		return BookChapterDto{
			ID:        chapter.ID,
			Order:     int(chapter.Order),
			Name:      chapter.Name,
			Words:     int(chapter.Words),
			CreatedAt: chapter.CreatedAt.Time,
			Summary:   chapter.Summary,
		}
	})

	return chapterDtos, nil
}

type ChapterNextPrevDto struct {
	ID    int64  `json:"id,string"`
	Name  string `json:"name"`
	Order int32  `json:"order"`
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
}

type ChapterWithDetails struct {
	Chapter ChapterDto
}

type GetBookChapterQuery struct {
	BookID      int64
	ChapterID   int64
	ActorUserID uuid.NullUUID
}

type GetBookChapterResult struct {
	Chapter ChapterWithDetails
}

func (s *bookService) GetBookChapter(ctx context.Context, query GetBookChapterQuery) (GetBookChapterResult, error) {
	chapter, err := s.queries.GetBookChapterWithDetails(ctx, store.GetBookChapterWithDetailsParams{
		ID:     query.ChapterID,
		BookID: query.BookID,
	})
	if err != nil {
		return GetBookChapterResult{}, err
	}

	var (
		prev Nullable[ChapterNextPrevDto]
		next Nullable[ChapterNextPrevDto]
	)

	if chapter.PrevChapterID.Valid {
		prev = Value(ChapterNextPrevDto{
			ID:    chapter.PrevChapterID.Int64,
			Name:  chapter.PrevChapterName.String,
			Order: int32(chapter.Order - 1),
		})
	}

	if chapter.NextChapterID.Valid {
		next = Value(ChapterNextPrevDto{
			ID:    chapter.NextChapterID.Int64,
			Name:  chapter.NextChapterName.String,
			Order: int32(chapter.Order + 1),
		})
	}

	return GetBookChapterResult{
		Chapter: ChapterWithDetails{
			Chapter: ChapterDto{
				ID:              chapter.ID,
				Name:            chapter.Name,
				Words:           chapter.Words,
				Content:         chapter.Content,
				IsAdultOverride: chapter.IsAdultOverride,
				CreatedAt:       chapter.CreatedAt.Time,
				Order:           chapter.Order,
				Summary:         chapter.Summary,
				PrevChapter:     prev,
				NextChapter:     next,
			},
		},
	}, nil
}
