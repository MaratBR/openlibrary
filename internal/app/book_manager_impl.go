package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"slices"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/MaratBR/openlibrary/internal/app/analytics"
	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/app/bookfont"
	"github.com/MaratBR/openlibrary/internal/app/content"
	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/app/imgconvert"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/minio/minio-go/v7"
)

type BookManagerServiceDeps struct {
	fx.In

	DB                 DB
	TagsService        TagsService
	UsersService       UserService
	UploadService      *UploadService
	BookReindexService BookReindexService
	MetricService      analytics.MetricService
	Log                *zap.SugaredLogger
	Markup             *content.MarkupEngine
	FontPolicy         bookfont.Policy
}

type bookManagerService struct {
	queries *store.Queries
	deps    BookManagerServiceDeps
}

const (
	BOOK_COVER_DIRECTORY = "book-covers"
)

func (s *bookManagerService) GetUserBooks(ctx context.Context, input ManagerGetUserBooksQuery) (ManagerGetUserBooksQuery_Result, error) {
	var (
		limit  int32
		offset int32
	)

	limit = int32(input.PageSize)
	offset = int32((input.Page - 1) * input.PageSize)
	if offset < 0 {
		offset = 0
	}

	books, err := s.queries.Book_ManagerGetUserBooks(ctx, store.Book_ManagerGetUserBooksParams{
		AuthorUserID: uuidDomainToDb(input.UserID),
		Limit:        limit,
		Offset:       offset,
		Search:       input.SearchQuery,
	})
	if err != nil {
		return ManagerGetUserBooksQuery_Result{}, apperror.WrapUnexpectedDBError(err)
	}

	userBooks, err := s.aggregateUserBooks(ctx, books)
	if err != nil {
		return ManagerGetUserBooksQuery_Result{}, err
	}

	count, err := s.queries.Book_Book_ManagerGetUserBooksCount(ctx, store.Book_Book_ManagerGetUserBooksCountParams{
		AuthorUserID: uuidDomainToDb(input.UserID),
		Search:       input.SearchQuery,
	})

	totalPages := uint32(math.Ceil(float64(count) / float64(input.PageSize)))

	if err != nil {
		return ManagerGetUserBooksQuery_Result{}, apperror.WrapUnexpectedDBError(err)
	}

	// load views
	bookIDs := MapSlice(userBooks, func(b ManagerBookDto) int64 { return b.ID })
	viewMappings, err := s.deps.MetricService.Get(ctx, analytics.MetricViews, bookIDs)
	if err != nil {
		return ManagerGetUserBooksQuery_Result{}, err
	}

	for i := range userBooks {
		dto := userBooks[i]
		views, ok := viewMappings[dto.ID]
		if ok {
			dto.Stats.Views = views
			userBooks[i] = dto
		}
	}

	return ManagerGetUserBooksQuery_Result{Books: userBooks, PageSize: input.PageSize, TotalPages: totalPages, Page: input.Page}, nil
}

func (s *bookManagerService) GetBook(ctx context.Context, query ManagerGetBookQuery) (ManagerGetBookQuery_Result, error) {
	book, err := s.queries.Book_Get(ctx, query.BookID)
	if err != nil {
		return ManagerGetBookQuery_Result{}, err
	}

	tags, err := s.deps.TagsService.GetTagsByIds(ctx, book.TagIds)
	if err != nil {
		return ManagerGetBookQuery_Result{}, err
	}

	ageRating := ageRatingFromDbValue(book.AgeRating)
	authorID := uuidDbToDomain(book.AuthorUserID)

	bookDto := ManagerBookDetailsDto{
		ID:              book.ID,
		Name:            book.Name,
		AgeRating:       ageRating,
		IsAdult:         ageRating.IsAdult(),
		Tags:            tags,
		Words:           int(book.Words),
		WordsPerChapter: getWordsPerChapter(int(book.Words), int(book.Chapters)),
		CreatedAt:       book.CreatedAt.Time,
		Collections:     []BookCollectionDto{},
		Chapters:        []ManagerBookChapterDto{},
		Author: BookDetailsAuthorDto{
			ID:   authorID,
			Name: book.AuthorName,
		},
		Summary:           book.Summary,
		IsPubliclyVisible: book.IsPubliclyVisible,
		IsBanned:          book.IsBanned,
		Cover:             getBookCover(s.deps.UploadService, book.Cover, book.ID),
	}

	{
		chapters, err := s.queries.GetAllBookChapters(ctx, query.BookID)
		if err != nil {
			return ManagerGetBookQuery_Result{}, err
		}
		bookDto.Chapters = MapSlice(chapters, func(chapter store.GetAllBookChaptersRow) ManagerBookChapterDto {
			var draftID Nullable[Int64String]

			if chapter.LatestDraftID != 0 {
				draftID = Value(Int64String(chapter.LatestDraftID))
			}

			return ManagerBookChapterDto{
				ID:                chapter.ID,
				Order:             chapter.Order,
				Name:              chapter.Name,
				Words:             int(chapter.Words),
				CreatedAt:         chapter.CreatedAt.Time,
				Summary:           chapter.Summary,
				IsPubliclyVisible: chapter.IsPubliclyVisible,
				DraftID:           draftID,
				ScheduledAt:       timeNullableDbToDomain(chapter.ScheduledAt),
			}
		})
	}

	{
		collections, err := s.queries.GetBookCollectionData(ctx, query.BookID)
		if err != nil {
			return ManagerGetBookQuery_Result{}, err
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

	return ManagerGetBookQuery_Result{
		Book: bookDto,
	}, nil
}

func (s *bookManagerService) CreateBook(ctx context.Context, input CreateBookCommand) (int64, error) {
	err := validateBookName(input.Name)
	if err != nil {
		return 0, err
	}

	err = validateBookSummary(input.Summary)
	if err != nil {
		return 0, err
	}

	tags, err := s.deps.TagsService.FindParentTagIds(ctx, input.Tags)
	if err != nil {
		return 0, err
	}

	id := GenID()
	err = s.queries.Book_Insert(ctx, store.Book_InsertParams{
		ID:                 id,
		Name:               input.Name,
		Slug:               makeSlug(input.Name),
		AuthorUserID:       uuidDomainToDb(input.UserID),
		CreatedAt:          timeToTimestamptz(time.Now()),
		TagIds:             tags.TagIds,
		CachedParentTagIds: tags.ParentTagIds,
		AgeRating:          ageRatingDbValue(input.AgeRating),
		Summary:            input.Summary,
		IsPubliclyVisible:  input.IsPubliclyVisible,
	})
	if err != nil {
		return 0, apperror.WrapUnexpectedDBError(err)
	}

	s.deps.BookReindexService.ScheduleReindex(ctx, id)

	return id, err
}

func (s *bookManagerService) UpdateBook(ctx context.Context, input UpdateBookCommand) error {
	err := validateBookName(input.Name)
	if err != nil {
		return err
	}

	summaryData, err := s.deps.Markup.Clean(input.Summary)
	if err != nil {
		return err
	}

	err = validateBookSummary(summaryData.Sanitized)
	if err != nil {
		return err
	}

	tags, err := s.deps.TagsService.FindParentTagIds(ctx, input.Tags)
	if err != nil {
		return err
	}

	err = s.queries.Book_Update(ctx, store.Book_UpdateParams{
		ID:                 input.BookID,
		Name:               input.Name,
		TagIds:             tags.TagIds,
		CachedParentTagIds: tags.ParentTagIds,
		AgeRating:          ageRatingDbValue(input.AgeRating),
		Summary:            summaryData.Sanitized,
		IsPubliclyVisible:  input.IsPubliclyVisible,
	})
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	s.deps.BookReindexService.ScheduleReindex(ctx, input.BookID)

	return nil
}

// UploadBookCover implements BookManagerService.
func (s *bookManagerService) UploadBookCover(ctx context.Context, input UploadBookCoverCommand) (result UploadBookCoverCommand_Result, err error) {
	ctx, span := startSpan(ctx, "BookManagerService.UploadBookCover")
	defer func() { endSpan(span, err) }()
	file, err := io.ReadAll(input.File)
	if err != nil {
		return
	}

	imgBytes, err := imgconvert.ConvertToJPEG(file)
	if err != nil {
		return
	}

	imgBytes, err = imgconvert.Resize(imgBytes, 300, 300)
	if err != nil {
		return
	}

	cover := fmt.Sprintf("%d_%d", input.BookID, GenID())
	path := fmt.Sprintf("%s/%s.jpeg", BOOK_COVER_DIRECTORY, cover)
	_, err = s.deps.UploadService.Client.PutObject(
		ctx,
		s.deps.UploadService.PublicBucket,
		path,
		bytes.NewReader(imgBytes),
		int64(len(imgBytes)),
		minio.PutObjectOptions{ContentType: "image/jpeg"},
	)
	if err != nil {
		return
	}

	err = s.queries.BookSetCover(ctx, store.BookSetCoverParams{
		ID:    input.BookID,
		Cover: cover,
	})
	if err != nil {
		return
	}

	result.URL = getBookCover(s.deps.UploadService, cover, input.BookID)

	return
}

func (s *bookManagerService) TrashBook(ctx context.Context, input TrashBookCommand) error {
	// TODO auth
	err := s.queries.Book_Trash(ctx, input.BookID)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil
}

func (s *bookManagerService) UntrashBook(ctx context.Context, input UntrashBookCommand) error {
	// TODO auth
	err := s.queries.Book_UnTrash(ctx, store.Book_UnTrashParams{
		ID:                input.BookID,
		IsPubliclyVisible: false,
	})
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil
}

// UpdateBookChaptersOrder updates the order of chapters in a book.
func (s *bookManagerService) UpdateBookChaptersOrder(ctx context.Context, input UpdateBookChapterOrdersCommand) (UpdateBookChapterOrdersCommand_Result, error) {
	// no modification is considered correct - just do nothing
	if len(input.Modifications) == 0 {
		return UpdateBookChapterOrdersCommand_Result{}, nil
	}

	tx, err := s.deps.DB.Begin(ctx)
	if err != nil {
		return UpdateBookChapterOrdersCommand_Result{}, err
	}

	queries := s.queries.WithTx(tx)

	oldOrder, err := queries.GetChaptersOrder(ctx, input.BookID)
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return UpdateBookChapterOrdersCommand_Result{}, err
	}
	oldIndices := make(map[int64]int)
	newOrder := make([]int64, len(oldOrder))
	for i := 0; i < len(oldOrder); i++ {
		oldIndices[oldOrder[i]] = i
		newOrder[i] = oldOrder[i]
	}

	for _, modification := range input.Modifications {
		_, ok := MoveItem(newOrder, modification.ChapterID, modification.NewPositionIndex)
		if !ok {
			s.deps.Log.Errorw("failed to apply chapter modification", "chapterID", modification.ChapterID, "newPositionIndex", modification.NewPositionIndex)
			continue
		}
	}

	modifiedPositions := make(map[int64]int)

	for i, chapterID := range newOrder {
		oldPosition, ok := oldIndices[chapterID]
		if ok && oldPosition == i {
			// not changed, ignore
			continue
		}

		modifiedPositions[chapterID] = i + 1

		s.deps.Log.Debugw("updating position of the chapter", "chapterID", chapterID, "index", i)
		err = queries.Book_SetChapterOrder(ctx, store.Book_SetChapterOrderParams{
			ID:    chapterID,
			Order: int32(i + 1),
		})
		if err != nil {
			dal.RollbackTx(ctx, tx)
			return UpdateBookChapterOrdersCommand_Result{}, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return UpdateBookChapterOrdersCommand_Result{}, err
	}

	return UpdateBookChapterOrdersCommand_Result{
		ModifiedPositions: modifiedPositions,
	}, nil
}

// MoveItem moves an item in the list, returns new position of item and a boolean representing
// the successfullness of the operation (true if successful), if operation was unsuccessful then
// newPosition returned SHOULD be -1
func MoveItem[T comparable](arr []T, item T, newPositionIndex int) (int, bool) {
	idx := slices.Index(arr, item)
	if idx == -1 {
		return -1, false
	}

	if newPositionIndex < 0 {
		return -1, false
	}

	if newPositionIndex >= len(arr) {
		newPositionIndex = len(arr) - 1
	}

	moveArrEl(arr, idx, newPositionIndex)
	return newPositionIndex, true
}

func (s *bookManagerService) aggregateUserBooks(ctx context.Context, rows []store.Book_ManagerGetUserBooksRow) ([]ManagerBookDto, error) {
	var (
		books   []ManagerBookDto = []ManagerBookDto{}
		book    ManagerBookDto
		tagsAgg = newTagsAggregator(s.deps.TagsService)
	)

	for _, row := range rows {
		if row.ID != book.ID {
			if book.ID != 0 {
				books = append(books, book)
			}

			tagsAgg.Add(book.ID, row.TagIds)

			book = ManagerBookDto{
				ID:                row.ID,
				Slug:              row.Slug,
				Name:              row.Name,
				CreatedAt:         row.CreatedAt.Time,
				AgeRating:         ageRatingFromDbValue(row.AgeRating),
				Tags:              nil, // will be set later
				Words:             int(row.Words),
				Chapters:          int(row.Chapters),
				WordsPerChapter:   getWordsPerChapter(int(row.Words), int(row.Chapters)),
				Collections:       []BookCollectionDto{},
				Summary:           row.Summary,
				IsPubliclyVisible: row.IsPubliclyVisible,
				IsBanned:          row.IsBanned,
				IsTrashed:         row.IsTrashed,
				Cover:             getBookCover(s.deps.UploadService, row.Cover, row.ID),
				Stats: ManagerBookDto_Stats{
					Ratings: row.TotalReviews,
					Reviews: row.TotalReviews,
				},
			}
		}

		if row.CollectionID.Valid {
			collection := BookCollectionDto{
				ID:       row.CollectionID.Int64,
				Name:     row.CollectionName.String,
				Position: int(row.CollectionPosition.Int32),
				Size:     int(row.CollectionSize.Int32),
			}
			book.Collections = append(book.Collections, collection)
		}
	}

	if book.ID != 0 {
		books = append(books, book)
	}

	tags, err := tagsAgg.Fetch(ctx)
	if err != nil {
		return []ManagerBookDto{}, err
	}

	for i := 0; i < len(books); i++ {
		bookTagIDs := tagsAgg.BookTags(books[i].ID)
		if bookTagIDs != nil {
			books[i].Tags = MapSlice(bookTagIDs, func(tagID int64) DefinedTagDto {
				return tags[tagID]
			})
		} else {
			books[i].Tags = []DefinedTagDto{}
		}
	}

	return books, nil
}

func (s *bookManagerService) CreateBookChapter(ctx context.Context, input CreateBookChapterCommand) (CreateBookChapterCommand_Result, error) {
	if err := validateChapterName(input.Name); err != nil {
		return CreateBookChapterCommand_Result{}, err
	}

	content, err := s.deps.Markup.Clean(input.Content)
	if err != nil {
		return CreateBookChapterCommand_Result{}, ErrTypeBookSanitizationFailed.Wrap(err, "failed to process content")
	}
	fonts := content.Fonts
	if err := s.deps.FontPolicy.Validate(fonts); err != nil {
		return CreateBookChapterCommand_Result{}, err
	}
	tx, err := s.deps.DB.Begin(ctx)
	if err != nil {
		return CreateBookChapterCommand_Result{}, apperror.WrapUnexpectedDBError(err)
	}
	queries := s.queries.WithTx(tx)
	lastOrder, err := queries.Book_GetLastChapterOrder(ctx, input.BookID)
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return CreateBookChapterCommand_Result{}, err
	}

	id := GenID()
	err = queries.Book_InsertChapter(ctx, store.Book_InsertChapterParams{
		ID:                id,
		BookID:            input.BookID,
		Name:              input.Name,
		CreatedAt:         timeToTimestamptz(time.Now()),
		Content:           content.Sanitized,
		Order:             lastOrder + 1,
		Words:             content.Words,
		Summary:           input.Summary,
		Fonts:             fonts,
		IsPubliclyVisible: input.IsPubliclyVisible,
	})
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return CreateBookChapterCommand_Result{}, err
	}
	err = queries.RecalculateBookStats(ctx, input.BookID)
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return CreateBookChapterCommand_Result{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return CreateBookChapterCommand_Result{}, apperror.WrapUnexpectedDBError(err)
	}
	s.deps.BookReindexService.ScheduleReindex(ctx, input.BookID)
	return CreateBookChapterCommand_Result{ID: id}, nil
}

func (s *bookManagerService) ReorderChapters(ctx context.Context, input ReorderChaptersCommand) error {

	var (
		oldChapterOrder map[int64]int
	)

	tx, err := s.deps.DB.Begin(ctx)
	if err != nil {
		return err
	}
	queries := s.queries.WithTx(tx)

	{
		chapterOrders, err := queries.Book_GetChapterOrder(ctx, input.BookID)
		if err != nil {
			dal.RollbackTx(ctx, tx)
			return err
		}

		oldChapterOrder = make(map[int64]int, len(chapterOrders))

		for _, v := range chapterOrders {
			oldChapterOrder[v.ID] = int(v.Order)
		}
	}

	var (
		newChapterOrder = make(map[int64]int)
	)

	{
		newChapterOrder = make(map[int64]int, len(oldChapterOrder))

		for i, chapterID := range input.ChapterIDs {
			if _, ok := oldChapterOrder[chapterID]; !ok {
				dal.RollbackTx(ctx, tx)
				return fmt.Errorf("chapter %d does not exist", chapterID)
			}
			if _, ok := newChapterOrder[chapterID]; ok {
				dal.RollbackTx(ctx, tx)
				return fmt.Errorf("chapter %d is duplicated", chapterID)
			}

			newChapterOrder[chapterID] = i + 1
		}
	}

	if len(newChapterOrder) < len(oldChapterOrder) {
		dal.RollbackTx(ctx, tx)
		return errors.New("not enough chapters provided")
	}

	for chapterID, newOrder := range newChapterOrder {
		if oldChapterOrder[chapterID] == newOrder {
			continue
		}
		err = queries.Book_SetChapterOrder(ctx, store.Book_SetChapterOrderParams{
			ID:    chapterID,
			Order: int32(newOrder),
		})
		if err != nil {
			dal.RollbackTx(ctx, tx)
			return err
		}
	}

	err = tx.Commit(ctx)

	return err
}

func (s *bookManagerService) GetBookChapters(ctx context.Context, query ManagerGetBookChaptersQuery) (ManagerGetBookChaptersQuery_Result, error) {
	rows, err := s.queries.GetAllBookChapters(ctx, query.BookID)
	if err != nil {
		return ManagerGetBookChaptersQuery_Result{}, err
	}

	var (
		chapters = make([]ManagerBookChapterDto, len(rows))
	)

	for i, row := range rows {
		chapters[i] = ManagerBookChapterDto{
			ID:                row.ID,
			Name:              row.Name,
			Summary:           row.Summary,
			CreatedAt:         row.CreatedAt.Time,
			Words:             int(row.Words),
			IsPubliclyVisible: row.IsPubliclyVisible,
			Order:             row.Order,
			ScheduledAt:       timeNullableDbToDomain(row.ScheduledAt),
		}
	}

	return ManagerGetBookChaptersQuery_Result{
		Chapters: chapters,
	}, nil
}

func (s *bookManagerService) GetChapter(ctx context.Context, query ManagerGetChapterQuery) (ManagerGetChapterQuery_Result, error) {
	chapter, err := s.queries.GetBookChapterWithDetails(ctx, store.GetBookChapterWithDetailsParams{
		ID:     query.ChapterID,
		BookID: query.BookID,
	})
	if err != nil {
		return ManagerGetChapterQuery_Result{}, err
	}

	return ManagerGetChapterQuery_Result{
		Chapter: ManagerBookChapterDetailsDto{
			ID:                chapter.ID,
			Name:              chapter.Name,
			Summary:           chapter.Summary,
			CreatedAt:         chapter.CreatedAt.Time,
			Words:             int(chapter.Words),
			Order:             chapter.Order,
			Content:           chapter.Content,
			Fonts:             chapter.Fonts,
			IsPubliclyVisible: true,
		},
	}, nil
}

func (s *bookManagerService) UpdateBookChapter(ctx context.Context, cmd UpdateBookChapterCommand) error {
	err := validateChapterName(cmd.Name)
	if err != nil {
		return err
	}

	summary, err := s.deps.Markup.Clean(cmd.Summary)
	if err != nil {
		return ErrTypeBookSanitizationFailed.Wrap(err, "failed to process chapter summary")
	}

	err = validateBookSummary(summary.Sanitized)
	if err != nil {
		return err
	}

	rows, err := s.queries.Chapter_UpdateDetails(ctx, store.Chapter_UpdateDetailsParams{
		Name:              cmd.Name,
		Summary:           summary.Sanitized,
		IsPubliclyVisible: cmd.IsPubliclyVisible,
		ChapterID:         cmd.ChapterID,
		BookID:            cmd.BookID,
		UserID:            uuidDomainToDb(cmd.UserID),
	})
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	if rows == 0 {
		return ErrTypeChapterDoesNotExist.New("chapter not found")
	}

	s.recalculateBookStats(ctx, cmd.BookID)
	s.deps.BookReindexService.ScheduleReindex(ctx, cmd.BookID)
	return nil
}

// GetDraft implements BookManagerService.
func (s *bookManagerService) GetDraft(ctx context.Context, query GetDraftQuery) (DraftDto, error) {
	err := s.authorizeDraftAccess(ctx, query.UserID, query.BookID, query.ChapterID, query.DraftID)
	if err != nil {
		return DraftDto{}, err
	}

	draft, err := s.queries.Draft_GetById(ctx, query.DraftID)
	if err != nil {
		if err == store.ErrNoRows {
			return DraftDto{}, ErrDraftNotFound
		}
		return DraftDto{}, apperror.WrapUnexpectedDBError(err)
	}

	user, err := s.deps.UsersService.GetUserSelfData(ctx, uuidDbToDomain(draft.CreatedBy))
	if err != nil {
		return DraftDto{}, apperror.WrapUnexpectedAppError(err)
	}

	return DraftDto{
		ID:          draft.ID,
		ChapterName: draft.ChapterName,
		Content:     draft.Content,
		CreatedAt:   draft.CreatedAt.Time,
		UpdatedAt:   timeNullableDbToDomain(draft.UpdatedAt),
		ScheduledAt: timeNullableDbToDomain(draft.ScheduledAt),
		Chapter: struct {
			ID               int64     "json:\"id,string\""
			ContentUpdatedAt time.Time "json:\"contentUpdatedAt\""
			Fonts            []string  "json:\"fonts\""
		}{
			ID:               draft.ChapterID,
			ContentUpdatedAt: timeDbToDomain(draft.ChapterContentUpdatedAt),
			Fonts:            draft.ChapterFonts,
		},
		CreatedBy: struct {
			ID   uuid.UUID `json:"id"`
			Name string    `json:"name"`
		}{
			ID:   user.ID,
			Name: user.Name,
		},
		Book: struct {
			ID   int64  `json:"id,string"`
			Name string `json:"name"`
		}{
			ID:   draft.BookID,
			Name: draft.BookName,
		},
		IsChapterPubliclyAvailable: draft.IsChapterPubliclyVisible,
	}, nil
}

func (s *bookManagerService) UpdateDraft(ctx context.Context, cmd UpdateDraftCommand) error {
	err := s.authorizeDraftAccess(ctx, cmd.UserID, cmd.BookID, cmd.ChapterID, cmd.DraftID)
	if err != nil {
		return err
	}
	if err := validateChapterName(cmd.Name); err != nil {
		return err
	}

	processedContent, err := s.deps.Markup.Clean(cmd.Content)

	if err != nil {
		return ErrTypeBookSanitizationFailed.Wrap(err, "failed to process content")
	}
	if err := s.deps.FontPolicy.Validate(processedContent.Fonts); err != nil {
		return err
	}

	err = s.queries.Draft_Update(ctx, store.Draft_UpdateParams{
		ID:          cmd.DraftID,
		Content:     processedContent.Sanitized,
		ChapterName: cmd.Name,
		Summary:     cmd.Summary,
		Words:       processedContent.Words,
	})
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	return nil
}

// DeleteDraft implements BookManagerService.
func (s *bookManagerService) DeleteDraft(ctx context.Context, cmd DeleteDraftCommand) error {
	err := s.authorizeDraftAccess(ctx, cmd.UserID, 0, 0, cmd.DraftID)
	if err != nil {
		return err
	}

	err = s.queries.Draft_Delete(ctx, cmd.DraftID)
	if err != nil {
		return err
	}
	return nil
}

// PublishDraft implements BookManagerService.
func (s *bookManagerService) PublishDraft(ctx context.Context, cmd PublishDraftCommand) (err error) {
	ctx, span := startSpan(ctx, "BookManagerService.PublishDraft")
	defer func() { endSpan(span, err) }()

	var (
		bookID int64
	)

	// get the draft and update the chapter
	draft, err := s.queries.Draft_GetById(ctx, cmd.DraftID)
	if err != nil {
		if err == store.ErrNoRows {
			return ErrDraftNotFound
		}
		return apperror.WrapUnexpectedDBError(err)
	}

	err = s.authorizeDraftAccess(ctx, cmd.UserID, draft.BookID, draft.ChapterID, cmd.DraftID)
	if err != nil {
		return err
	}
	content, err := s.deps.Markup.Clean(draft.Content)
	if err != nil {
		return ErrTypeBookSanitizationFailed.Wrap(err, "failed to process content")
	}
	fonts := content.Fonts
	if err := s.deps.FontPolicy.Validate(fonts); err != nil {
		return err
	}

	tx, err := s.deps.DB.Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	// update the chapter and mark draft as published
	queries := s.queries.WithTx(tx)

	isChapterPublic := draft.IsChapterPubliclyVisible

	if cmd.MakePublic {
		isChapterPublic = true
	}

	bookID, err = queries.Chapter_Update(ctx, store.Chapter_UpdateParams{
		ID:                draft.ChapterID,
		Name:              draft.ChapterName,
		Summary:           draft.Summary,
		Content:           content.Sanitized,
		ContentUpdatedAt:  draft.UpdatedAt,
		Fonts:             fonts,
		Words:             draft.Words,
		IsPubliclyVisible: isChapterPublic,
	})
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}

	err = queries.Draft_MarkAsPublished(ctx, cmd.DraftID)
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}

	err = tx.Commit(ctx)

	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	s.recalculateBookStats(ctx, bookID)

	return nil
}

func (s *bookManagerService) ScheduleDraft(ctx context.Context, cmd ScheduleDraftCommand) error {
	if !cmd.ScheduledAt.After(time.Now()) {
		return apperror.ValidationError.New("scheduled publishing time must be in the future")
	}
	if err := s.authorizeDraftAccess(ctx, cmd.UserID, cmd.BookID, cmd.ChapterID, cmd.DraftID); err != nil {
		return err
	}

	tx, err := s.deps.DB.Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	queries := s.queries.WithTx(tx)
	if err = queries.Draft_ClearChapterSchedules(ctx, cmd.ChapterID); err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}
	if err = queries.Draft_Schedule(ctx, store.Draft_ScheduleParams{
		ID:          cmd.DraftID,
		ScheduledAt: timeToTimestamptz(cmd.ScheduledAt),
	}); err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	return nil
}

func (s *bookManagerService) recalculateBookStats(ctx context.Context, bookID int64) {
	err := s.queries.RecalculateBookStats(ctx, bookID)
	if err != nil {
		s.deps.Log.Errorw("failed to recalculate book stats", "err", err, "bookID", bookID)
	}
}

// GetLatestDraft implements BookManagerService.
func (s *bookManagerService) GetLatestDraft(ctx context.Context, cmd GetLatestDraftQuery) (Nullable[int64], error) {
	if err := s.authorizeChapterEdit(ctx, cmd.UserID, cmd.BookID, cmd.ChapterID); err != nil {
		return Null[int64](), err
	}

	draftID, err := s.queries.Draft_GetLatestID(ctx, cmd.ChapterID)
	if err != nil {
		if err == store.ErrNoRows {
			return Null[int64](), nil
		}
		return Nullable[int64]{}, apperror.WrapUnexpectedDBError(err)
	}

	return Value(draftID), nil
}

// CreateDraft implements BookManagerService.
func (s *bookManagerService) CreateDraft(ctx context.Context, cmd CreateDraftCommand) (int64, error) {
	err := s.authorizeChapterEdit(ctx, cmd.UserID, cmd.BookID, cmd.ChapterID)
	if err != nil {
		return 0, err
	}

	chapter, err := s.queries.GetBookChapterWithDetails(ctx, store.GetBookChapterWithDetailsParams{
		ID:     cmd.ChapterID,
		BookID: cmd.BookID,
	})

	if err != nil {
		if err == store.ErrNoRows {
			return 0, ErrTypeChapterDoesNotExist.New("chapter not found")
		}

		return 0, apperror.WrapUnexpectedDBError(err)
	}

	id := GenID()

	err = s.queries.Draft_Insert(ctx, store.Draft_InsertParams{
		ID:          id,
		CreatedBy:   uuidDomainToDb(cmd.UserID),
		ChapterID:   cmd.ChapterID,
		ChapterName: chapter.Name,
		Content:     chapter.Content,
		UpdatedAt:   timeToTimestamptz(time.Now()),
		CreatedAt:   timeToTimestamptz(time.Now()),
	})

	if err != nil {
		return 0, apperror.WrapUnexpectedDBError(err)
	}

	return id, nil
}

// UpdateDraftChapterName implements BookManagerService.
func (s *bookManagerService) UpdateDraftChapterName(ctx context.Context, cmd UpdateDraftChapterNameCommand) error {
	err := s.authorizeDraftAccess(ctx, cmd.UserID, cmd.BookID, cmd.ChapterID, cmd.DraftID)
	if err != nil {
		return err
	}
	if err := validateChapterName(cmd.ChapterName); err != nil {
		return err
	}

	err = s.queries.Draft_UpdateChapterName(ctx, store.Draft_UpdateChapterNameParams{
		ID:          cmd.DraftID,
		ChapterName: cmd.ChapterName,
	})

	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil

}

// UpdateDraftContent implements BookManagerService.
func (s *bookManagerService) UpdateDraftContent(ctx context.Context, cmd UpdateDraftContentCommand) error {
	err := s.authorizeDraftAccess(ctx, cmd.UserID, cmd.BookID, cmd.ChapterID, cmd.DraftID)
	if err != nil {
		return err
	}

	data, err := s.deps.Markup.Clean(cmd.Content)

	if err != nil {
		return err
	}
	if err := s.deps.FontPolicy.Validate(data.Fonts); err != nil {
		return err
	}

	err = s.queries.Draft_UpdateContent(ctx, store.Draft_UpdateContentParams{
		ID:      cmd.DraftID,
		Content: data.Sanitized,
		Words:   data.Words,
	})

	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil
}

func (s *bookManagerService) authorizeDraftAccess(ctx context.Context, userID uuid.UUID, bookID, chapterID, draftID int64) error {
	allowed, err := s.queries.Draft_UserCanAccess(ctx, store.Draft_UserCanAccessParams{
		DraftID:   draftID,
		UserID:    uuidDomainToDb(userID),
		ChapterID: chapterID,
		BookID:    bookID,
	})
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	if !allowed {
		return apperror.ErrGenericForbidden
	}
	return nil
}

func (s *bookManagerService) authorizeChapterEdit(ctx context.Context, userID uuid.UUID, bookID, chapterID int64) error {
	allowed, err := s.queries.Chapter_UserCanEdit(ctx, store.Chapter_UserCanEditParams{
		ChapterID: chapterID,
		BookID:    bookID,
		UserID:    uuidDomainToDb(userID),
	})
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	if !allowed {
		return apperror.ErrGenericForbidden
	}
	return nil
}

func NewBookManagerService(deps BookManagerServiceDeps, db dal.DB) BookManagerService {
	return &bookManagerService{
		queries: store.New(db),
		deps:    deps,
	}
}
