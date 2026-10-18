package public

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"uuid"
)

type booksSearchRequest struct {
	app.BookSearchFilters
	Query string
	Sort  app.BookSearchSort

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

type bookSearchRequestParsingOptions struct {
	DisableUsers bool
}

func getBooksSearchRequest(r *http.Request, opts *bookSearchRequestParsingOptions) (search booksSearchRequest) {
	source := r.URL.Query()

	search.Query = source.Get("q")
	search.Sort = app.ParseBookSearchSort(source.Get("sort"))

	search.Words = olhttp.GetInt32RangeFromQuery(source, "w")
	search.Chapters = olhttp.GetInt32RangeFromQuery(source, "c")
	search.WordsPerChapter = olhttp.GetInt32RangeFromQuery(source, "wc")

	search.IncludeTags = olhttp.GetInt64Array(source, "it")
	search.ExcludeTags = olhttp.GetInt64Array(source, "et")

	if opts == nil || !opts.DisableUsers {
		search.IncludeUsers = olhttp.GetUUIDArray(source, "iu")
		search.ExcludeUsers = olhttp.GetUUIDArray(source, "eu")
	}

	// pagination and page size
	search.Page = olhttp.GetPage(source, "p")
	search.PageSize = uint(olhttp.GetPageSize(source, "ps", 1, 100, 20))
	search.Rating = source.Get("rating")
	switch search.Rating {
	case "?", "G", "PG", "PG-13", "R", "NC-17":
	default:
		search.Rating = ""
	}
	for key, target := range map[string]*app.SearchNumberRange{
		"views": &search.Views, "readers": &search.Readers, "reviews": &search.Reviews,
		"popularity": &search.Popularity, "weightedRating": &search.WeightedRating,
	} {
		for suffix, bound := range map[string]*app.Nullable[float64]{"min": &target.Min, "max": &target.Max} {
			v, err := strconv.ParseFloat(source.Get(key+"."+suffix), 64)
			if err == nil && !math.IsNaN(v) && !math.IsInf(v, 0) && v >= 0 {
				*bound = app.Value(v)
			}
		}
	}
	for key, target := range map[string]*app.SearchDateRange{"createdAt": &search.CreatedAt, "updatedAt": &search.UpdatedAt} {
		for suffix, bound := range map[string]*string{"min": &target.Min, "max": &target.Max} {
			v := source.Get(key + "." + suffix)
			if _, err := time.Parse("2006-01-02", v); err == nil {
				*bound = v
			}
		}
	}
	return
}
