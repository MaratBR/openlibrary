-- name: SiteConfig_Get :one
-- Assemble the typed application snapshot from independent settings rows.
select jsonb_object_agg(key, value)::jsonb as value
from site_config
having count(*) > 0;

-- name: SiteConfig_Set :exec
-- A single statement keeps settings and their migration version consistent.
insert into site_config (key, value)
select setting.key, setting.value
from jsonb_each(sqlc.arg(settings)::jsonb) setting
on conflict (key) do update set value = excluded.value;
