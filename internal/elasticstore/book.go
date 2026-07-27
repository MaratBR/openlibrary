package elasticstore

import (
	"context"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/k3a/html2text"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type BookIndex struct {
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Rating            string    `json:"rating"`
	AuthorID          uuid.UUID `json:"authorId"`
	Tags              []int64   `json:"tags"`
	Chapters          int32     `json:"chapters"`
	Words             int32     `json:"words"`
	Slug              string    `json:"slug"`
	WordsPerChapter   int32     `json:"wordsPerChapter"`
	IsPubliclyVisible bool      `json:"isPubliclyVisible"`
	IsTrashed         bool      `json:"isTrashed"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	Views          int64   `json:"views"`
	Popularity     float64 `json:"popularity"`
	ReadersCount   int64   `json:"readersCount"`
	ReviewsCount   int64   `json:"reviewsCount"`
	WeightedRating float64 `json:"weightedRating"`
}

func (c *BookIndex) Normalize() {
	c.Description = html2text.HTML2Text(c.Description)
}

const (
	BOOKS_INDEX_NAME = "books"
)

func createBookIndex(ctx context.Context, client *opensearchapi.Client) error {
	var err error

	_, err = client.Indices.Create(ctx, opensearchapi.IndicesCreateReq{
		Index: BOOKS_INDEX_NAME,
		Body: strings.NewReader(`
{
  "settings": {
    "max_result_window": 100000
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "long"
      },
      "name": {
        "type": "text"
      },
      "description": {
        "type": "text"
      },
      "rating": {
        "type": "keyword"
      },
      "authorId": {
        "type": "keyword"
      },
      "tags": {
        "type": "long"
      },
      "chapters": {
        "type": "integer"
      },
      "words": {
        "type": "integer"
      },
      "wordsPerChapter": {
        "type": "integer"
      }
    }
  }
}
		`),
	})

	if err != nil {
		return err
	}

	return nil
}
