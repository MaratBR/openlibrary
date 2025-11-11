-- name: SiteConfig_Get :one
select "value" from site_config
where "key" = 'main';

-- name: SiteConfig_Set :exec
insert into site_config ("key", "value")
values ('main', $1)
on conflict ("key") do update set "value" = EXCLUDED."value";