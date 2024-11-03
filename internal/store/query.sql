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



-- name: InsertBook :exec
insert into books 
(id, name, summary, author_user_id, created_at, age_rating, tag_ids, cached_parent_tag_ids)
values ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateBook :exec
update books
set name = $2, age_rating = $3, tag_ids = $4, cached_parent_tag_ids = $5, summary = $6
where id = $1;

-- name: GetBook :one
select books.*, users.name as author_name
from books
join users on books.author_user_id = users.id
where books.id = $1
limit 1;

-- name: GetUserBooks :many
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

-- name: GetBookChapters :many
select c.id, c.name, c.words, c."order", c.created_at, c.summary
from book_chapters c
where book_id = $1
order by "order";

-- name: GetBookChaptersMinimal :many
select id, name
from book_chapters
where book_id = $1
order by "order";

-- name: GetBookCollections :many
select collections.id, collections.name, collections.books_count as size, collection_books."order" as position, collections.created_at, users.name as user_name, collections.user_id
from collections
join collection_books on collections.id = collection_books.collection_id
join users on collections.user_id = users.id
where collection_books.book_id = $1
order by collections.created_at desc;



-- name: InsertBookChapter :exec
insert into book_chapters
(id, name, book_id, content, "order", created_at, words, summary)
values ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateBookChapter :one
update book_chapters
set name = $2, content = $3, words = $4
where id = $1
returning book_chapters.book_id;

-- name: RecalculateBookStats :exec
update books
set words = stat.words, chapters = stat.chapters
from (select sum(words) as words, count(1) as chapters from book_chapters where book_id = $1) as stat
where books.id = $1;

-- name: GetBookChapterWithDetails :one
select book_chapters.*
from book_chapters
join books on book_chapters.book_id = books.id
join users on users.id = books.author_user_id
where book_chapters.id = $1 and book_chapters.book_id = $2;

-- name: ReorderChapters :exec
update book_chapters
set "order" = c.new_order
from (values ($1::int[])) as v(arr)
join unnest(v.arr) with ordinality as c (value, new_order)
on c.value = book_chapters.id
where book_chapters.book_id = $2;

-- name: GetLastChapterOrder :one
select cast(coalesce(max("order"), 0) as int4) as last_order
from book_chapters
where book_id = $1;

-- name: SearchDefinedTags :many
select *
from defined_tags
where lowercased_name like $1
limit $2;

-- name: SearchDefinedTagsWithType :many
select *
from defined_tags
where lowercased_name like $1 and tag_type = $3
limit $2;

-- name: InsertDefinedTag :exec
insert into defined_tags
(id, name, description, is_spoiler, is_adult, created_at, tag_type, synonym_of)
values ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: InsertDefinedTagEnMasse :copyfrom
insert into defined_tags
(id, name, description, is_spoiler, is_adult, created_at, tag_type, synonym_of)
values ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: DefinedTagsAreInitialized :one
select exists(select 1 from defined_tags) as initialized;

-- name: GetTag :one
select * from defined_tags where id = $1;

-- name: GetTagParent :one
select t0.* 
from defined_tags t0
where t0.id = (select coalesce(t1.synonym_of, t1.id) from defined_tags t1 where t1.id = $1);

-- name: GetTagsByName :many
select * from defined_tags where name = ANY(sqlc.arg('names')::text[]);

-- name: GetTagsByIds :many
select * from defined_tags where id = ANY(sqlc.arg('ids')::int8[]);

-- name: ImportTags :exec
insert into defined_tags
(id, name, description, is_spoiler, is_adult, created_at, tag_type, synonym_of, is_default)
select
    id, name, description, is_spoiler, is_adult, created_at, tag_type, synonym_of, true
from jsonb_to_recordset($1::jsonb)
as json_set (
    id int8, 
    name text, 
    description text, 
    is_spoiler boolean, 
    is_adult boolean, 
    created_at timestamptz, 
    tag_type tag_type, 
    synonym_of int8)
where not exists (select 1 from defined_tags where name = json_set.name);

-- name: RemoveUnusedDefaultTags :exec
delete from defined_tags d
where 
    d.name <> ANY(sqlc.arg(names)::text[]) and
    d.is_default = true and 
    not exists (select 1 from defined_tags where synonym_of = d.id) and
    not exists (select 1 from books where tag_ids @> array[d.id]);

-- -- name: ReplaceArbitraryTagWithDefinedTag :exec
-- update books
-- set 
--     arbitrary_tags = array_remove(arbitrary_tags, $1),
--     tags = case
--         when $2 = any(tags) then tags
--         else array_append(tags, $2)
--     end
-- where id = $3 and arbitrary_tags @> array[$1];
