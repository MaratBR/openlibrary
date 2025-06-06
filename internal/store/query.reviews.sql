-- name: InsertOrUpdateRate :exec
INSERT INTO ratings (book_id, user_id, rating)
VALUES ($1, $2, $3)
ON CONFLICT (book_id, user_id)
DO UPDATE SET rating = $3;

-- name: GetRate :one
SELECT rating FROM ratings WHERE book_id = $1 AND user_id = $2;

-- name: RecalculateBookRating :exec
UPDATE books
SET 
    rating = (SELECT avg(rating::float8) FROM ratings WHERE ratings.book_id = $1), 
    total_ratings = (SELECT COUNT(*) FROM ratings WHERE ratings.book_id = $1), 
    total_reviews = (SELECT COUNT(*) FROM reviews WHERE reviews.book_id = $1)
WHERE id = $1;

-- name: DeleteRate :exec
DELETE FROM ratings WHERE book_id = $1 AND user_id = $2;

-- name: InsertOrUpdateReview :one
insert into reviews (user_id, book_id, content)
values ($1, $2, $3)
on conflict (user_id, book_id) do update set content = $3
returning *;

-- name: DeleteReview :exec
delete from reviews where user_id = $1 and book_id = $2;

-- name: GetReviewAndRating :one
select reviews.*, ratings.rating, ratings.updated_at as rating_updated_at
from reviews
join ratings on ratings.book_id = reviews.book_id and ratings.user_id = reviews.user_id
where reviews.user_id = $1 and reviews.book_id = $2;

-- name: GetReview :one
select reviews.*
from reviews
where reviews.user_id = $1 and reviews.book_id = $2;

-- name: GetRating :one
select ratings.*
from ratings
where ratings.user_id = $1 and ratings.book_id = $2;

-- name: GetBookReviews :many
select reviews.*, ratings.rating, ratings.updated_at as rating_updated_at, users.name as user_name
from reviews
join ratings on ratings.book_id = reviews.book_id and ratings.user_id = reviews.user_id
join users on users.id = reviews.user_id
where reviews.book_id = $1
order by reviews.created_at desc
limit sqlc.arg('limit') offset sqlc.arg('offset');

-- name: GetBookReviewsDistribution :many
select rating, count(*) as count
from reviews
join ratings on ratings.book_id = reviews.book_id and ratings.user_id = reviews.user_id
where reviews.book_id = $1
group by rating
order by rating;