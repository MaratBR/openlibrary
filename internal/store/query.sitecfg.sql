-- name: SiteConfig_All :many
select * from site_config;

-- name: SiteConfig_Set :exec
insert into site_config ("key", "value")
values ($1, $2)
on conflict ("key") do update set "value" = EXCLUDED."value";