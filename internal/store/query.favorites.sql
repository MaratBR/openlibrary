-- name: RecalculateBookFavorites :exec
update books
set favorites = coalesce((select count(*) from favorites where book_id = $1 and is_favorite), 0)
where id = $1;

-- name: SetUserFavourite :exec
insert into favorites (user_id, book_id, is_favorite)
values ($1, $2, $3)
on conflict (user_id, book_id) do update set is_favorite = $3;

-- name: IsFavoritedBy :one
select is_favorite
from favorites
where user_id = $1 and book_id = $2;