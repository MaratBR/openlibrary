package app

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
)

type searchService struct {
	db            store.DBTX
	queries       *store.Queries
	tagsService   TagsService
	uploadService *UploadService
}

func int32RangeToInt4Range(r Int32Range) store.Int4Range {
	return store.Int4Range{
		Min: pgtype.Int4{Int32: r.Min.Int32, Valid: r.Min.Valid},
		Max: pgtype.Int4{Int32: r.Max.Int32, Valid: r.Max.Valid},
	}
}

// SearchBooks implements SearchService.
func (s *searchService) SearchBooks(ctx context.Context, req BookSearchQuery) (*BookSearchResult, error) {
	dbReq, err := s.constructBookSearchRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	start := time.Now()

	result, err := s.performSearch(ctx, dbReq)
	if err != nil {
		return nil, err
	}

	result.TookUS = time.Since(start).Microseconds()
	return result, nil

}

func (s *searchService) GetBookExtremes(ctx context.Context) (*BookExtremes, error) {
	result, err := store.GetBooksFilterExtremes(ctx, s.db, &store.BookSearchFilter{
		IncludeBanned: false,
		IncludeHidden: false,
		IncludeEmpty:  false,
	})
	if err != nil {
		return nil, err
	}

	return &BookExtremes{
		Words:           Int32Range{Min: Int32{Int32: int32(result.WordsMin), Valid: true}, Max: Int32{Int32: int32(result.WordsMax), Valid: true}},
		Chapters:        Int32Range{Min: Int32{Int32: int32(result.ChaptersMin), Valid: true}, Max: Int32{Int32: int32(result.ChaptersMax), Valid: true}},
		WordsPerChapter: Int32Range{Min: Int32{Int32: int32(result.WordsPerChapterMin), Valid: true}, Max: Int32{Int32: int32(result.WordsPerChapterMax), Valid: true}},
		Favorites:       Int32Range{Min: Int32{Int32: int32(result.FavoritesMin), Valid: true}, Max: Int32{Int32: int32(result.FavoritesMax), Valid: true}},
	}, nil
}

func (s *searchService) performSearch(ctx context.Context, dbReq store.BookSearchRequest) (*BookSearchResult, error) {
	books, err := store.SearchBooks(ctx, s.db, dbReq)
	if err != nil {
		return nil, err
	}

	result := new(BookSearchResult)
	result.Cache.Key = getSearchRequestCacheKey(&dbReq)
	result.Cache.Hit = false
	result.Books = make([]BookSearchItem, len(books))
	bookIds := make([]int64, len(books))

	tagsAgg := newTagsAggregator(s.tagsService)

	for i, book := range books {
		bookIds[i] = book.ID
		result.Books[i] = BookSearchItem{
			ID:   book.ID,
			Name: book.Name,
			Author: BookDetailsAuthorDto{
				ID:   uuidDbToDomain(book.AuthorUserID),
				Name: book.AuthorName,
			},
			CreatedAt:       book.CreatedAt.Time,
			AgeRating:       ageRatingFromDbValue(book.AgeRating),
			Words:           int(book.Words),
			Chapters:        int(book.Chapters),
			WordsPerChapter: getWordsPerChapter(int(book.Words), int(book.Chapters)),
			Summary:         book.Summary,
			Favorites:       book.Favorites,
			Cover:           getBookCoverURL(s.uploadService, book.ID, book.HasCover),
		}

		tagsAgg.Add(book.ID, book.TagIds)
	}

	{
		tags, err := tagsAgg.Fetch(ctx)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(result.Books); i++ {
			result.Books[i].Tags = mapSlice(books[i].TagIds, func(id int64) DefinedTagDto { return tags[id] })
		}
	}

	{
		bookCollectionsArr, err := s.queries.GetBooksCollections(ctx, bookIds)
		if err != nil {
			return nil, err
		}
		bookCollectionMap := make(map[int64][]BookCollectionDto, 0)

		for _, bookCollection := range bookCollectionsArr {
			item := BookCollectionDto{
				ID:       bookCollection.ID,
				Name:     bookCollection.Name,
				Position: int(bookCollection.Position),
				Size:     int(bookCollection.Size),
			}

			if _, ok := bookCollectionMap[bookCollection.BookID]; !ok {
				bookCollectionMap[bookCollection.BookID] = []BookCollectionDto{item}
			} else {
				bookCollectionMap[bookCollection.BookID] = append(bookCollectionMap[bookCollection.BookID], item)
			}
		}

		for i := 0; i < len(result.Books); i++ {
			bookCollections, ok := bookCollectionMap[result.Books[i].ID]
			if ok {
				result.Books[i].Collections = bookCollections
			} else {
				result.Books[i].Collections = []BookCollectionDto{}
			}
		}
	}

	return result, nil
}

func (s *searchService) constructBookSearchRequest(ctx context.Context, req BookSearchQuery) (dbReq store.BookSearchRequest, err error) {
	{
		var includeTags BookTags
		includeTags, err = s.tagsService.FindBookTags(ctx, req.IncludeTags)
		if err != nil {
			return
		}
		dbReq.IncludeParentTags = includeTags.ParentTagIds
	}

	{
		var excludeTags BookTags
		excludeTags, err = s.tagsService.FindBookTags(ctx, req.ExcludeTags)
		if err != nil {
			return
		}
		dbReq.ExcludeParentTags = excludeTags.ParentTagIds
	}

	dbReq.Words = int32RangeToInt4Range(req.Words)
	dbReq.WordsPerChapter = int32RangeToInt4Range(req.WordsPerChapter)
	dbReq.Chapters = int32RangeToInt4Range(req.Chapters)
	dbReq.Favorites = int32RangeToInt4Range(req.Favorites)

	dbReq.IncludeAuthors = arrUuidDomainToDb(req.IncludeUsers)
	dbReq.ExcludeAuthors = arrUuidDomainToDb(req.ExcludeUsers)
	dbReq.IncludeBanned = req.IncludeBanned
	dbReq.IncludeHidden = req.IncludeHidden
	dbReq.IncludeEmpty = req.IncludeEmpty
	dbReq.Limit = req.Limit
	dbReq.Offset = req.Offset

	return
}

func NewSearchService(db store.DBTX, tagsService TagsService, uploadService *UploadService) SearchService {
	return &searchService{
		db:            db,
		queries:       store.New(db),
		tagsService:   tagsService,
		uploadService: uploadService,
	}
}

func writeInt4Range(w io.Writer, r store.Int4Range) {
	var bytes [10]byte

	if r.Max.Valid {
		bytes[0] = 1
		binary.BigEndian.PutUint32(bytes[1:], uint32(r.Max.Int32))
	} else {
		bytes[0] = 0
	}

	if r.Min.Valid {
		bytes[5] = 1
		binary.BigEndian.PutUint32(bytes[6:], uint32(r.Min.Int32))
	} else {
		bytes[5] = 0
	}

	w.Write(bytes[:])
}

func getSearchRequestCacheKey(req *store.BookSearchRequest) string {
	h := sha512.New()
	writeInt4Range(h, req.Words)
	// writeInt4Range(h, req.WordsPerChapter)
	// writeInt4Range(h, req.Chapters)
	// writeInt4Range(h, req.Favorites)

	{
		var buf [8]byte
		binary.BigEndian.PutUint64(buf[:], uint64(req.Limit))
		h.Write(buf[:])
		binary.BigEndian.PutUint64(buf[:], uint64(req.Offset))
		h.Write(buf[:])

		binary.BigEndian.PutUint64(buf[:], uint64(len(req.IncludeAuthors)))
		h.Write(buf[:])
		for _, id := range req.IncludeAuthors {
			h.Write(id.Bytes[:])
		}

		binary.BigEndian.PutUint64(buf[:], uint64(len(req.ExcludeAuthors)))
		h.Write(buf[:])
		for _, id := range req.ExcludeAuthors {
			h.Write(id.Bytes[:])
		}

		binary.BigEndian.PutUint64(buf[:], uint64(len(req.IncludeParentTags)))
		h.Write(buf[:])
		for _, id := range req.IncludeParentTags {
			var tagIdBuf [8]byte
			binary.BigEndian.PutUint64(tagIdBuf[:], uint64(id))
			h.Write(tagIdBuf[:])
		}

		binary.BigEndian.PutUint64(buf[:], uint64(len(req.ExcludeParentTags)))
		h.Write(buf[:])
		for _, id := range req.ExcludeParentTags {
			var tagIdBuf [8]byte
			binary.BigEndian.PutUint64(tagIdBuf[:], uint64(id))
			h.Write(tagIdBuf[:])
		}
	}

	{
		var buf2 [3]byte
		if req.IncludeEmpty {
			buf2[0] = 1
		}
		if req.IncludeBanned {
			buf2[1] = 1
		}
		if req.IncludeHidden {
			buf2[2] = 1
		}
		h.Write(buf2[:])
	}

	hash := h.Sum(nil)
	hashStr := fmt.Sprintf("BookSearchRequest:1:sha512:%s", base64.RawURLEncoding.EncodeToString(hash))
	return hashStr
}
