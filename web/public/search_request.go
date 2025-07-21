package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/gofrs/uuid"
)

type booksSearchRequest struct {
	Query string

	IncludeUsers []uuid.UUID
	ExcludeUsers []uuid.UUID
	IncludeTags  []int64
	ExcludeTags  []int64

	Words           app.Int32Range
	Chapters        app.Int32Range
	WordsPerChapter app.Int32Range

	Page     uint32
	PageSize uint
}

func getBooksSearchRequest(r *http.Request) (search booksSearchRequest) {
	source := r.URL.Query()

	search.Query = source.Get("q")

	search.Words = olhttp.GetInt32RangeFromQuery(source, "w")
	search.Chapters = olhttp.GetInt32RangeFromQuery(source, "c")
	search.WordsPerChapter = olhttp.GetInt32RangeFromQuery(source, "wc")

	search.IncludeTags = olhttp.GetInt64Array(source, "it")
	search.ExcludeTags = olhttp.GetInt64Array(source, "et")

	search.IncludeUsers = olhttp.GetUUIDArray(source, "iu")
	search.ExcludeUsers = olhttp.GetUUIDArray(source, "eu")

	// pagination and page size
	search.Page = olhttp.GetPage(source, "p")
	pageSize := olhttp.GetInt32FromQuery(source, "ps")
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
