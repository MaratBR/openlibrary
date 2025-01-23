package public

import (
	"net/http"
	"net/url"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/gofrs/uuid"
)

type searchRequest struct {
	IncludeUsers []uuid.UUID
	ExcludeUsers []uuid.UUID
	IncludeTags  []int64
	ExcludeTags  []int64

	Words           app.Int32Range
	Chapters        app.Int32Range
	WordsPerChapter app.Int32Range
	Favorites       app.Int32Range

	Page     uint
	PageSize uint
}

func parseSearchRequest(source url.Values) (search searchRequest) {
	search.Words = getInt32RangeFromQuery(source, "w")
	search.Favorites = getInt32RangeFromQuery(source, "f")
	search.Chapters = getInt32RangeFromQuery(source, "c")
	search.WordsPerChapter = getInt32RangeFromQuery(source, "wc")

	search.IncludeTags = getInt64Array(source, "it")
	search.ExcludeTags = getInt64Array(source, "et")

	search.IncludeUsers = getUUIDArray(source, "iu")
	search.ExcludeUsers = getUUIDArray(source, "eu")

	// pagination and page size
	search.Page = getPage(source, "p")
	pageSize := getInt32FromQuery(source, "ps")
	if pageSize.Valid {
		if pageSize.Int32 <= 0 {
			search.PageSize = 20
		} else if pageSize.Int32 > 100 {
			search.PageSize = 100
		} else {
			search.PageSize = uint(pageSize.Int32)
		}
	} else {
		search.PageSize = 20
	}

	return
}

func getSearchRequest(r *http.Request) searchRequest {
	query := r.URL.Query()
	req := parseSearchRequest(query)
	return req
}
