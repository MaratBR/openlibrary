-- name: InsertCollection :exec
insert into collections (id, name, user_id)
values ($1, $2, $3);

-- name: GetUserCollections :many
select *
from collections
where user_id = $3
order by created_at desc
limit $1 offset $2;

-- name: GetCollectionBooks :many
select b.*, cb."order" as order_within_collection
from collection_books cb
join books b on b.id = cb.book_id
order by cb."order";


-- name: GetLatestUserCollections :many
select *
from collections
where user_id = $1
order by last_updated_at desc
limit $2;

-- name: GetBookCollections :many
select c.*
from collections c
join collection_books cb on cb.collection_id = c.id
where c.user_id = $1 and cb.book_id = $2
order by last_updated_at desc;

-- name: AddBookToCollection :exec
insert into collection_books (book_id, collection_id, "order")
values (
    $1, $2,
    (select coalesce(max("order"), 0) + 1 from collection_books where collection_id = $2)    
)
on conflict (book_id, collection_id) do nothing;

-- name: DeleteBookFromCollection :exec
delete from collection_books where book_id = $1 and collection_id = $2;

-- name: GetMaxOrderInCollection :one
select cast(coalesce(max("order"), -1) as int4)
from collection_books
where collection_id = $1;

-- name: GetCollection :one
select *
from collections
where id = $1;


-- name: GetCollections :many
select *
from collections
where id = ANY($1::int8[]);
