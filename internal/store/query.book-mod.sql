-- name: ModGetBookInfo :one
select is_banned, is_shadow_banned, is_perm_removed, name, summary
from books 
where id = $1;

-- name: ModSetBookBanned :exec
update books 
set is_banned = $1
where id = $2;

-- name: ModSetBookShadowBanned :exec
update books 
set is_shadow_banned = $1
where id = $2;


-- name: ModGetBookModState :one
select id, is_banned, is_shadow_banned
from books
where id = $1;

-- name: ModAddBookLog :exec
insert into book_logs (id, "time", book_id, action_type, payload, actor_user_id, reason) values ($1, $2, $3, $4, $5, $6, $7);

-- name: ModGetBookLog :many
select *
from book_logs
where book_id = $1
order by "time" desc
limit $2 offset $3;

-- name: ModGetBookLogOfType :many
select *
from book_logs
where book_id = $1 and action_type = $4
order by "time" desc
limit $2 offset $3;

-- name: ModPermRemoveBook :exec
update books
set 
    is_perm_removed = true,
    name = '[DELETED]',
    summary = '',
    created_at = now(),
    updated_at = now(),
    age_rating = '?',
    is_publicly_visible = false,
    words = 0,
    -- chapters = 0,
    tag_ids = '{}',
    cached_parent_tag_ids = '{}',
    has_cover = false,
    view = 0,
    rating = null,
    total_reviews = 0,
    total_ratings = 0,
    is_pinned = false,
    author_user_id = $2
where id = $1;
delete from book_view where book_id = $1;


-- name: ModGetBookLogFiltered :many
select book_logs.*, users.name as actor_user_name
from book_logs
join users on users.id = book_logs.actor_user_id
where 
    book_id = $1 and 
    (action_type = ANY(CAST(sqlc.arg('actionTypes') as book_action_type[])) or sqlc.arg('actionTypes') is null)
order by "time" desc
limit $2 offset $3;

-- name: ModCountBookLogFiltered :one
select count(*)
from book_logs
join users on users.id = book_logs.actor_user_id
where 
    book_id = $1 and 
    (action_type = ANY(CAST(sqlc.arg('actionTypes') as book_action_type[])) or sqlc.arg('actionTypes') is null);