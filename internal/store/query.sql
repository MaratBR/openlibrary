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



-- name: InsertBook :exec
insert into books 
(id, name, author_user_id, created_at)
values ($1, $2, $3, $4);

-- name: InsertBookChapter :exec
insert into book_chapters
(id, name, book_id, content, "order", created_at)
values ($1, $2, $3, $4, $5, $6);

-- name: GetBook :one
select *
from books
where id = $1
limit 1;

-- name: GetBookChapters :many
select *
from book_chapters
where book_id = $1
order by "order";

-- name: GetBookChaptersMinimal :many
select id, name
from book_chapters
where book_id = $1
order by "order";