-- name: Moderation_GetAuditLog :many
select moderation_logs.*, users.name as actor_user_name
from moderation_logs
left join users on users.id = moderation_logs.actor_user_id
where sqlc.arg('target_type')::text = '' or moderation_logs.target_type = sqlc.arg('target_type')
order by moderation_logs.time desc, moderation_logs.id desc
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: Moderation_CountAuditLog :one
select count(*) from moderation_logs
where sqlc.arg('target_type')::text = '' or target_type = sqlc.arg('target_type');
