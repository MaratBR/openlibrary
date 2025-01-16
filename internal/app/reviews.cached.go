package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/MaratBR/openlibrary/internal/app/cache"
)

type cachedReviewsService struct {
	inner ReviewsService
	cache *cache.Cache
}

// UpdateRating implements ReviewsService.
func (c *cachedReviewsService) UpdateRating(ctx context.Context, cmd UpdateRatingCommand) error {
	c.invalidate(ctx, cmd.BookID)
	return c.inner.UpdateRating(ctx, cmd)
}

// GetBookReviewsDistribution implements ReviewsService.
func (c *cachedReviewsService) GetBookReviewsDistribution(ctx context.Context, bookID int64) (GetBookReviewsDistributionResult, error) {
	queryKey := fmt.Sprintf("GetBookReviewsDistribution:%d", bookID)
	return cache.GetOrSet(ctx, c.cache, queryKey, func(entry *cache.CacheEntry) (GetBookReviewsDistributionResult, error) {
		result, err := c.inner.GetBookReviewsDistribution(ctx, bookID)
		return result, err
	})

}

// DeleteReview implements ReviewsService.
func (c *cachedReviewsService) DeleteReview(ctx context.Context, cmd DeleteReviewCommand) error {
	return c.inner.DeleteReview(ctx, cmd)
}

// GetReview implements ReviewsService.
func (c *cachedReviewsService) GetReview(ctx context.Context, query GetReviewQuery) (RatingAndReview, error) {
	return c.inner.GetReview(ctx, query)
}

// GetBookReviews implements ReviewsService.
func (c *cachedReviewsService) GetBookReviews(ctx context.Context, query GetBookReviewsQuery) (GetBookReviewsResult, error) {
	if query.Page == 1 && query.PageSize == 20 {
		queryKey := fmt.Sprintf("GetBookReviews:%d:1:20", query.BookID)
		result, err := cache.GetOrSet(ctx, c.cache, queryKey, func(entry *cache.CacheEntry) (GetBookReviewsResult, error) {
			result, err := c.inner.GetBookReviews(ctx, query)
			return result, err
		})
		return result, err
	}
	result, err := c.inner.GetBookReviews(ctx, query)
	return result, err
}

// UpdateReview implements ReviewsService.
func (c *cachedReviewsService) UpdateReview(ctx context.Context, cmd UpdateReviewCommand) (ReviewDto, error) {
	review, err := c.inner.UpdateReview(ctx, cmd)
	c.invalidate(ctx, cmd.BookID)
	return review, err
}

func (c *cachedReviewsService) invalidate(ctx context.Context, bookID int64) {
	_, err := c.cache.Remove(ctx, fmt.Sprintf("GetBookReviews:%d:1:20", bookID))
	if err != nil {
		slog.Error("failed to remove GetBookReviews from cache", "err", err)
	}
}

func NewCachedReviewsService(inner ReviewsService, cache *cache.Cache) ReviewsService {
	return &cachedReviewsService{
		inner: inner,
		cache: cache,
	}
}
