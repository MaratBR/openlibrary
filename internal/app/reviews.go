package app

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

type RatingValue uint8

func (v RatingValue) ToUint16() int16 {
	return int16(uint8(v))
}

func (v RatingValue) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(v.ToUint16()), 10)), nil
}

func (v *RatingValue) UnmarshalJSON(data []byte) error {
	var b int16
	err := json.Unmarshal(data, &b)
	if err != nil {
		return err
	}
	*v = CreateRatingValue(b)
	return nil
}

func CreateRatingValue(v int16) RatingValue {
	if v < 1 {
		v = 1
	} else if v > 10 {
		v = 10
	}
	return RatingValue(v)
}

type GetBookReviewsQuery struct {
	BookID   int64
	PageSize int32
	Page     int32
}

type ReviewUserDto struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar"`
}

type ReviewDto struct {
	User      ReviewUserDto       `json:"user"`
	Rating    RatingValue         `json:"rating"`
	Content   string              `json:"content"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt Nullable[time.Time] `json:"updatedAt"`
	Likes     int32               `json:"likes"`
}

type GetBookReviewsResult struct {
	Reviews    []ReviewDto
	Pagination PaginationOptions
}

type UpdateReviewCommand struct {
	BookID  int64
	UserID  uuid.UUID
	Rating  RatingValue
	Content string
}

type GetReviewQuery struct {
	UserID uuid.UUID
	BookID int64
}

type DeleteReviewCommand struct {
	UserID uuid.UUID
	BookID int64
}

type BookReviewsDistribution [10]int32

type GetBookReviewsDistributionResult struct {
	Distribution BookReviewsDistribution
}

type ReviewsService interface {
	GetBookReviewsDistribution(ctx context.Context, bookID int64) (GetBookReviewsDistributionResult, error)
	GetBookReviews(ctx context.Context, query GetBookReviewsQuery) (GetBookReviewsResult, error)
	UpdateReview(ctx context.Context, cmd UpdateReviewCommand) (ReviewDto, error)
	GetReview(ctx context.Context, query GetReviewQuery) (Nullable[ReviewDto], error)
	DeleteReview(ctx context.Context, cmd DeleteReviewCommand) error
}
