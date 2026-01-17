package elasticstore

import (
	"context"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func Setup(ctx context.Context, client *opensearchapi.Client) error {
	var err error

	err = createBookIndex(ctx, client)
	if err != nil {
		return err
	}

	return nil
}
