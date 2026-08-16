-- name: Moderation_GetChapter :one
select * from book_chapters where id = $1;

-- name: Moderation_SetChapterVisibility :exec
update book_chapters set is_publicly_visible = $2, updated_at = now() where id = $1;

-- name: Moderation_SetCommentRemoved :exec
update comments
set deleted_at = case when sqlc.arg('removed')::bool then now() else null end,
    updated_at = now()
where id = $1;

-- name: Moderation_AddLog :exec
insert into moderation_logs (id, "time", "type", target_type, target_id, payload, actor_user_id, reason)
values ($1, $2, $3, $4, $5, $6, $7, $8);

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

-- name: Moderation_GetUserLoginHistory :many
select id, user_id, created_at, user_agent, ip_address
from sessions
where user_id = $1
order by created_at desc;

-- name: Moderation_SearchUserLoginHistory :many
select sessions.id, sessions.user_id, users.name as user_name,
       sessions.created_at, sessions.user_agent, sessions.ip_address,
       sessions.expires_at, sessions.is_terminated,
       sessions.location_country, sessions.location_region, sessions.location_city
from sessions
join users on users.id = sessions.user_id
where (cardinality(sqlc.arg('user_ids')::uuid[]) = 0 or sessions.user_id = any(sqlc.arg('user_ids')::uuid[]))
  and (sqlc.arg('search')::text = ''
       or ip_address ilike '%' || sqlc.arg('search') || '%'
       or user_agent ilike '%' || sqlc.arg('search') || '%'
       or location_country ilike '%' || sqlc.arg('search') || '%'
       or location_region ilike '%' || sqlc.arg('search') || '%'
       or location_city ilike '%' || sqlc.arg('search') || '%')
  and (sqlc.narg('date_from')::timestamptz is null or created_at >= sqlc.narg('date_from'))
  and (sqlc.narg('date_to')::timestamptz is null or created_at < sqlc.narg('date_to'))
  and (sqlc.arg('session_status')::text = ''
       or (sqlc.arg('session_status') = 'active' and not is_terminated and expires_at > now())
       or (sqlc.arg('session_status') = 'expired' and not is_terminated and expires_at <= now())
       or (sqlc.arg('session_status') = 'terminated' and is_terminated))
order by sessions.created_at desc, sessions.id desc
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: Moderation_CountUserLoginHistory :one
select count(*)
from sessions
where (cardinality(sqlc.arg('user_ids')::uuid[]) = 0 or sessions.user_id = any(sqlc.arg('user_ids')::uuid[]))
  and (sqlc.arg('search')::text = ''
       or ip_address ilike '%' || sqlc.arg('search') || '%'
       or user_agent ilike '%' || sqlc.arg('search') || '%'
       or location_country ilike '%' || sqlc.arg('search') || '%'
       or location_region ilike '%' || sqlc.arg('search') || '%'
       or location_city ilike '%' || sqlc.arg('search') || '%')
  and (sqlc.narg('date_from')::timestamptz is null or created_at >= sqlc.narg('date_from'))
  and (sqlc.narg('date_to')::timestamptz is null or created_at < sqlc.narg('date_to'))
  and (sqlc.arg('session_status')::text = ''
       or (sqlc.arg('session_status') = 'active' and not is_terminated and expires_at > now())
       or (sqlc.arg('session_status') = 'expired' and not is_terminated and expires_at <= now())
       or (sqlc.arg('session_status') = 'terminated' and is_terminated));

-- name: Moderation_GetRecentLoginSessions :many
select created_at, location_country, location_region, location_city
from sessions
where user_id = $1
order by created_at desc;
