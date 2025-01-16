package app

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type reviewsService struct {
	db                DB
	userService       UserService
	backgroundService BookBackgroundService
}

// GetBookReviewsDistribution implements ReviewsService.
func (r *reviewsService) GetBookReviewsDistribution(ctx context.Context, bookID int64) (GetBookReviewsDistributionResult, error) {
	queries := store.New(r.db)
	rows, err := queries.GetBookReviewsDistribution(ctx, bookID)
	if err != nil {
		return GetBookReviewsDistributionResult{}, wrapUnexpectedDBError(err)
	}

	result := GetBookReviewsDistributionResult{}

	for _, row := range rows {
		result.Distribution[row.Rating-1] = int32(row.Count)
	}

	return result, nil
}

// DeleteReview implements ReviewsService.
func (r *reviewsService) DeleteReview(ctx context.Context, cmd DeleteReviewCommand) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	queries := store.New(r.db).WithTx(tx)
	err = queries.DeleteRate(ctx, store.DeleteRateParams{
		UserID: uuidDomainToDb(cmd.UserID),
		BookID: cmd.BookID,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}

	err = queries.DeleteReview(ctx, store.DeleteReviewParams{
		UserID: uuidDomainToDb(cmd.UserID),
		BookID: cmd.BookID,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	r.backgroundService.ScheduleBookRecalculation(cmd.BookID)

	return nil
}

// GetReview implements ReviewsService.
func (r *reviewsService) GetReview(ctx context.Context, query GetReviewQuery) (RatingAndReview, error) {
	queries := store.New(r.db)

	rating, err := queries.GetRating(ctx, store.GetRatingParams{
		UserID: uuidDomainToDb(query.UserID),
		BookID: query.BookID,
	})
	if err != nil {
		if err == store.ErrNoRows {
			return RatingAndReview{}, nil
		}
		return RatingAndReview{}, wrapUnexpectedDBError(err)
	}

	review, err := queries.GetReview(ctx, store.GetReviewParams{
		UserID: uuidDomainToDb(query.UserID),
		BookID: query.BookID,
	})
	if err != nil {
		if err == store.ErrNoRows {
			return RatingAndReview{
				Rating: Value(CreateRatingValue(rating.Rating)),
			}, nil
		}
		return RatingAndReview{}, wrapUnexpectedDBError(err)
	}

	user, err := r.userService.GetUserSelfData(ctx, query.UserID)
	if err != nil {
		return RatingAndReview{}, wrapUnexpectedAppError(err)
	}

	reviewDto := ReviewDto{
		User: ReviewUserDto{
			ID:     uuidDbToDomain(review.UserID),
			Name:   user.Name,
			Avatar: getUserAvatar(user.Name, 84),
		},
		Rating:    CreateRatingValue(rating.Rating),
		Content:   review.Content,
		CreatedAt: review.CreatedAt.Time,
		UpdatedAt: timeNullableDbToDomain(review.LastUpdatedAt),
		Likes:     review.Likes,
	}
	return RatingAndReview{
		Rating: Value(CreateRatingValue(rating.Rating)),
		Review: Value(reviewDto),
	}, nil
}

func (r *reviewsService) getDtoFromRow(ctx context.Context, userID uuid.UUID, bookID int64) (ReviewDto, error) {
	queries := store.New(r.db)
	review, err := queries.GetReviewAndRating(ctx, store.GetReviewAndRatingParams{
		UserID: uuidDomainToDb(userID),
		BookID: bookID,
	})
	if err != nil {
		if err == store.ErrNoRows {
			return ReviewDto{}, err
		}
		return ReviewDto{}, wrapUnexpectedDBError(err)
	}

	user, err := r.userService.GetUserSelfData(ctx, userID)
	if err != nil {
		return ReviewDto{}, wrapUnexpectedAppError(err)
	}

	return ReviewDto{
		User: ReviewUserDto{
			ID:     uuidDbToDomain(review.UserID),
			Name:   user.Name,
			Avatar: getUserAvatar(user.Name, 84),
		},
		Rating:    CreateRatingValue(review.Rating),
		Content:   review.Content,
		CreatedAt: review.CreatedAt.Time,
		UpdatedAt: timeNullableDbToDomain(review.LastUpdatedAt),
		Likes:     review.Likes,
	}, nil
}

// GetBookReviews implements ReviewsService.
func (r *reviewsService) GetBookReviews(ctx context.Context, query GetBookReviewsQuery) (GetBookReviewsResult, error) {
	queries := store.New(r.db)

	pagination := DefaultPaginationRestrictions.Validate(query.Page, query.PageSize)

	reviews, err := queries.GetBookReviews(ctx, store.GetBookReviewsParams{
		BookID: query.BookID,
		Limit:  pagination.Limit(),
		Offset: pagination.Offset(),
	})
	if err != nil {
		return GetBookReviewsResult{}, wrapUnexpectedDBError(err)
	}

	reviewsDto := make([]ReviewDto, 0, len(reviews))
	for _, review := range reviews {
		reviewsDto = append(reviewsDto, ReviewDto{
			User: ReviewUserDto{
				ID:     uuidDbToDomain(review.UserID),
				Name:   review.UserName,
				Avatar: getUserAvatar(review.UserName, 84),
			},
			Rating:    CreateRatingValue(review.Rating),
			Content:   review.Content,
			CreatedAt: review.CreatedAt.Time,
			UpdatedAt: timeNullableDbToDomain(review.LastUpdatedAt),
			Likes:     review.Likes,
		})
	}

	return GetBookReviewsResult{
		Reviews:    reviewsDto,
		Pagination: pagination,
	}, nil
}

// UpdateReview implements ReviewsService.
func (r *reviewsService) UpdateReview(ctx context.Context, cmd UpdateReviewCommand) (ReviewDto, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return ReviewDto{}, wrapUnexpectedDBError(err)
	}

	queries := store.New(r.db).WithTx(tx)
	_, err = queries.InsertOrUpdateReview(ctx, store.InsertOrUpdateReviewParams{
		UserID:  uuidDomainToDb(cmd.UserID),
		BookID:  cmd.BookID,
		Content: cmd.Content,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return ReviewDto{}, wrapUnexpectedDBError(err)
	}

	err = queries.InsertOrUpdateRate(ctx, store.InsertOrUpdateRateParams{
		UserID: uuidDomainToDb(cmd.UserID),
		BookID: cmd.BookID,
		Rating: cmd.Rating.ToUint16(),
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return ReviewDto{}, wrapUnexpectedDBError(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return ReviewDto{}, wrapUnexpectedDBError(err)
	}

	r.backgroundService.ScheduleBookRecalculation(cmd.BookID)

	reviewDto, err := r.getDtoFromRow(ctx, cmd.UserID, cmd.BookID)
	return reviewDto, err
}

// UpdateRating implements ReviewsService.
func (r *reviewsService) UpdateRating(ctx context.Context, cmd UpdateRatingCommand) error {
	queries := store.New(r.db)
	err := queries.InsertOrUpdateRate(ctx, store.InsertOrUpdateRateParams{
		UserID: uuidDomainToDb(cmd.UserID),
		BookID: cmd.BookID,
		Rating: cmd.Rating.ToUint16(),
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	return nil
}

func NewReviewsService(db DB, userService UserService, backgroundService BookBackgroundService) ReviewsService {
	return &reviewsService{db: db, userService: userService, backgroundService: backgroundService}
}
