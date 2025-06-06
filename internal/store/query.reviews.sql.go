// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.reviews.sql

package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteRate = `-- name: DeleteRate :exec
DELETE FROM ratings WHERE book_id = $1 AND user_id = $2
`

type DeleteRateParams struct {
	BookID int64
	UserID pgtype.UUID
}

func (q *Queries) DeleteRate(ctx context.Context, arg DeleteRateParams) error {
	_, err := q.db.Exec(ctx, deleteRate, arg.BookID, arg.UserID)
	return err
}

const deleteReview = `-- name: DeleteReview :exec
delete from reviews where user_id = $1 and book_id = $2
`

type DeleteReviewParams struct {
	UserID pgtype.UUID
	BookID int64
}

func (q *Queries) DeleteReview(ctx context.Context, arg DeleteReviewParams) error {
	_, err := q.db.Exec(ctx, deleteReview, arg.UserID, arg.BookID)
	return err
}

const getBookReviews = `-- name: GetBookReviews :many
select reviews.user_id, reviews.book_id, reviews.content, reviews.created_at, reviews.last_updated_at, reviews.likes, ratings.rating, ratings.updated_at as rating_updated_at, users.name as user_name
from reviews
join ratings on ratings.book_id = reviews.book_id and ratings.user_id = reviews.user_id
join users on users.id = reviews.user_id
where reviews.book_id = $1
order by reviews.created_at desc
limit $3 offset $2
`

type GetBookReviewsParams struct {
	BookID int64
	Offset int32
	Limit  int32
}

type GetBookReviewsRow struct {
	UserID          pgtype.UUID
	BookID          int64
	Content         string
	CreatedAt       pgtype.Timestamptz
	LastUpdatedAt   pgtype.Timestamptz
	Likes           int32
	Rating          int16
	RatingUpdatedAt pgtype.Timestamptz
	UserName        string
}

func (q *Queries) GetBookReviews(ctx context.Context, arg GetBookReviewsParams) ([]GetBookReviewsRow, error) {
	rows, err := q.db.Query(ctx, getBookReviews, arg.BookID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBookReviewsRow
	for rows.Next() {
		var i GetBookReviewsRow
		if err := rows.Scan(
			&i.UserID,
			&i.BookID,
			&i.Content,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.Likes,
			&i.Rating,
			&i.RatingUpdatedAt,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBookReviewsDistribution = `-- name: GetBookReviewsDistribution :many
select rating, count(*) as count
from reviews
join ratings on ratings.book_id = reviews.book_id and ratings.user_id = reviews.user_id
where reviews.book_id = $1
group by rating
order by rating
`

type GetBookReviewsDistributionRow struct {
	Rating int16
	Count  int64
}

func (q *Queries) GetBookReviewsDistribution(ctx context.Context, bookID int64) ([]GetBookReviewsDistributionRow, error) {
	rows, err := q.db.Query(ctx, getBookReviewsDistribution, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBookReviewsDistributionRow
	for rows.Next() {
		var i GetBookReviewsDistributionRow
		if err := rows.Scan(&i.Rating, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRate = `-- name: GetRate :one
SELECT rating FROM ratings WHERE book_id = $1 AND user_id = $2
`

type GetRateParams struct {
	BookID int64
	UserID pgtype.UUID
}

func (q *Queries) GetRate(ctx context.Context, arg GetRateParams) (int16, error) {
	row := q.db.QueryRow(ctx, getRate, arg.BookID, arg.UserID)
	var rating int16
	err := row.Scan(&rating)
	return rating, err
}

const getRating = `-- name: GetRating :one
select ratings.user_id, ratings.book_id, ratings.rating, ratings.updated_at
from ratings
where ratings.user_id = $1 and ratings.book_id = $2
`

type GetRatingParams struct {
	UserID pgtype.UUID
	BookID int64
}

func (q *Queries) GetRating(ctx context.Context, arg GetRatingParams) (Rating, error) {
	row := q.db.QueryRow(ctx, getRating, arg.UserID, arg.BookID)
	var i Rating
	err := row.Scan(
		&i.UserID,
		&i.BookID,
		&i.Rating,
		&i.UpdatedAt,
	)
	return i, err
}

const getReview = `-- name: GetReview :one
select reviews.user_id, reviews.book_id, reviews.content, reviews.created_at, reviews.last_updated_at, reviews.likes
from reviews
where reviews.user_id = $1 and reviews.book_id = $2
`

type GetReviewParams struct {
	UserID pgtype.UUID
	BookID int64
}

func (q *Queries) GetReview(ctx context.Context, arg GetReviewParams) (Review, error) {
	row := q.db.QueryRow(ctx, getReview, arg.UserID, arg.BookID)
	var i Review
	err := row.Scan(
		&i.UserID,
		&i.BookID,
		&i.Content,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.Likes,
	)
	return i, err
}

const getReviewAndRating = `-- name: GetReviewAndRating :one
select reviews.user_id, reviews.book_id, reviews.content, reviews.created_at, reviews.last_updated_at, reviews.likes, ratings.rating, ratings.updated_at as rating_updated_at
from reviews
join ratings on ratings.book_id = reviews.book_id and ratings.user_id = reviews.user_id
where reviews.user_id = $1 and reviews.book_id = $2
`

type GetReviewAndRatingParams struct {
	UserID pgtype.UUID
	BookID int64
}

type GetReviewAndRatingRow struct {
	UserID          pgtype.UUID
	BookID          int64
	Content         string
	CreatedAt       pgtype.Timestamptz
	LastUpdatedAt   pgtype.Timestamptz
	Likes           int32
	Rating          int16
	RatingUpdatedAt pgtype.Timestamptz
}

func (q *Queries) GetReviewAndRating(ctx context.Context, arg GetReviewAndRatingParams) (GetReviewAndRatingRow, error) {
	row := q.db.QueryRow(ctx, getReviewAndRating, arg.UserID, arg.BookID)
	var i GetReviewAndRatingRow
	err := row.Scan(
		&i.UserID,
		&i.BookID,
		&i.Content,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.Likes,
		&i.Rating,
		&i.RatingUpdatedAt,
	)
	return i, err
}

const insertOrUpdateRate = `-- name: InsertOrUpdateRate :exec
INSERT INTO ratings (book_id, user_id, rating)
VALUES ($1, $2, $3)
ON CONFLICT (book_id, user_id)
DO UPDATE SET rating = $3
`

type InsertOrUpdateRateParams struct {
	BookID int64
	UserID pgtype.UUID
	Rating int16
}

func (q *Queries) InsertOrUpdateRate(ctx context.Context, arg InsertOrUpdateRateParams) error {
	_, err := q.db.Exec(ctx, insertOrUpdateRate, arg.BookID, arg.UserID, arg.Rating)
	return err
}

const insertOrUpdateReview = `-- name: InsertOrUpdateReview :one
insert into reviews (user_id, book_id, content)
values ($1, $2, $3)
on conflict (user_id, book_id) do update set content = $3
returning user_id, book_id, content, created_at, last_updated_at, likes
`

type InsertOrUpdateReviewParams struct {
	UserID  pgtype.UUID
	BookID  int64
	Content string
}

func (q *Queries) InsertOrUpdateReview(ctx context.Context, arg InsertOrUpdateReviewParams) (Review, error) {
	row := q.db.QueryRow(ctx, insertOrUpdateReview, arg.UserID, arg.BookID, arg.Content)
	var i Review
	err := row.Scan(
		&i.UserID,
		&i.BookID,
		&i.Content,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.Likes,
	)
	return i, err
}

const recalculateBookRating = `-- name: RecalculateBookRating :exec
UPDATE books
SET 
    rating = (SELECT avg(rating::float8) FROM ratings WHERE ratings.book_id = $1), 
    total_ratings = (SELECT COUNT(*) FROM ratings WHERE ratings.book_id = $1), 
    total_reviews = (SELECT COUNT(*) FROM reviews WHERE reviews.book_id = $1)
WHERE id = $1
`

func (q *Queries) RecalculateBookRating(ctx context.Context, bookID int64) error {
	_, err := q.db.Exec(ctx, recalculateBookRating, bookID)
	return err
}
