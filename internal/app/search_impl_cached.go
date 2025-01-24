package app

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/cache"
)

type cachedSearchService struct {
	inner SearchService
	cache *cache.Cache
}

// ExplainSearchQuery implements SearchService.
func (c *cachedSearchService) ExplainSearchQuery(ctx context.Context, req BookSearchQuery) (DetailedBookSearchQuery, error) {
	return c.inner.ExplainSearchQuery(ctx, req)
}

// GetBookExtremes implements SearchService.
func (c *cachedSearchService) GetBookExtremes(ctx context.Context) (*BookExtremes, error) {
	return c.inner.GetBookExtremes(ctx)
}

// SearchBooks implements SearchService.
func (c *cachedSearchService) SearchBooks(ctx context.Context, req BookSearchQuery) (*BookSearchResult, error) {
	if !GlobalFeatureFlags.DisableCache {
		return c.inner.SearchBooks(ctx, req)
	}

	cacheKey := getSearchRequestCacheKey(&req)
	result := new(BookSearchResult)
	t := time.Now()
	err := c.cache.GetJSON(ctx, cacheKey, result)
	if err == nil {
		result.Meta.CacheHit = true
		result.Meta.CacheKey = cacheKey
		result.Meta.CacheTookUS = time.Since(t).Microseconds()
		return result, nil
	} else if err != cache.ErrCacheMiss {
		return nil, err
	}
	result, err = c.inner.SearchBooks(ctx, req)
	if err != nil {
		return nil, err
	}

	err = c.cache.PutJSON(ctx, cacheKey, result, time.Now().Add(5*time.Minute))
	if err != nil {
		return nil, err
	}

	return result, err
}

func NewCachedSearchService(inner SearchService, cache *cache.Cache) SearchService {
	return &cachedSearchService{
		inner: inner,
		cache: cache,
	}
}

func writeInt32Range(w io.Writer, r Int32Range) {
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

func getSearchRequestCacheKey(req *BookSearchQuery) string {
	h := sha512.New()
	writeInt32Range(h, req.Words)
	writeInt32Range(h, req.WordsPerChapter)
	writeInt32Range(h, req.Chapters)
	writeInt32Range(h, req.Favorites)

	{
		var buf [8]byte
		binary.BigEndian.PutUint64(buf[:], uint64(req.Page))
		h.Write(buf[:])
		binary.BigEndian.PutUint64(buf[:], uint64(req.PageSize))
		h.Write(buf[:])

		binary.BigEndian.PutUint64(buf[:], uint64(len(req.IncludeUsers)))
		h.Write(buf[:])
		for _, id := range req.IncludeUsers {
			h.Write(id.Bytes())
		}

		binary.BigEndian.PutUint64(buf[:], uint64(len(req.ExcludeUsers)))
		h.Write(buf[:])
		for _, id := range req.ExcludeUsers {
			h.Write(id.Bytes())
		}

		binary.BigEndian.PutUint64(buf[:], uint64(len(req.IncludeTags)))
		h.Write(buf[:])
		for _, id := range req.IncludeTags {
			binary.BigEndian.PutUint64(buf[:], uint64(id))
			h.Write(buf[:])
		}

		binary.BigEndian.PutUint64(buf[:], uint64(len(req.ExcludeTags)))
		h.Write(buf[:])
		for _, id := range req.ExcludeTags {
			binary.BigEndian.PutUint64(buf[:], uint64(id))
			h.Write(buf[:])
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
	hashStr := fmt.Sprintf("search:1:sha512:%s", base64.RawURLEncoding.EncodeToString(hash))
	return hashStr
}
