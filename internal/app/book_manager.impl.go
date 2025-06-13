package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/imgconvert"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/minio/minio-go/v7"
)

type bookManagerService struct {
	queries       *store.Queries
	db            DB
	tagsService   TagsService
	usersService  UserService
	uploadService *UploadService
}

const (
	BOOK_COVER_DIRECTORY = "book-covers"
)

func (s *bookManagerService) GetUserBooks(ctx context.Context, input GetUserBooksQuery) (GetUserBooksResult, error) {
	books, err := s.queries.ManagerGetUserBooks(ctx, store.ManagerGetUserBooksParams{
		AuthorUserID: uuidDomainToDb(input.UserID),
		Limit:        int32(input.Limit),
		Offset:       int32(input.Offset),
	})
	if err != nil {
		return GetUserBooksResult{}, err
	}

	userBooks, err := s.aggregateUserBooks(ctx, books)
	if err != nil {
		return GetUserBooksResult{}, err
	}

	return GetUserBooksResult{Books: userBooks}, nil
}

func (s *bookManagerService) GetBook(ctx context.Context, query ManagerGetBookQuery) (ManagerGetBookResult, error) {
	book, err := s.queries.GetBook(ctx, query.BookID)
	if err != nil {
		return ManagerGetBookResult{}, err
	}

	tags, err := s.tagsService.GetTagsByIds(ctx, book.TagIds)
	if err != nil {
		return ManagerGetBookResult{}, err
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
		Chapters:        []BookChapterDto{},
		Author: BookDetailsAuthorDto{
			ID:   authorID,
			Name: book.AuthorName,
		},
		Summary:           book.Summary,
		IsPubliclyVisible: book.IsPubliclyVisible,
		IsBanned:          book.IsBanned,
		Cover:             getBookCoverURL(s.uploadService, book.ID, book.HasCover),
	}

	{
		chapters, err := s.queries.GetBookChapters(ctx, query.BookID)
		if err != nil {
			return ManagerGetBookResult{}, err
		}
		bookDto.Chapters = mapSlice(chapters, func(chapter store.GetBookChaptersRow) BookChapterDto {
			return BookChapterDto{
				ID:        chapter.ID,
				Order:     int(chapter.Order),
				Name:      chapter.Name,
				Words:     int(chapter.Words),
				CreatedAt: chapter.CreatedAt.Time,
				Summary:   chapter.Summary,
			}
		})
	}

	{
		collections, err := s.queries.GetBookCollections(ctx, query.BookID)
		if err != nil {
			return ManagerGetBookResult{}, err
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

	return ManagerGetBookResult{
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

	tags, err := s.tagsService.FindParentTagIds(ctx, input.Tags)
	if err != nil {
		return 0, err
	}

	id := GenID()
	err = s.queries.InsertBook(ctx, store.InsertBookParams{
		ID:                 id,
		Name:               input.Name,
		AuthorUserID:       uuidDomainToDb(input.UserID),
		CreatedAt:          timeToTimestamptz(time.Now()),
		TagIds:             tags.TagIds,
		CachedParentTagIds: tags.ParentTagIds,
		AgeRating:          ageRatingDbValue(input.AgeRating),
		Summary:            input.Summary,
		IsPubliclyVisible:  input.IsPubliclyVisible,
	})
	return id, err
}

func (s *bookManagerService) UpdateBook(ctx context.Context, input UpdateBookCommand) error {
	err := validateBookName(input.Name)
	if err != nil {
		return err
	}

	summaryData, err := ProcessContent(input.Summary)
	if err != nil {
		return err
	}

	err = validateBookSummary(summaryData.Sanitized)
	if err != nil {
		return err
	}

	tags, err := s.tagsService.FindParentTagIds(ctx, input.Tags)
	if err != nil {
		return err
	}

	return s.queries.UpdateBook(ctx, store.UpdateBookParams{
		ID:                 input.BookID,
		Name:               input.Name,
		TagIds:             tags.TagIds,
		CachedParentTagIds: tags.ParentTagIds,
		AgeRating:          ageRatingDbValue(input.AgeRating),
		Summary:            summaryData.Sanitized,
		IsPubliclyVisible:  input.IsPubliclyVisible,
	})
}

// UploadBookCover implements BookManagerService.
func (s *bookManagerService) UploadBookCover(ctx context.Context, input UploadBookCoverCommand) (result UploadBookCoverResult, err error) {
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

	path := fmt.Sprintf("%s/%d.jpeg", BOOK_COVER_DIRECTORY, input.BookID)
	_, err = s.uploadService.Client.PutObject(
		ctx,
		s.uploadService.PublicBucket,
		path,
		bytes.NewReader(imgBytes),
		int64(len(imgBytes)),
		minio.PutObjectOptions{ContentType: "image/jpeg"},
	)
	if err != nil {
		return
	}

	err = s.queries.BookSetHasCover(ctx, store.BookSetHasCoverParams{
		ID:       input.BookID,
		HasCover: true,
	})
	if err != nil {
		return
	}

	result.URL = getBookCoverURL(s.uploadService, input.BookID, true)

	return
}

// UpdateBookChaptersOrder updates the order of chapters in a book.
func (s *bookManagerService) UpdateBookChaptersOrder(ctx context.Context, input UpdateBookChaptersOrders) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	queries := s.queries.WithTx(tx)

	chapters, err := queries.GetChaptersOrder(ctx, input.BookID)
	if err != nil {
		rollbackTx(ctx, tx)
		return err
	}

	isEqualSet := commonutil.ContainsSameAndNoDuplicates(chapters, input.ChapterIDs)
	if !isEqualSet {
		rollbackTx(ctx, tx)
		return ErrTypeChaptersReorder.New("chapters do not match")
	}

	for i, chapterID := range input.ChapterIDs {
		err = queries.UpdateChaptersOrder(ctx, store.UpdateChaptersOrderParams{
			ID:    chapterID,
			Order: int32(i + 1),
		})
		if err != nil {
			rollbackTx(ctx, tx)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *bookManagerService) aggregateUserBooks(ctx context.Context, rows []store.ManagerGetUserBooksRow) ([]ManagerAuthorBookDto, error) {
	var (
		books   []ManagerAuthorBookDto = []ManagerAuthorBookDto{}
		book    ManagerAuthorBookDto
		tagsAgg = newTagsAggregator(s.tagsService)
	)

	for _, row := range rows {
		if row.ID != book.ID {
			if book.ID != 0 {
				books = append(books, book)
			}

			tagsAgg.Add(book.ID, row.TagIds)

			book = ManagerAuthorBookDto{
				ID:                row.ID,
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
				Cover:             getBookCoverURL(s.uploadService, row.ID, row.HasCover),
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
		return []ManagerAuthorBookDto{}, err
	}

	for i := 0; i < len(books); i++ {
		bookTagIDs := tagsAgg.BookTags(books[i].ID)
		if bookTagIDs != nil {
			books[i].Tags = mapSlice(bookTagIDs, func(tagID int64) DefinedTagDto {
				return tags[tagID]
			})
		} else {
			books[i].Tags = []DefinedTagDto{}
		}
	}

	return books, nil
}

func (s *bookManagerService) CreateBookChapter(ctx context.Context, input CreateBookChapterCommand) (CreateBookChapterResult, error) {
	lastOrder, err := s.queries.GetLastChapterOrder(ctx, input.BookID)
	if err != nil {
		return CreateBookChapterResult{}, err
	}

	id := GenID()
	content, err := ProcessContent(input.Content)
	if err != nil {
		return CreateBookChapterResult{}, ErrTypeBookSanitizationFailed.Wrap(err, "failed to process content")
	}
	err = s.queries.InsertBookChapter(ctx, store.InsertBookChapterParams{
		ID:        id,
		BookID:    input.BookID,
		Name:      input.Name,
		CreatedAt: timeToTimestamptz(time.Now()),
		Content:   content.Sanitized,
		Order:     lastOrder + 1,
		Words:     content.Words,
		Summary:   input.Summary,
	})
	if err != nil {
		return CreateBookChapterResult{}, err
	}
	err = s.queries.RecalculateBookStats(ctx, input.BookID)
	if err != nil {
		return CreateBookChapterResult{}, err
	}
	return CreateBookChapterResult{ID: id}, nil
}

func (s *bookManagerService) ReorderChapters(ctx context.Context, input ReorderChaptersCommand) error {

	var (
		oldChapterOrder map[int64]int
	)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	queries := s.queries.WithTx(tx)

	{
		chapterOrders, err := queries.GetChapterOrder(ctx, input.BookID)
		if err != nil {
			rollbackTx(ctx, tx)
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
				rollbackTx(ctx, tx)
				return fmt.Errorf("chapter %d does not exist", chapterID)
			}
			if _, ok := newChapterOrder[chapterID]; ok {
				rollbackTx(ctx, tx)
				return fmt.Errorf("chapter %d is duplicated", chapterID)
			}

			newChapterOrder[chapterID] = i + 1
		}
	}

	if len(newChapterOrder) < len(oldChapterOrder) {
		rollbackTx(ctx, tx)
		return errors.New("not enough chapters provided")
	}

	for chapterID, newOrder := range newChapterOrder {
		if oldChapterOrder[chapterID] == newOrder {
			continue
		}
		err = queries.SetChapterOrder(ctx, store.SetChapterOrderParams{
			ID:    chapterID,
			Order: int32(newOrder),
		})
		if err != nil {
			rollbackTx(ctx, tx)
			return err
		}
	}

	err = tx.Commit(ctx)

	return err
}

func (s *bookManagerService) GetBookChapters(ctx context.Context, query ManagerGetBookChaptersQuery) (ManagerGetBookChapterResult, error) {
	rows, err := s.queries.GetBookChapters(ctx, query.BookID)
	if err != nil {
		return ManagerGetBookChapterResult{}, err
	}

	var (
		chapters = make([]ManagerBookChapterDto, len(rows))
	)

	for i, row := range rows {
		chapters[i] = ManagerBookChapterDto{
			ID:              row.ID,
			Name:            row.Name,
			Summary:         row.Summary,
			CreatedAt:       row.CreatedAt.Time,
			Words:           int(row.Words),
			IsAdultOverride: row.IsAdultOverride,
			Order:           row.Order,
		}
	}

	return ManagerGetBookChapterResult{
		Chapters: chapters,
	}, nil
}

func (s *bookManagerService) GetChapter(ctx context.Context, query ManagerGetChapterQuery) (ManagerGetChapterResult, error) {
	chapter, err := s.queries.GetBookChapterWithDetails(ctx, store.GetBookChapterWithDetailsParams{
		ID:     query.ChapterID,
		BookID: query.BookID,
	})
	if err != nil {
		return ManagerGetChapterResult{}, err
	}

	return ManagerGetChapterResult{
		Chapter: ManagerBookChapterDetailsDto{
			ID:                chapter.ID,
			Name:              chapter.Name,
			Summary:           chapter.Summary,
			CreatedAt:         chapter.CreatedAt.Time,
			Words:             int(chapter.Words),
			IsAdultOverride:   chapter.IsAdultOverride,
			Order:             chapter.Order,
			Content:           chapter.Content,
			IsPubliclyVisible: true,
		},
	}, nil
}

// GetDraft implements BookManagerService.
func (s *bookManagerService) GetDraft(ctx context.Context, query GetDraftQuery) (DraftDto, error) {
	draft, err := s.queries.GetDraftById(ctx, query.DraftID)
	if err != nil {
		if err == store.ErrNoRows {
			return DraftDto{}, ErrDraftNotFound
		}
		return DraftDto{}, wrapUnexpectedDBError(err)
	}

	user, err := s.usersService.GetUserSelfData(ctx, uuidDbToDomain(draft.CreatedBy))
	if err != nil {
		return DraftDto{}, wrapUnexpectedAppError(err)
	}

	return DraftDto{
		ID:          draft.ID,
		ChapterName: draft.ChapterName,
		Content:     draft.Content,
		CreatedAt:   draft.CreatedAt.Time,
		UpdatedAt:   draft.UpdatedAt.Time,
		ChapterID:   draft.ChapterID,
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
	}, nil
}

func (s *bookManagerService) UpdateDraft(ctx context.Context, cmd UpdateDraftCommand) error {
	content, err := ProcessContent(cmd.Content)

	if err != nil {
		return ErrTypeBookSanitizationFailed.Wrap(err, "failed to process content")
	}

	err = s.queries.UpdateDraft(ctx, store.UpdateDraftParams{
		ID:              cmd.DraftID,
		Content:         content.Sanitized,
		ChapterName:     cmd.Name,
		Summary:         cmd.Summary,
		IsAdultOverride: cmd.IsAdultOverride,
		Words:           content.Words,
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	return nil
}

// DeleteDraft implements BookManagerService.
func (s *bookManagerService) DeleteDraft(ctx context.Context, cmd DeleteDraftCommand) error {
	err := s.authorizeDraftDelete(cmd.UserID, cmd.DraftID)
	if err != nil {
		return err
	}

	err = s.queries.DeleteDraft(ctx, cmd.DraftID)
	if err != nil {
		return err
	}
	return nil
}

// PublishDraft implements BookManagerService.
func (s *bookManagerService) PublishDraft(ctx context.Context, cmd PublishDraftCommand) error {

	var (
		bookID int64
	)

	// get the draft and update the chapter
	draft, err := s.queries.GetDraftById(ctx, cmd.DraftID)
	if err != nil {
		if err == store.ErrNoRows {
			return ErrDraftNotFound
		}
		return wrapUnexpectedDBError(err)
	}

	err = s.authorizeDraftPublish(cmd.UserID, draft.ChapterID, cmd.DraftID)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	// update the chapter and mark draft as published
	queries := s.queries.WithTx(tx)

	bookID, err = queries.UpdateBookChapter(ctx, store.UpdateBookChapterParams{
		ID:      draft.ChapterID,
		Name:    draft.ChapterName,
		Summary: draft.Summary,
		Content: draft.Content,
		Words:   draft.Words,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}

	err = queries.MarkDraftAsPublished(ctx, cmd.DraftID)
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}

	err = tx.Commit(ctx)

	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	s.recalculateBookStats(ctx, bookID)

	return nil
}

func (s *bookManagerService) recalculateBookStats(ctx context.Context, bookID int64) {
	err := s.queries.RecalculateBookStats(ctx, bookID)
	if err != nil {
		slog.Error("failed to recalculate book stats", "err", err, "bookID", bookID)
	}
}

// GetLatestDraft implements BookManagerService.
func (s *bookManagerService) GetLatestDraft(ctx context.Context, cmd GetLatestDraftQuery) (Nullable[int64], error) {
	draftID, err := s.queries.GetLatestDraftID(ctx, cmd.ChapterID)
	if err != nil {
		if err == store.ErrNoRows {
			return Null[int64](), nil
		}
		return Nullable[int64]{}, wrapUnexpectedDBError(err)
	}

	return Value(draftID), nil
}

// CreateDraft implements BookManagerService.
func (s *bookManagerService) CreateDraft(ctx context.Context, cmd CreateDraftCommand) (int64, error) {
	err := s.authorizeDraftCreate(cmd.UserID, cmd.ChapterID)
	if err != nil {
		return 0, err
	}

	chapter, err := s.queries.GetBookChapterWithDetails(ctx, store.GetBookChapterWithDetailsParams{
		ID:     cmd.ChapterID,
		BookID: 0,
	})

	if err != nil {
		if err == store.ErrNoRows {
			return 0, ErrTypeChapterDoesNotExist.New("chapter not found")
		}

		return 0, wrapUnexpectedDBError(err)
	}

	id := GenID()

	err = s.queries.InsertDraft(ctx, store.InsertDraftParams{
		ID:          id,
		CreatedBy:   uuidDomainToDb(cmd.UserID),
		ChapterID:   cmd.ChapterID,
		ChapterName: chapter.Name,
		Content:     chapter.Content,
		UpdatedAt:   timeToTimestamptz(time.Now()),
		CreatedAt:   timeToTimestamptz(time.Now()),
	})

	if err != nil {
		return 0, wrapUnexpectedDBError(err)
	}

	return id, nil
}

// UpdateDraftContent implements BookManagerService.
func (s *bookManagerService) UpdateDraftContent(ctx context.Context, cmd UpdateDraftContentCommand) error {

	err := s.authorizeDraftUpdate(cmd.UserID, cmd.ChapterID, cmd.DraftID)
	if err != nil {
		return err
	}

	err = s.queries.UpdateDraftContent(ctx, store.UpdateDraftContentParams{
		ID:      cmd.DraftID,
		Content: cmd.Content,
	})

	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil
}

// authorization stuff here

func (s *bookManagerService) authorizeDraftUpdate(userID uuid.UUID, chapterID, draftID int64) error {
	return nil
}

func (s *bookManagerService) authorizeDraftPublish(userID uuid.UUID, chapterID, draftID int64) error {
	return nil
}

func (s *bookManagerService) authorizeDraftDelete(userID uuid.UUID, draftID int64) error {
	return nil
}

func (s *bookManagerService) authorizeDraftCreate(userID uuid.UUID, chapterID int64) error {
	return nil
}

func NewBookManagerService(db DB, tagsService TagsService, uploadService *UploadService, usersService UserService) BookManagerService {
	return &bookManagerService{queries: store.New(db), tagsService: tagsService, db: db, uploadService: uploadService, usersService: usersService}
}
