-- name: Moderation_GetAuditLog :many
select moderation_logs.*, users.name as actor_user_name
from moderation_logs
left join users on users.id = moderation_logs.actor_user_id
where (sqlc.arg('target_type')::text = '' or moderation_logs.target_type = sqlc.arg('target_type'))
  and (sqlc.arg('target_id')::text = '' or target_id = sqlc.arg('target_id'))
  and (sqlc.arg('actor_id')::text = '' or actor_user_id::text = sqlc.arg('actor_id'))
  and (sqlc.arg('action')::text = '' or type = sqlc.arg('action'))
  and (sqlc.narg('date_from')::timestamptz is null or time >= sqlc.narg('date_from'))
  and (sqlc.narg('date_to')::timestamptz is null or time < sqlc.narg('date_to'))
order by moderation_logs.time desc, moderation_logs.id desc
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: Moderation_CountAuditLog :one
select count(*) from moderation_logs
where (sqlc.arg('target_type')::text = '' or target_type = sqlc.arg('target_type'))
  and (sqlc.arg('target_id')::text = '' or target_id = sqlc.arg('target_id'))
  and (sqlc.arg('actor_id')::text = '' or actor_user_id::text = sqlc.arg('actor_id'))
  and (sqlc.arg('action')::text = '' or type = sqlc.arg('action'))
  and (sqlc.narg('date_from')::timestamptz is null or time >= sqlc.narg('date_from'))
  and (sqlc.narg('date_to')::timestamptz is null or time < sqlc.narg('date_to'));
