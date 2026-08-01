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
insert into moderation_logs (id, "time", "type", target_type, target_id, payload, actor_user_id, reason)
values (sqlc.arg('id'), sqlc.arg('time'), sqlc.arg('type'), 'book', sqlc.arg('targetID')::bigint::text,
        sqlc.arg('payload'), sqlc.arg('actorUserID'), sqlc.arg('reason'));

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
    cover = '',
    view = 0,
    rating = null,
    total_reviews = 0,
    total_ratings = 0,
    is_pinned = false,
    author_user_id = $2
where id = $1;
delete from book_view where book_id = $1;


-- name: ModGetBookLogFiltered :many
select moderation_logs.*, users.name as actor_user_name
from moderation_logs
join users on users.id = moderation_logs.actor_user_id
where 
    target_type = 'book' and target_id = sqlc.arg('bookID')::bigint::text and
    ("type" = ANY(CAST(sqlc.arg('actionTypes') as text[])) or sqlc.arg('actionTypes') is null)
order by "time" desc
limit sqlc.arg('limit') offset sqlc.arg('offset');

-- name: ModCountBookLogFiltered :one
select count(*)
from moderation_logs
join users on users.id = moderation_logs.actor_user_id
where 
    target_type = 'book' and target_id = sqlc.arg('bookID')::bigint::text and
    ("type" = ANY(CAST(sqlc.arg('actionTypes') as text[])) or sqlc.arg('actionTypes') is null);
