-- name: GetBook :one
select books.*, users.name as author_name
from books
join users on books.author_user_id = users.id
where books.id = $1
limit 1;

-- name: GetPubliclyVisibleBookChapters :many
select c.id, c.name, c.words, c."order", c.created_at, c.summary, c.is_adult_override
from book_chapters c
where book_id = $1 and is_publicly_visible = true
order by "order";

-- name: GetAllBookChapters :many
select c.*, 
  cast(coalesce((select id from drafts where drafts.chapter_id = c.id order by created_at desc limit 1), 0) as int8) as latest_draft_id
from book_chapters c
where book_id = $1
order by "order";

-- name: GetUserBooks :many
select b.*
from books b
where b.author_user_id = $1 and chapters > 0
order by b.is_pinned desc, b.created_at asc
limit $2 offset $3;

-- name: GetBookCollectionData :many
select collections.id, collections.name, collections.books_count as size, collection_books."order" as position, collections.created_at, users.name as user_name, collections.user_id
from collections
join collection_books on collections.id = collection_books.collection_id
join users on collections.user_id = users.id
where collection_books.book_id = $1
order by collections.created_at desc;

-- name: GetBooksCollectionData :many
select collections.id, collections.name, collections.books_count as size, collection_books.book_id, collection_books."order" as position, collections.created_at, users.name as user_name, collections.user_id
from collections
join collection_books on collections.id = collection_books.collection_id
join users on collections.user_id = users.id
where collection_books.book_id = ANY($1::int8[])
order by collections.created_at desc;

-- name: GetBookChapterWithDetails :one
select 
    bc.*,
    coalesce(prev_chapter.id, 0) as prev_chapter_id,
    coalesce(prev_chapter.name, '') as prev_chapter_name,
    coalesce(next_chapter.id, 0) as next_chapter_id,
    coalesce(next_chapter.name, '') as next_chapter_name
from book_chapters bc
left join lateral (
    select id, name
    from book_chapters
    where book_id = bc.book_id
      and "order" < bc."order"
      and is_publicly_visible = true
    order by "order" desc
    limit 1
) prev_chapter on true
left join lateral (
    select id, name
    from book_chapters
    where book_id = bc.book_id
      and "order" > bc."order"
      and is_publicly_visible = true
    order by "order" asc
    limit 1
) next_chapter on true
where bc.id = $1
  and (bc.book_id = $2 or $2 = 0);

-- name: GetTopUserBooks :many
select *
from books
where author_user_id = $1 and is_publicly_visible
order by rating desc limit $2;


-- name: GetRandomPublicBookIDs :many
select id
from books
where is_publicly_visible and age_rating not in ('R', 'NC-17') and not is_banned and chapters > 0
order by random()
limit $1;


-- name: GetChapterBookID :one
select book_id
from book_chapters
where id = $1;

-- name: GetAllBooks :many
select id, name, summary, author_user_id, created_at, age_rating, cached_parent_tag_ids, is_publicly_visible, chapters, words
from books
where id > $1
order by id asc
limit $2;

-- name: GetBookSearchRelatedData :many
select created_at, has_cover, id
from books
where id = any(sqlc.arg(ids)::int8[]);