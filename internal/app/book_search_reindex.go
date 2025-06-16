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
)

type BookFullReindexService struct {
	mutex  sync.Mutex
	db     store.DBTX
	cancel context.CancelFunc
	client *elasticsearch.TypedClient
}

func NewBookFullReindexService(db store.DBTX, client *elasticsearch.TypedClient) *BookFullReindexService {
	return &BookFullReindexService{db: db, client: client}
}

func (s *BookFullReindexService) ScheduleReindexAll() error {
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

func (s *BookFullReindexService) reindexAll(ctx context.Context, batchSize int) {
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
				_ = enc.Encode(elasticstore.BookIndex{
					Name:        book.Name,
					Description: book.Summary,
					Rating:      string(book.AgeRating),
					AuthorID:    uuidDbToDomain(book.AuthorUserID).String(),
					Tags:        book.CachedParentTagIds,
					Chapters:    book.Chapters,
					Words:       book.Words,
					WordsPerChapter: int32(getWordsPerChapter(
						int(book.Words),
						int(book.Chapters))),
				})
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
