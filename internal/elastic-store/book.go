package elasticstore

import (
	"context"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

type BookIndex struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Rating          string  `json:"rating"`
	AuthorID        string  `json:"authorId"`
	Tags            []int64 `json:"tags"`
	Chapters        int32   `json:"chapters"`
	Words           int32   `json:"words"`
	WordsPerChapter int32   `json:"wordsPerChapter"`
}

const (
	BOOKS_INDEX_NAME = "books"
)

func createBookIndex(ctx context.Context, client *elasticsearch.TypedClient) error {
	var err error

	exists, err := client.Indices.Exists(BOOKS_INDEX_NAME).Do(ctx)

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = client.Indices.Create(BOOKS_INDEX_NAME).
		Request(&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"id":              types.NewLongNumberProperty(),
					"name":            types.NewTextProperty(),
					"description":     types.NewTextProperty(),
					"rating":          types.NewKeywordProperty(),
					"authorId":        types.NewKeywordProperty(),
					"tags":            types.NewLongNumberProperty(),
					"chapters":        types.NewIntegerNumberProperty(),
					"words":           types.NewIntegerNumberProperty(),
					"wordsPerChapter": types.NewIntegerNumberProperty(),
				},
			},
		}).
		Do(ctx)
	if err != nil {
		return err
	}

	return nil
}
