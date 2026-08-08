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

-- name: Report_GetByID :one
select reports.*, users.name as reporter_user_name
from reports
join users on users.id = reports.reporter_user_id
where reports.id = $1;

-- name: Report_GetBookContext :one
select
    reports.book_chapter_id,
    reports.book_excerpt,
    books.id as book_id,
    books.name as title,
    authors.name as author,
    books.cover,
    book_chapters.name as chapter,
    books.age_rating,
    coalesce(array_agg(distinct defined_tags.name) filter (where defined_tags.id is not null), '{}')::text[] as warnings,
    books.is_perm_removed,
    books.is_banned,
    books.is_trashed,
    books.is_publicly_visible,
    books.created_at as book_created_at,
    book_chapters.created_at as chapter_created_at,
    book_chapters.updated_at as chapter_updated_at,
    book_chapters.content_updated_at as chapter_content_updated_at,
    (select count(*) from reports related where related.target_type = 'book' and related.target_id = reports.target_id)::int as related_reports
from reports
join books on books.id::text = reports.target_id
join users authors on authors.id = books.author_user_id
left join book_chapters on book_chapters.id = reports.book_chapter_id and book_chapters.book_id = books.id
left join defined_tags on defined_tags.id = any(books.tag_ids) and defined_tags.tag_type = 'warning'
where reports.id = $1 and reports.target_type = 'book'
group by reports.id, books.id, authors.name, book_chapters.id;

-- name: Report_Search :many
select reports.*, users.name as reporter_user_name
from reports
join users on users.id = reports.reporter_user_id
where
    (sqlc.arg('search')::text = ''
        or reports.number ilike '%' || sqlc.arg('search') || '%'
        or reports.reason ilike '%' || sqlc.arg('search') || '%'
        or reports.description ilike '%' || sqlc.arg('search') || '%'
        or reports.target_id ilike '%' || sqlc.arg('search') || '%')
    and (sqlc.arg('target_type')::text = '' or reports.target_type = sqlc.arg('target_type'))
order by reports."time" desc, reports.id desc
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: Report_CountSearch :one
select count(*)
from reports
where
    (sqlc.arg('search')::text = ''
        or reports.number ilike '%' || sqlc.arg('search') || '%'
        or reports.reason ilike '%' || sqlc.arg('search') || '%'
        or reports.description ilike '%' || sqlc.arg('search') || '%'
        or reports.target_id ilike '%' || sqlc.arg('search') || '%')
    and (sqlc.arg('target_type')::text = '' or reports.target_type = sqlc.arg('target_type'));

-- name: Report_UserExists :one
select exists(select 1 from users where id = $1);

-- name: Report_BookExists :one
select exists(select 1 from books where id = $1);

-- name: Report_CommentExists :one
select exists(select 1 from comments where id = $1);
