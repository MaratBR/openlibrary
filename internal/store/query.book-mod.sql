-- name: ModGetBookInfo :one
select books.is_banned, books.is_shadow_banned, books.is_perm_removed, books.name, books.summary,
       books.author_user_id, users.name as author_user_name, books.created_at,
       books.age_rating, books.is_publicly_visible, books.words, books.chapters,
       (select count(*) from reports where target_type = 'book' and target_id = books.id::text)::bigint as reports_count,
       coalesce(latest_report.id, 0)::bigint as latest_pending_report_id,
       coalesce(latest_report.number, '')::text as latest_pending_report_number,
       coalesce(latest_report.reason, '')::text as latest_pending_report_reason,
       latest_report.time as latest_pending_report_time,
       coalesce(latest_ban.reason, '')::text as ban_reason
from books
join users on users.id = books.author_user_id
left join lateral (
    select id, number, reason, time from reports
    where target_type = 'book' and target_id = books.id::text and status = 'unreviewed'
    order by time desc, id desc limit 1
) latest_report on true
left join lateral (
    select reason from moderation_logs
    where target_type = 'book' and target_id = books.id::text and type = 'ban'
    order by time desc limit 1
) latest_ban on books.is_banned
where books.id = $1;

-- name: ModSearchBooks :many
select books.id, books.name, books.created_at, books.is_banned, books.is_shadow_banned,
       books.is_trashed, books.is_perm_removed, books.is_publicly_visible,
       books.words, books.chapters, users.id as author_user_id, users.name as author_user_name,
       (select count(*) from reports where target_type = 'book' and target_id = books.id::text)::bigint as reports_count
from books
join users on users.id = books.author_user_id
where (sqlc.arg('search')::text = ''
       or (sqlc.narg('search_id')::bigint is not null and books.id = sqlc.narg('search_id'))
       or (sqlc.narg('search_id')::bigint is null and
           case when sqlc.arg('exact_name')::bool then lower(books.name) = lower(sqlc.arg('search'))
                else books.name ilike '%' || sqlc.arg('search') || '%' end))
  and (sqlc.arg('include_banned')::bool or (not books.is_banned and not books.is_shadow_banned))
  and (sqlc.arg('include_deleted')::bool or (not books.is_trashed and not books.is_perm_removed))
order by books.name, books.id
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: ModCountBooks :one
select count(*)
from books
where (sqlc.arg('search')::text = ''
       or (sqlc.narg('search_id')::bigint is not null and books.id = sqlc.narg('search_id'))
       or (sqlc.narg('search_id')::bigint is null and
           case when sqlc.arg('exact_name')::bool then lower(books.name) = lower(sqlc.arg('search'))
                else books.name ilike '%' || sqlc.arg('search') || '%' end))
  and (sqlc.arg('include_banned')::bool or (not books.is_banned and not books.is_shadow_banned))
  and (sqlc.arg('include_deleted')::bool or (not books.is_trashed and not books.is_perm_removed));

-- name: ModGetBookChapters :many
select book_chapters.id, book_chapters.name, book_chapters.created_at, book_chapters.updated_at,
       book_chapters.words, book_chapters.is_publicly_visible,
       exists(select 1 from reports where reports.target_type = 'book'
              and reports.target_id = book_chapters.book_id::text
              and reports.book_chapter_id = book_chapters.id and reports.status = 'unreviewed') as has_pending_reports
from book_chapters
where book_chapters.book_id = $1
order by book_chapters."order", book_chapters.id;

-- name: ModChangeBookAgeRating :exec
update books set age_rating = $2 where id = $1;

-- name: ModChangeBookSummary :exec
update books set summary = $2 where id = $1;

-- name: ModSetBookBanned :exec
update books 
set is_banned = $1
where id = $2;

-- name: ModSetBookShadowBanned :exec
update books 
set is_shadow_banned = $1
where id = $2;


-- name: ModGetBookModState :one
select id, is_banned, is_shadow_banned
from books
where id = $1;

-- name: ModAddBookLog :exec
insert into moderation_logs (id, "time", "type", target_type, target_id, payload, actor_user_id, reason)
values (sqlc.arg('id'), sqlc.arg('time'), sqlc.arg('type'), 'book', sqlc.arg('targetID')::bigint::text,
        sqlc.arg('payload'), sqlc.arg('actorUserID'), sqlc.arg('reason'));

-- name: ModPermRemoveBook :exec
update books
set 
    is_perm_removed = true,
    name = '[DELETED]',
    summary = '',
    created_at = now(),
    updated_at = now(),
    age_rating = '?',
    is_publicly_visible = false,
    words = 0,
    -- chapters = 0,
    tag_ids = '{}',
    cached_parent_tag_ids = '{}',
    cover = '',
    view = 0,
    rating = null,
    total_reviews = 0,
    total_ratings = 0,
    is_pinned = false,
    author_user_id = $2
where id = $1;
delete from book_view where book_id = $1;


-- name: ModGetBookLogFiltered :many
select moderation_logs.*, users.name as actor_user_name
from moderation_logs
join users on users.id = moderation_logs.actor_user_id
where 
    target_type = 'book' and target_id = sqlc.arg('bookID')::bigint::text and
    ("type" = ANY(CAST(sqlc.arg('actionTypes') as text[])) or sqlc.arg('actionTypes') is null)
order by "time" desc
limit sqlc.arg('limit') offset sqlc.arg('offset');

-- name: ModCountBookLogFiltered :one
select count(*)
from moderation_logs
join users on users.id = moderation_logs.actor_user_id
where 
    target_type = 'book' and target_id = sqlc.arg('bookID')::bigint::text and
    ("type" = ANY(CAST(sqlc.arg('actionTypes') as text[])) or sqlc.arg('actionTypes') is null);
