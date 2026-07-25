-- name: Moderation_GetChapter :one
select * from book_chapters where id = $1;

-- name: Moderation_SetChapterVisibility :exec
update book_chapters set is_publicly_visible = $2, updated_at = now() where id = $1;

-- name: Moderation_SetCommentRemoved :exec
update comments
set deleted_at = case when sqlc.arg('removed')::bool then now() else null end,
    updated_at = now()
where id = $1;

-- name: Moderation_AddCommentLog :exec
insert into comment_logs (id, "time", comment_id, action_type, payload, actor_user_id, reason)
values ($1, $2, $3, $4, $5, $6, $7);

-- name: Moderation_SetUserBanned :exec
update users set is_banned = $2 where id = $1;

-- name: Moderation_GetLatestUserBanExpiry :one
select expires_at from user_bans where user_id = $1 order by created_at desc limit 1;

-- name: Moderation_AddUserBan :exec
insert into user_bans (id, user_id, created_at, banned_by_user_id, note, expires_at)
values ($1, $2, $3, $4, $5, $6);

-- name: Moderation_RenameUser :exec
update users set name = $2 where id = $1;

-- name: Moderation_ChangeUserAbout :exec
update users set about = $2 where id = $1;

-- name: Moderation_AddUserLog :exec
insert into user_logs (id, user_id, actor_user_id, action_type, payload, "time", reason)
values ($1, $2, $3, $4, $5, $6, $7);

-- name: Moderation_GetUserLoginHistory :many
select id, user_id, created_at, user_agent, ip_address
from sessions
where user_id = $1
order by created_at desc;
