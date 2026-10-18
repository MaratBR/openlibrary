-- name: Notification_Get :one
select * from notifications where key = $1;

-- name: Notification_Upsert :one
insert into notifications (key, user_id, created_at, "type", title, content, metadata)
values ($1, $2, $3, $4, $5, $6, $7)
on conflict (user_id, key) do update set
    created_at = EXCLUDED.created_at,
    content = EXCLUDED.content,
    metadata = EXCLUDED.metadata,
    title = EXCLUDED.title,
    "type" = EXCLUDED."type"
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
select sqlc.arg(key), users.id, now(), sqlc.arg(type), sqlc.arg(title), sqlc.arg(content), sqlc.arg(metadata)
from users
where (sqlc.narg(recipient_id)::uuid is not null and users.id = sqlc.narg(recipient_id)::uuid)
   or (sqlc.narg(recipient_role)::text is not null and users.role::text = sqlc.narg(recipient_role)::text);

-- name: Notification_History :many
select * from notifications
where user_id = $1 and (sqlc.arg(before_id)::bigint = 0 or id < sqlc.arg(before_id)::bigint)
order by id desc
limit 100;

-- name: Notification_Status :one
select
    (select count(*) from notifications n
     where n.user_id = sqlc.arg(user_id)
       and n.created_at > coalesce((select (data #>> '{}')::timestamptz from user_data
                                   where user_id = sqlc.arg(user_id) and key = 'notifications:last-read-at'), '-infinity'::timestamptz)
       and n.created_at <= statement_timestamp())::bigint as unread_count,
    (select (data #>> '{}')::timestamptz from user_data
     where user_id = sqlc.arg(user_id) and key = 'notifications:last-read-at')::timestamptz as last_read_at,
    statement_timestamp()::timestamptz as observed_at;

-- name: Notification_MarkRead :exec
insert into user_data (user_id, key, data)
values (sqlc.arg(user_id), 'notifications:last-read-at',
        to_jsonb(least(sqlc.arg(read_through)::timestamptz, statement_timestamp())))
on conflict (key, user_id) do update
set data = to_jsonb(greatest((user_data.data #>> '{}')::timestamptz,
                           (excluded.data #>> '{}')::timestamptz));
