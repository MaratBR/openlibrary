-- name: Report_Create :one
with next_number as (
    insert into report_number_counters ("day", counter)
    values ((timezone('UTC', statement_timestamp()))::date, 1)
    on conflict ("day") do update
        set counter = report_number_counters.counter + 1
    returning "day", counter
)
insert into reports (number, "time", reporter_user_id, target_type, target_id, reason, description)
select
    'R-' || to_char("day", 'YYYY-MMDD') || '-' || counter::text,
    sqlc.arg('time'), sqlc.arg('reporter_user_id'), sqlc.arg('target_type'),
    sqlc.arg('target_id'), sqlc.arg('reason'), sqlc.arg('description')
from next_number
returning id, number;

-- name: Report_UserExists :one
select exists(select 1 from users where id = $1);

-- name: Report_BookExists :one
select exists(select 1 from books where id = $1);

-- name: Report_CommentExists :one
select exists(select 1 from comments where id = $1);
