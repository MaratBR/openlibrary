package app

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
)

type bookService struct {
	queries            *store.Queries
	tagsService        TagsService
	uploadService      *UploadService
	readingListService ReadingListService
	reviewService      ReviewsService
	log                *zap.SugaredLogger
}

// GetRandomBookID implements BookService.
func (s *bookService) GetRandomBookID(ctx context.Context) (Nullable[int64], error) {
	ids, err := s.queries.GetRandomPublicBookIDs(ctx, 1)
	if err != nil {
		return Nullable[int64]{}, apperror.WrapUnexpectedDBError(err)
	}
	if len(ids) == 0 {
		return Null[int64](), nil
	}
	return Value(ids[0]), nil
}

func NewBookService(
	db store.DBTX,
	tagsService TagsService,
	uploadService *UploadService,
	readingListService ReadingListService,
	reviewService ReviewsService,
	log *zap.SugaredLogger,
) BookService {
	return &bookService{
		queries:            store.New(db),
		tagsService:        tagsService,
		uploadService:      uploadService,
		readingListService: readingListService,
		reviewService:      reviewService,
		log:                log,
	}
}

func getWordsPerChapter(words, chapters int) int {
	if chapters == 0 {
		return 0
	}

	return words / chapters
}

func (s *bookService) GetBookDetails(ctx context.Context, query GetBookQuery) (result BookDetailsDto, err error) {
	ctx, span := startSpan(ctx, "BookService.GetBookDetails")
	defer func() { endSpan(span, err) }()
	book, err := s.queries.Book_Get(ctx, query.ID)
	if err != nil {
		if err == store.ErrNoRows {
			return BookDetailsDto{}, ErrTypeBookNotFound.New("book with id %d not found", query.ID)
		}

		return BookDetailsDto{}, err
	}

	userPermissionState := getUserBookPermissionsState(getUserBookPermissionsStateRequest{
		IsPubliclyVisible: book.IsPubliclyVisible,
		UserID:            query.ActorUserID,
		BookAuthorID:      uuidDbToDomain(book.AuthorUserID),
	})
	override, err := s.authorizeBookOverride(ctx, query.ActorUserID, query.AdminOverride)
	if err != nil {
		return BookDetailsDto{}, err
	}
	reasons := bookVisibilityReasons(book.IsPubliclyVisible, book.IsBanned, book.IsShadowBanned, book.IsTrashed, book.IsPermRemoved)
	if !override && (!userPermissionState.CanView || book.IsBanned || book.IsShadowBanned || book.IsTrashed || book.IsPermRemoved) {
		return BookDetailsDto{}, ErrTypeBookPrivated.New("book %d cannot be seen", book.ID)
	}

	ageRating := ageRatingFromDbValue(book.AgeRating)
	authorID := uuidDbToDomain(book.AuthorUserID)
	tags, err := s.tagsService.GetTagsByIds(ctx, book.TagIds)
	if err != nil {
		return BookDetailsDto{}, err
	}

	var firstChapterID Nullable[int64]
	firstChapterIDInt, err := s.queries.Book_GetFirstChapterID(ctx, store.Book_GetFirstChapterIDParams{BookID: book.ID, AdminOverride: override})
	if err != nil {
		if err != store.ErrNoRows {
			s.log.Warnw("failed to execute Book_GetFirstChapterID", "err", err)
		}
	} else {
		firstChapterID = Value(firstChapterIDInt)
	}

	bookDto := BookDetailsDto{
		AdminOverride: override,
		IsBanned:      book.IsBanned, IsShadowBanned: book.IsShadowBanned, IsTrashed: book.IsTrashed, IsPermRemoved: book.IsPermRemoved,
		ID:              book.ID,
		Name:            book.Name,
		Slug:            book.Slug,
		AgeRating:       ageRating,
		IsAdult:         ageRating.IsAdult(),
		Tags:            tags,
		Summary:         book.Summary,
		Words:           int(book.Words),
		WordsPerChapter: getWordsPerChapter(int(book.Words), int(book.Chapters)),
		Chapters:        int(book.Chapters),
		CreatedAt:       book.CreatedAt.Time,
		// TODO: Populate ExternalLinks when book source-link persistence is added.
		ExternalLinks: []BookExternalLinkDto{},
		Collections:   []BookCollectionDto{},
		Author: BookDetailsAuthorDto{
			ID:   authorID,
			Name: book.AuthorName,
		},
		Permissions:         BookUserPermissions{CanEdit: query.ActorUserID.Valid && authorID == query.ActorUserID.Value},
		Cover:               getBookCover(s.uploadService, book.Cover, book.ID),
		Rating:              float64ToNullable(book.Rating),
		Reviews:             book.TotalReviews,
		Votes:               book.TotalRatings,
		IsPubliclyAvailable: book.IsPubliclyVisible,
		FirstChapterID:      firstChapterID,
	}

	if override {
		bookDto.VisibilityReasons = reasons
	}
	if userPermissionState.IsOwner && !override {
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

	{
		collections, err := s.queries.GetBookCollectionData(ctx, query.ID)
		if err != nil {
			return BookDetailsDto{}, err
		}
		bookDto.Collections = MapSlice(collections, func(collection store.GetBookCollectionDataRow) BookCollectionDto {
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
func (s *bookService) GetBookChapters(ctx context.Context, query GetBookChaptersQuery) ([]ChapterListDto, error) {
	if _, err := s.GetBookDetails(ctx, GetBookQuery{ID: query.ID, ActorUserID: query.ActorUserID, AdminOverride: query.AdminOverride}); err != nil {
		return nil, err
	}
	if query.AdminOverride {
		chapters, err := s.queries.GetAllBookChapters(ctx, query.ID)
		if err != nil {
			return nil, err
		}
		return MapSlice(chapters, func(chapter store.GetAllBookChaptersRow) ChapterListDto {
			return ChapterListDto{ID: chapter.ID, Order: int(chapter.Order), Name: chapter.Name, Words: int(chapter.Words), CreatedAt: chapter.CreatedAt.Time, Summary: chapter.Summary}
		}), nil
	}
	chapters, err := s.queries.Book_GetPubliclyVisibleChapters(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	chapterDtos := MapSlice(chapters, func(chapter store.Book_GetPubliclyVisibleChaptersRow) ChapterListDto {
		return ChapterListDto{
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

func (s *bookService) GetBookChapter(ctx context.Context, query GetBookChapterQuery) (result GetBookChapterResult, err error) {
	ctx, span := startSpan(ctx, "BookService.GetBookChapter")
	defer func() { endSpan(span, err) }()
	book, err := s.GetBookDetails(ctx, GetBookQuery{ID: query.BookID, ActorUserID: query.ActorUserID, AdminOverride: query.AdminOverride})
	if err != nil {
		return GetBookChapterResult{}, err
	}
	chapter, err := s.queries.GetBookChapterWithDetails(ctx, store.GetBookChapterWithDetailsParams{
		AdminOverride: book.AdminOverride,
		ID:            query.ChapterID,
		BookID:        query.BookID,
	})
	if err != nil {
		return GetBookChapterResult{}, err
	}

	if !chapter.IsPubliclyVisible && !book.AdminOverride {
		return GetBookChapterResult{}, ErrTypeBookPrivated.New("chapter cannot be seen")
	}
	reasons := append([]string{}, book.VisibilityReasons...)
	if !chapter.IsPubliclyVisible {
		reasons = append(reasons, "hiddenChapter")
	}
	var (
		prev Nullable[ChapterNextPrevDto]
		next Nullable[ChapterNextPrevDto]
	)

	if chapter.PrevChapterID != 0 {
		prev = Value(ChapterNextPrevDto{
			ID:    chapter.PrevChapterID,
			Name:  chapter.PrevChapterName,
			Order: int32(chapter.Order - 1),
		})
	}

	if chapter.NextChapterID != 0 {
		next = Value(ChapterNextPrevDto{
			ID:    chapter.NextChapterID,
			Name:  chapter.NextChapterName,
			Order: int32(chapter.Order + 1),
		})
	}

	return GetBookChapterResult{
		VisibilityReasons: reasons,
		Chapter: ChapterDto{
			ID:            chapter.ID,
			Name:          chapter.Name,
			Words:         chapter.Words,
			Content:       chapter.Content,
			CreatedAt:     chapter.CreatedAt.Time,
			Order:         chapter.Order,
			Summary:       chapter.Summary,
			PrevChapter:   prev,
			NextChapter:   next,
			BookID:        chapter.BookID,
			CommentsCount: 42, // TODO
			Fonts:         chapter.Fonts,
		},
	}, nil
}

// GetPinnedBooks implements BookService.
func (s *bookService) GetPinnedBooks(ctx context.Context, input GetPinnedUserBooksQuery) (GetPinnedUserBooksResult, error) {
	rows, err := s.queries.Book_GetByUser(ctx, store.Book_GetByUserParams{
		AuthorUserID: uuidDomainToDb(input.UserID),
		Offset:       int32(input.Offset),
		Limit:        int32(input.Limit + 1),
	})
	if err != nil {
		return GetPinnedUserBooksResult{}, apperror.WrapUnexpectedDBError(err)
	}

	hasMore := len(rows) == input.Limit+1
	books := make([]BookListDto, 0, min(len(rows), input.Limit))
	for i := 0; i < min(len(rows), input.Limit); i++ {
		books = append(books, BookListDto{
			ID:        rows[i].ID,
			Name:      rows[i].Name,
			CreatedAt: rows[i].CreatedAt.Time,
			AgeRating: ageRatingFromDbValue(rows[i].AgeRating),
			Words:     int(rows[i].Words),
			WordsPerChapter: getWordsPerChapter(
				int(rows[i].Words),
				int(rows[i].Chapters)),
			Chapters: int(rows[i].Chapters),
			Cover:    getBookCover(s.uploadService, rows[i].Cover, rows[i].ID),
			IsPinned: rows[i].IsPinned,
		})
	}

	return GetPinnedUserBooksResult{
		Books:   books,
		HasMore: hasMore,
	}, nil
}

func (s *bookService) GetBooksById(ctx context.Context, ids []int64) ([]BookListDto, error) {
	rows, err := s.queries.Book_GetByIds(ctx, ids)
	if err != nil {
		return nil, apperror.WrapUnexpectedDBError(err)
	}

	books := make([]BookListDto, 0, len(rows))
	for i := range rows {
		books = append(books, BookListDto{
			ID:   rows[i].ID,
			Slug: rows[i].Slug,
			Name: rows[i].Name,
			Author: BookDetailsAuthorDto{
				ID:   uuidDbToDomain(rows[i].AuthorUserID),
				Name: rows[i].AuthorName,
			},
			CreatedAt: rows[i].CreatedAt.Time,
			AgeRating: ageRatingFromDbValue(rows[i].AgeRating),
			Words:     int(rows[i].Words),
			WordsPerChapter: getWordsPerChapter(
				int(rows[i].Words),
				int(rows[i].Chapters)),
			Chapters: int(rows[i].Chapters),
			Cover:    getBookCover(s.uploadService, rows[i].Cover, rows[i].ID),
			IsPinned: rows[i].IsPinned,
		})
	}

	return books, nil

}

// authorizeBookOverride checks the current persisted role, independently of HTTP flags.
func (s *bookService) authorizeBookOverride(ctx context.Context, actor Nullable[uuid.UUID], requested bool) (bool, error) {
	if !requested {
		return false, nil
	}
	if !actor.Valid {
		return false, ErrTypeBookPrivated.New("administrator privileges required")
	}
	user, err := s.queries.User_Get(ctx, uuidDomainToDb(actor.Value))
	if err != nil {
		return false, err
	}
	if !UserRole(user.Role).IsAdmin() {
		return false, ErrTypeBookPrivated.New("administrator privileges required")
	}
	return true, nil
}

func bookVisibilityReasons(public, banned, shadowBanned, trashed, removed bool) []string {
	reasons := []string{}
	for _, flag := range []struct {
		value  bool
		reason string
	}{{!public, "private"}, {banned, "banned"}, {shadowBanned, "shadowBanned"}, {trashed, "deleted"}, {removed, "removed"}} {
		if flag.value {
			reasons = append(reasons, flag.reason)
		}
	}
	return reasons
}
