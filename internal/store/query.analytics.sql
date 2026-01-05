-- name: Analytics_IncrView :exec
insert into ol_analytics.view_bucket ("period", book_id, "count")
values ($1, $2, $3)
on conflict ("period", book_id)
do update set "count" = EXCLUDED."count" + ol_analytics.view_bucket."count";

-- name: Analytics_GetViews :many
select "period", count
from ol_analytics.view_bucket
where 
    book_id = sqlc.arg('book_id') and (
        "period" = 0 or
        "period" = sqlc.arg('year_period') or
        "period" = sqlc.arg('month_period') or
        "period" = sqlc.arg('week_period') or
        "period" = sqlc.arg('day_period') or
        "period" = sqlc.arg('hour_period')
    );

-- name: Analytics_GetViewBuckets :many
select "period", count
from ol_analytics.view_bucket
where book_id = $1 and "period" >= sqlc.arg('from') and "period" <= sqlc.arg('to');

-- name: Analytics_GetTotalViews :one
select count
from ol_analytics.view_bucket
where book_id = $1 and "period" = 0; 

-- name: Analytics_GetMostViewedBooks :many
select book_id, count
from ol_analytics.view_bucket
where "period" = $1
order by count desc
limit $2;