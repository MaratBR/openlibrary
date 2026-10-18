-- name: Notification_Get :one
select * from notifications where key = $1;

-- name: Notification_Upsert :one
insert into notifications (key, user_id, created_at, "type", title, content, metadata)
values ($1, $2, $3, $4, $5, $6, $7)
on conflict (user_id, key) do update set
    created_at = EXCLUDED.created_at,
    content = EXCLUDED.content,
    metadata = EXCLUDED.metadata,
    title = EXCLUDED.title
returning id;

-- name: Notification_DeleteByKey :exec
delete from notifications where key = $1;

-- name: Notification_DeleteByUserID :exec
delete from notifications where user_id = sqlc.arg(user_id);

-- name: Notification_QueryAfter :many
select *
from notifications
where user_id = $1 and created_at > $2
order by created_at desc
limit $3;

-- name: Notification_QueryBefore :many
select *
from notifications
where user_id = $1 and created_at < $2
order by created_at desc
limit $3;

-- name: Notification_Query :many
select *
from notifications
where user_id = $1
order by created_at desc
limit $2;
-- name: Notification_Stream :many
select * from notifications
where user_id = $1 and id > $2
order by id asc
limit 100;

-- name: Notification_Send :execrows
insert into notifications (key, user_id, created_at, "type", title, content, metadata)
select sqlc.arg(key), users.id, now(), 'info', sqlc.arg(title), sqlc.arg(content), '{}'::jsonb
from users
where (sqlc.narg(recipient_id)::uuid is not null and users.id = sqlc.narg(recipient_id)::uuid)
   or (sqlc.narg(recipient_role)::text is not null and users.role::text = sqlc.narg(recipient_role)::text);

-- name: Notification_History :many
select * from notifications
where user_id = $1 and (sqlc.arg(before_id)::bigint = 0 or id < sqlc.arg(before_id)::bigint)
order by id desc
limit 100;
