-- name: GetUser :one
select *
from users
where id = $1
limit 1;

-- name: FindUserByUsername :one
select *
from users
where name = $1
limit 1;

-- name: UserExistsByUsername :one
select exists(select 1
from users
where name = $1);

-- name: InsertUser :exec
insert into users
(id, name, password_hash, joined_at)
values ($1, $2, $3, $4);



-- name: InsertSession :exec
insert into sessions
(id, user_id, created_at, user_agent, ip_address, expires_at)
values ($1, $2, $3, $4, $5, $6);

-- name: UpdateSession :exec
update sessions
set expires_at = $1
where id = $2;

-- name: GetSessionInfo :one
select s.*, u.id as user_id, u.name as user_name, u.joined_at as user_joined_at
from sessions s
join users u on s.user_id = u.id
where s.id = $1;

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
