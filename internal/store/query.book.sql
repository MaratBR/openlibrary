-- name: GetBook :one
select books.*, users.name as author_name
from books
join users on books.author_user_id = users.id
where books.id = $1
limit 1;

-- name: GetBookChapters :many
select c.id, c.name, c.words, c."order", c.created_at, c.summary, c.is_adult_override
from book_chapters c
where book_id = $1
order by "order";

-- name: GetBookCollections :many
select collections.id, collections.name, collections.books_count as size, collection_books."order" as position, collections.created_at, users.name as user_name, collections.user_id
from collections
join collection_books on collections.id = collection_books.collection_id
join users on collections.user_id = users.id
where collection_books.book_id = $1
order by collections.created_at desc;

-- name: GetBooksCollections :many
select collections.id, collections.name, collections.books_count as size, collection_books.book_id, collection_books."order" as position, collections.created_at, users.name as user_name, collections.user_id
from collections
join collection_books on collections.id = collection_books.collection_id
join users on collections.user_id = users.id
where collection_books.book_id = ANY($1::int8[])
order by collections.created_at desc;

-- name: GetBookChapterWithDetails :one
select 
    bc.*,
    prev_chapter.id as prev_chapter_id,
    prev_chapter.name as prev_chapter_name,
    next_chapter.id as next_chapter_id,
    next_chapter.name as next_chapter_name
from book_chapters bc
left join book_chapters prev_chapter on prev_chapter.book_id = bc.book_id and prev_chapter."order" = bc."order" - 1
left join book_chapters next_chapter on next_chapter.book_id = bc.book_id and next_chapter."order" = bc."order" + 1
join books on bc.book_id = books.id
join users on users.id = books.author_user_id
where bc.id = $1 and bc.book_id = $2;

-- name: GetTopUserBooks :many
select *
from books
where author_user_id = $1 and is_publicly_visible
order by favorites desc limit $2;