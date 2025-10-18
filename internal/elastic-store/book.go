package elasticstore

import (
	"context"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/gofrs/uuid"
	"github.com/k3a/html2text"
)

type BookIndex struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Rating          string    `json:"rating"`
	AuthorID        uuid.UUID `json:"authorId"`
	Tags            []int64   `json:"tags"`
	Chapters        int32     `json:"chapters"`
	Words           int32     `json:"words"`
	Slug            string    `json:"slug"`
	WordsPerChapter int32     `json:"wordsPerChapter"`
}

func (c *BookIndex) Normalize() {
	c.Description = html2text.HTML2Text(c.Description)
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

	maxResultWindow := 100000

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
			Settings: &types.IndexSettings{
				MaxResultWindow: &maxResultWindow,
			},
		}).
		Do(ctx)
	if err != nil {
		return err
	}

	return nil
}
