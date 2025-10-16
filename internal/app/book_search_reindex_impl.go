package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sync"

	elasticstore "github.com/MaratBR/openlibrary/internal/elastic-store"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

type bookReindexService struct {
	mutex  sync.Mutex
	db     store.DBTX
	cancel context.CancelFunc
	client *elasticsearch.TypedClient
}

func NewBookFullReindexService(db store.DBTX, client *elasticsearch.TypedClient) BookReindexService {
	return &bookReindexService{db: db, client: client}
}

func (s *bookReindexService) ScheduleReindexAll() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	go s.reindexAll(ctx, 10000)

	return nil
}

func (s *bookReindexService) ScheduleReindex(ctx context.Context, id int64) {
	go func() {
		err := s.Reindex(ctx, id)
		if err != nil {
			slog.Error("failed to reindex book", "err", err)
		}
	}()
}

func (s *bookReindexService) Reindex(ctx context.Context, id int64) error {
	queries := store.New(s.db)
	book, err := queries.GetBook(ctx, id)
	if err != nil {
		return err
	}

	idx := elasticstore.BookIndex{
		Name:        book.Name,
		Description: book.Summary,
		Rating:      string(book.AgeRating),
		AuthorID:    uuidDbToDomain(book.AuthorUserID),
		Tags:        book.CachedParentTagIds,
		Chapters:    book.Chapters,
		Words:       book.Words,
		WordsPerChapter: int32(getWordsPerChapter(
			int(book.Words),
			int(book.Chapters))),
	}
	idx.Normalize()
	_, err = s.client.Index(elasticstore.BOOKS_INDEX_NAME).Request(idx).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *bookReindexService) reindexAll(ctx context.Context, batchSize int) {
	_, err := s.client.DeleteByQuery(elasticstore.BOOKS_INDEX_NAME).Query(&types.Query{
		MatchAll: types.NewMatchAllQuery(),
	}).Do(ctx)

	if err != nil {
		slog.Error("failed to delete all docs from index", "index", elasticstore.BOOKS_INDEX_NAME, "err", err)
	}

	var cursor int64

	queries := store.New(s.db)
	num := 1

	slog.Debug("reindexing all books", "logger", "BookFullReindexService")

	body := make(map[string]any, batchSize)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		books, err := queries.GetAllBooks(ctx, store.GetAllBooksParams{
			Limit: int32(batchSize),
			ID:    cursor,
		})
		if err != nil {
			slog.Error("failed to get all books", "logger", "BookFullReindexService", "err", err)
			return
		}

		if len(books) == 0 {
			return
		}

		cursor = books[len(books)-1].ID

		clear(body)

		for _, book := range books {
			body[fmt.Sprintf("%d", book.ID)] = book
		}

		var body bytes.Buffer
		{
			enc := json.NewEncoder(&body)
			var header struct {
				Index struct {
					Index string `json:"_index"`
					ID    string `json:"_id"`
				} `json:"index"`
			}
			header.Index.Index = elasticstore.BOOKS_INDEX_NAME
			for _, book := range books {
				header.Index.ID = fmt.Sprintf("%d", book.ID)
				_ = enc.Encode(header)
				idx := elasticstore.BookIndex{
					Name:        book.Name,
					Description: book.Summary,
					Rating:      string(book.AgeRating),
					AuthorID:    uuidDbToDomain(book.AuthorUserID),
					Tags:        book.CachedParentTagIds,
					Chapters:    book.Chapters,
					Words:       book.Words,
					WordsPerChapter: int32(getWordsPerChapter(
						int(book.Words),
						int(book.Chapters))),
				}
				idx.Normalize()
				_ = enc.Encode(idx)
			}
		}

		res, err := esapi.BulkRequest{
			Index:  elasticstore.BOOKS_INDEX_NAME,
			Body:   &body,
			Pretty: true,
		}.Do(ctx, s.client)

		if err != nil {
			slog.Error("failed to index books batch", "logger", "BookFullReindexService", "err", err)
			return
		}

		if res.IsError() {
			resBytes, _ := io.ReadAll(res.Body)
			slog.Error("failed to index books batch", "logger", "BookFullReindexService", "err", string(resBytes))
			return
		}

		slog.Debug("indexed books batch", "num", num, "logger", "BookFullReindexService", "count", len(books))
		num++
	}
}

type dummyBookReindexService struct{}

func NewDummyBookReindexService() BookReindexService {
	return &dummyBookReindexService{}
}

func (s *dummyBookReindexService) ScheduleReindex(ctx context.Context, id int64) {
}

func (s *dummyBookReindexService) Reindex(ctx context.Context, id int64) error {
	return nil
}

func (s *dummyBookReindexService) ScheduleReindexAll() error {
	return nil
}
