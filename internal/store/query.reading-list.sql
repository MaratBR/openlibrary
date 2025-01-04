-- name: GetBookReadingListState :one
SELECT *
FROM reading_list
WHERE book_id = $1 and user_id = $2;

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

-- name: GetFirstChapterID :one
select id
from book_chapters
where book_id = $1 and "order" = 0;

-- name: GetLastChapterID :one
select c.id
from book_chapters c
where c.book_id = $1 and "order" = (select max("order") from book_chapters where book_id = $1);