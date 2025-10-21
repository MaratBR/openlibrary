-- name: Collection_Insert :exec
insert into collections (id, name, slug, user_id)
values ($1, $2, $3, $4);

-- name: Collection_GetByUser :many
select *
from collections
where user_id = $3
order by created_at desc
limit $1 offset $2;

-- name: Collections_CountByUser :one
select count(*)
from collections 
where user_id = $1;

-- name: Collection_GetBooks :many
select b.*, author.name as author_name, cb."order" as order_within_collection
from collection_books cb
join books b on b.id = cb.book_id
join users author on author.id = b.author_user_id
where cb.collection_id = $3
order by cb."order"
limit $1 offset $2;

-- name: Collection_GetRecentByUser :many
select *
from collections
where user_id = $1
order by last_updated_at desc
limit $2;

-- name: Collection_GetByBook :many
select c.*
from collections c
join collection_books cb on cb.collection_id = c.id
where c.user_id = $1 and cb.book_id = $2
order by last_updated_at desc;

-- name: Collection_AddBookToCollection :exec
insert into collection_books (book_id, collection_id, "order")
values (
    $1, $2,
    (select coalesce(max("order"), 0) + 1 from collection_books where collection_id = $2)    
)
on conflict (book_id, collection_id) do nothing;

-- name: Collection_DeleteBookFromCollection :exec
delete from collection_books where book_id = $1 and collection_id = $2;

-- name: Collection_GetMaxOrder :one
select cast(coalesce(max("order"), -1) as int4)
from collection_books
where collection_id = $1;

-- name: Collection_Get :one
select c.*, u.name as user_name
from collections c
join users u on c.user_id = u.id
where c.id = $1;


-- name: Collections_ListByID :many
select *
from collections
where id = ANY($1::int8[]);


-- name: Collection_RecalculateCounter :exec
update collections
set books_count = coalesce((select count(*) from collection_books where collection_id = sqlc.arg('collection_id')), 0)
where id = sqlc.arg('collection_id');