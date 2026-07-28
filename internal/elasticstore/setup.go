package elasticstore

import (
	"context"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func Setup(ctx context.Context, client *opensearchapi.Client) error {
	return RunMigrations(ctx, client, []Migration{
		{
			Name: "0001_create_books_index",
			Run:  createBookIndex,
		},
	})
}
