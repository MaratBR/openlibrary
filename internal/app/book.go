package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type BookService struct {
	queries     *store.Queries
	tagsService *TagsService
}

func NewBookService(db store.DBTX, tagsService *TagsService) *BookService {
	return &BookService{queries: store.New(db), tagsService: tagsService}
}

type AuthorBookDto struct {
	ID              int64               `json:"id,string"`
	Name            string              `json:"name"`
	CreatedAt       time.Time           `json:"createdAt"`
	AgeRating       AgeRating           `json:"ageRating"`
	Tags            []DefinedTagDto     `json:"tags"`
	Words           int                 `json:"words"`
	WordsPerChapter int                 `json:"wordsPerChapter"`
	Chapters        int                 `json:"chapters"`
	Collections     []BookCollectionDto `json:"collections"`
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

type GetBookQuery struct {
	ID          int64
	ActorUserID uuid.NullUUID
}

type BookChapterDto struct {
	ID        int64     `json:"id,string"`
	Order     int       `json:"order"`
	Name      string    `json:"name"`
	Words     int       `json:"words"`
	CreatedAt time.Time `json:"createdAt"`
	Summary   string    `json:"summary"`
}

type BookDetailsAuthorDto struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type BookUserPermissions struct {
	CanEdit bool `json:"canEdit"`
}

type BookDetailsDto struct {
	ID              int64                `json:"id,string"`
	Name            string               `json:"name"`
	AgeRating       AgeRating            `json:"ageRating"`
	IsAdult         bool                 `json:"isAdult"`
	Tags            []DefinedTagDto      `json:"tags"`
	Words           int                  `json:"words"`
	WordsPerChapter int                  `json:"wordsPerChapter"`
	CreatedAt       time.Time            `json:"createdAt"`
	Collections     []BookCollectionDto  `json:"collections"`
	Chapters        []BookChapterDto     `json:"chapters"`
	Author          BookDetailsAuthorDto `json:"author"`
	Permissions     BookUserPermissions  `json:"permissions"`
}

func (s *BookService) GetBook(ctx context.Context, query GetBookQuery) (BookDetailsDto, error) {
	book, err := s.queries.GetBook(ctx, query.ID)
	if err != nil {
		return BookDetailsDto{}, err
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
		Words:           int(book.Words),
		WordsPerChapter: getWordsPerChapter(int(book.Words), int(book.Chapters)),
		CreatedAt:       book.CreatedAt.Time,
		Collections:     []BookCollectionDto{},
		Chapters:        []BookChapterDto{},
		Author: BookDetailsAuthorDto{
			ID:   authorID,
			Name: book.AuthorName,
		},
		Permissions: BookUserPermissions{CanEdit: query.ActorUserID.Valid && authorID == query.ActorUserID.UUID},
	}

	{
		chapters, err := s.queries.GetBookChapters(ctx, query.ID)
		if err != nil {
			return BookDetailsDto{}, err
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

type embedTagsObjectEnvelope struct {
	Version int             `json:"v"`
	Data    json.RawMessage `json:"d"`
}

type ChapterDto struct {
	ID              int64     `json:"id,string"`
	Name            string    `json:"name"`
	Words           int32     `json:"words"`
	Content         string    `json:"content"`
	IsAdultOverride bool      `json:"isAdultOverride"`
	CreatedAt       time.Time `json:"createdAt"`
	Order           int32     `json:"order"`
	Summary         string    `json:"summary"`
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

func (s BookService) GetBookChapter(ctx context.Context, query GetBookChapterQuery) (GetBookChapterResult, error) {
	chapter, err := s.queries.GetBookChapterWithDetails(ctx, store.GetBookChapterWithDetailsParams{
		ID:     query.ChapterID,
		BookID: query.BookID,
	})
	if err != nil {
		return GetBookChapterResult{}, err
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
			},
		},
	}, nil
}
