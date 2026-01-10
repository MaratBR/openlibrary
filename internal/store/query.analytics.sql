-- name: Analytics_IncrView :exec
insert into ol_analytics.view_counter ("period", book_id, view_count, entity_type, entity_id)
values (sqlc.arg('period'), sqlc.arg('book_id'), sqlc.arg('incr_by'), sqlc.arg('entity_type'), sqlc.arg('entity_id'))
on conflict (period, entity_type, entity_id)
do update set view_count = EXCLUDED.view_count + ol_analytics.view_bucket.view_count;

-- name: Analytics_GetViews :many
select "period", view_count
from ol_analytics.view_counter
where 
    book_id = sqlc.arg('book_id') and entity_type = 0 and (
        "period" = 0 or
        "period" = sqlc.arg('year_period') or
        "period" = sqlc.arg('month_period') or
        "period" = sqlc.arg('week_period') or
        "period" = sqlc.arg('day_period') or
        "period" = sqlc.arg('hour_period')
    );

-- name: Analytics_GetChapterViews :many
select "period", sum(view_count) as agg_view_count
from ol_analytics.view_counter
where
    book_id = sqlc.arg('book_id') and entity_type = 1 and (
        "period" = 0 or
        "period" = sqlc.arg('year_period') or
        "period" = sqlc.arg('month_period') or
        "period" = sqlc.arg('week_period') or
        "period" = sqlc.arg('day_period') or
        "period" = sqlc.arg('hour_period')
    )
group by entity_id;

-- name: Analytics_GetTotalViews :one
select view_count
from ol_analytics.view_counter
where book_id = $1 and "period" = 0 and entity_type = 0; 

-- name: Analytics_GetMostViewedBooksByBookViewsOnly :many
select book_id, view_count
from ol_analytics.view_counter
where "period" = $1 and entity_type = 0
order by view_count desc
limit $2;