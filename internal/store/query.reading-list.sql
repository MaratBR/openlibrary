-- name: GetBookReadingListState :one
SELECT rl.*, bc."order" as chapter_order
FROM reading_list rl
LEFT JOIN book_chapters bc ON rl.last_accessed_chapter_id = bc.id
WHERE rl.book_id = $1 and rl.user_id = $2;

-- name: SetBookReadingListStatusAndChapter :one
INSERT INTO reading_list (book_id, user_id, status, last_accessed_chapter_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (book_id, user_id)
DO UPDATE SET status = $3, last_accessed_chapter_id = $4, last_updated_at = now()
RETURNING *;

-- name: SetBookReadingListStatus :one
INSERT INTO reading_list (book_id, user_id, status, last_accessed_chapter_id)
VALUES ($1, $2, $3, null)
ON CONFLICT (book_id, user_id)
DO UPDATE SET status = $3, last_updated_at = now()
RETURNING *;

-- name: SetBookReadingListChapter :exec
INSERT INTO reading_list (book_id, user_id, status, last_accessed_chapter_id)
VALUES ($1, $2, 'reading', $3)
ON CONFLICT (book_id, user_id)
DO UPDATE SET last_accessed_chapter_id = $3, last_updated_at = now();

-- name: GetFirstChapterID :one
select id
from book_chapters
where book_id = $1 and "order" = 0;

-- name: GetLastChapterID :one
select c.id
from book_chapters c
where c.book_id = $1 and "order" = (select max("order") from book_chapters where book_id = $1);

-- name: GetUserLibrary :many
select 
    books.id, books.name, books.has_cover, books.age_rating, 
    reading_list.last_updated_at,
    last_chapter."order" as chapter_order, last_chapter.name as chapter_name, last_chapter.id as chapter_id
from reading_list
join books on reading_list.book_id = books.id
left join book_chapters last_chapter on last_chapter.id = reading_list.last_accessed_chapter_id
where reading_list.user_id = $1 and reading_list.status = $2
order by reading_list.last_updated_at
limit $3;