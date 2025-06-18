package app

import "github.com/elastic/go-elasticsearch/v9"

type BackgroundServices struct {
	Book        BookBackgroundService
	BookReindex BookReindexService
}

func NewBackgroundServices(db DB, esClient *elasticsearch.TypedClient) *BackgroundServices {
	return &BackgroundServices{
		Book:        NewBookBackgroundService(db),
		BookReindex: NewBookFullReindexService(db, esClient),
	}
}

func (s *BackgroundServices) Stop() {
	s.Book.Stop()
}

func (s *BackgroundServices) Start() error {
	if err := s.Book.Start(); err != nil {
		return err
	}
	return nil
}
