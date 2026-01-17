package app

import (
	"errors"

	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func getOpenSearchError(r *opensearchapi.Response) error {
	if !r.IsError() {
		return nil
	}

	return errors.New("TODO: make normal error for OS")
}
