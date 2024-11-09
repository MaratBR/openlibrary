-- name: InsertBook :exec
insert into books 
(
    id, name, summary, author_user_id, created_at, age_rating, tag_ids, cached_parent_tag_ids,
    is_publicly_visible
)
values 
(
    $1, $2, $3, $4, $5, $6, $7, $8, 
    $9
);

-- name: UpdateBook :exec
update books
set name = $2, age_rating = $3, tag_ids = $4, cached_parent_tag_ids = $5, summary = $6, is_publicly_visible = $7
where id = $1;

-- name: RecalculateBookStats :exec
update books
set words = stat.words, chapters = stat.chapters
from (select sum(words) as words, count(1) as chapters from book_chapters where book_id = $1) as stat
where books.id = $1;

-- name: ManagerGetUserBooks :many
select 
    books.*,
    collections.id as collection_id,
    collections.name as collection_name,
    collection_books."order" as collection_position,
    collections.books_count as collection_size
from books
left join collection_books on books.id = collection_books.book_id
left join collections on collection_books.collection_id = collections.id
where author_user_id = $1
order by books.created_at desc
limit $2 offset $3;

-- name: SetChapterOrder :exec
update book_chapters
set "order" = $2
where id = $1;

-- name: GetChapterOrder :many
select "order", id
from book_chapters
where book_id = $1
order by "order";

-- name: GetLastChapterOrder :one
select cast(coalesce(max("order"), 0) as int4) as last_order
from book_chapters
where book_id = $1;

-- name: InsertBookChapter :exec
insert into book_chapters
(id, name, book_id, content, "order", created_at, words, summary)
values ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateBookChapter :one
update book_chapters
set name = $2, content = $3, words = $4, summary = $5
where id = $1
returning book_chapters.book_id;