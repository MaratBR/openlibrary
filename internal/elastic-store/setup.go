package elasticstore

import (
	"context"

	"github.com/elastic/go-elasticsearch/v9"
)

func Setup(ctx context.Context, client *elasticsearch.TypedClient) error {
	var err error

	err = createBookIndex(ctx, client)
	if err != nil {
		return err
	}

	return nil
}
