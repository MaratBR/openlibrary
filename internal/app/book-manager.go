package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type BookManagerService struct {
	queries     *store.Queries
	tagsService *TagsService
}

func NewBookManagerService(db store.DBTX, tagsService *TagsService) *BookManagerService {
	return &BookManagerService{queries: store.New(db), tagsService: tagsService}
}

type CreateBookCommand struct {
	Name      string
	UserID    uuid.UUID
	AgeRating AgeRating
	Tags      []string
}

func (s *BookManagerService) CreateBook(ctx context.Context, input CreateBookCommand) (int64, error) {
	id := GenID()
	err := s.queries.InsertBook(ctx, store.InsertBookParams{
		ID:           id,
		Name:         input.Name,
		AuthorUserID: uuidDomainToDb(input.UserID),
		CreatedAt:    timeToTimestamptz(time.Now()),
		Tags:         input.Tags,
		AgeRating:    ageRatingDbValue(input.AgeRating),
	})
	return id, err
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

func (s *BookManagerService) CreateBookChapter(ctx context.Context, input CreateBookChapterCommand) (CreateBookChapterResult, error) {
	lastOrder, err := s.queries.GetLastChapterOrder(ctx, input.BookID)
	if err != nil {
		return CreateBookChapterResult{}, err
	}

	id := GenID()
	content := CleanUpContent(input.Content)
	err = s.queries.InsertBookChapter(ctx, store.InsertBookChapterParams{
		ID:        id,
		BookID:    input.BookID,
		Name:      input.Name,
		CreatedAt: timeToTimestamptz(time.Now()),
		Content:   content,
		Order:     lastOrder + 1,
		Words:     CountWords(content),
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

type UpdateBookChapterCommand struct {
	ID              int64
	Name            string
	Content         string
	IsAdultOverride bool
}

func (s *BookManagerService) UpdateBookChapter(ctx context.Context, input UpdateBookChapterCommand) error {
	content := CleanUpContent(input.Content)
	bookID, err := s.queries.UpdateBookChapter(ctx, store.UpdateBookChapterParams{
		ID:      input.ID,
		Name:    input.Name,
		Content: content,
		Words:   CountWords(content),
	})
	if err != nil {
		return err
	}
	err = s.queries.RecalculateBookStats(ctx, bookID)
	if err != nil {
		return err
	}
	return nil
}

type ManagerGetBookQuery struct {
	ActorUserID uuid.UUID
	BookID      int64
}

type ManagerBookDetailsDto struct {
	ID              int64                `json:"id,string"`
	Name            string               `json:"name"`
	AgeRating       AgeRating            `json:"ageRating"`
	IsAdult         bool                 `json:"isAdult"`
	Tags            []TagDto             `json:"tags"`
	Words           int                  `json:"words"`
	WordsPerChapter int                  `json:"wordsPerChapter"`
	CreatedAt       time.Time            `json:"createdAt"`
	Collections     []BookCollectionDto  `json:"collections"`
	Chapters        []BookChapterDto     `json:"chapters"`
	Author          BookDetailsAuthorDto `json:"author"`
}

type ManagerGetBookResult struct {
	Book ManagerBookDetailsDto
}

func (s *BookManagerService) GetBook(ctx context.Context, query ManagerGetBookQuery) (ManagerGetBookResult, error) {
	book, err := s.queries.GetBook(ctx, query.BookID)
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
		Tags:            aggregateTags(s.tagsService, book.Tags),
		Words:           int(book.Words),
		WordsPerChapter: getWordsPerChapter(int(book.Words), int(book.Chapters)),
		CreatedAt:       book.CreatedAt.Time,
		Collections:     []BookCollectionDto{},
		Chapters:        []BookChapterDto{},
		Author: BookDetailsAuthorDto{
			ID:   authorID,
			Name: book.AuthorName,
		},
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

type GetUserBooksQuery struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

type GetUserBooksResult struct {
	Books []AuthorBookDto
}

func (s *BookManagerService) GetUserBooks(ctx context.Context, input GetUserBooksQuery) (GetUserBooksResult, error) {
	books, err := s.queries.GetUserBooks(ctx, store.GetUserBooksParams{
		AuthorUserID: uuidDomainToDb(input.UserID),
		Limit:        int32(input.Limit),
		Offset:       int32(input.Offset),
	})
	if err != nil {
		return GetUserBooksResult{}, err
	}

	return GetUserBooksResult{Books: s.aggregateUserBooks(books)}, nil
}

func (s *BookManagerService) aggregateUserBooks(rows []store.GetUserBooksRow) []AuthorBookDto {
	books := []AuthorBookDto{}
	var (
		book AuthorBookDto
	)

	for _, row := range rows {
		if row.ID != book.ID {
			if book.ID != 0 {
				books = append(books, book)
			}

			book = AuthorBookDto{
				ID:              row.ID,
				Name:            row.Name,
				CreatedAt:       row.CreatedAt.Time,
				AgeRating:       AgeRatingPG13,
				Tags:            aggregateTags(s.tagsService, row.Tags),
				Words:           int(row.Words),
				Chapters:        int(row.Chapters),
				WordsPerChapter: getWordsPerChapter(int(row.Words), int(row.Chapters)),
				Collections:     []BookCollectionDto{},
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

	return books
}
