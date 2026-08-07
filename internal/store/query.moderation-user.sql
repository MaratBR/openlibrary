-- name: Moderation_GetUserInfo :one
select
    users.id,
    users.name,
    users.email,
    users.about,
    users.email_verified,
    users.joined_at,
    users.role,
    users.is_banned,
    (select count(*) from books where author_user_id = users.id) as books_total,
    (select count(*) from comments where user_id = users.id) as comments_total,
    (select count(*) from user_follower where followed_id = users.id) as followers_total
from users
where users.id = $1;

-- name: Moderation_SearchUsers :many
select
    users.id,
    users.name,
    users.joined_at,
    users.role,
    users.is_banned,
    (select sessions.created_at from sessions where sessions.user_id = users.id order by sessions.created_at desc limit 1) as last_visit_at,
    latest_ban.created_at as banned_at,
    coalesce(latest_ban.note, '') as ban_reason
from users
left join lateral (
    select user_bans.created_at, user_bans.note
    from user_bans
    where user_bans.user_id = users.id
    order by user_bans.created_at desc
    limit 1
) latest_ban on true
where
    (sqlc.arg('search')::text = '' or (sqlc.narg('search_id')::uuid is not null and users.id = sqlc.narg('search_id')) or (sqlc.narg('search_id')::uuid is null and users.name ilike '%' || sqlc.arg('search') || '%'))
    and (sqlc.arg('banned_status')::text = '' or users.is_banned = (sqlc.arg('banned_status') = 'banned'))
    and (sqlc.arg('role')::text = '' or users.role::text = sqlc.arg('role'))
order by users.name, users.id
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: Moderation_CountUsers :one
select count(*)
from users
where
    (sqlc.arg('search')::text = '' or (sqlc.narg('search_id')::uuid is not null and users.id = sqlc.narg('search_id')) or (sqlc.narg('search_id')::uuid is null and users.name ilike '%' || sqlc.arg('search') || '%'))
    and (sqlc.arg('banned_status')::text = '' or users.is_banned = (sqlc.arg('banned_status') = 'banned'))
    and (sqlc.arg('role')::text = '' or users.role::text = sqlc.arg('role'));

-- name: Moderation_GetUserBooks :many
select id, name, created_at, is_publicly_visible, is_banned, is_trashed
from books
where author_user_id = sqlc.arg('user_id')
order by created_at desc
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: Moderation_GetUserComments :many
select comments.id, comments.content, comments.created_at, comments.deleted_at,
       book_chapters.id as chapter_id, book_chapters.name as chapter_name,
       books.id as book_id, books.name as book_name
from comments
join book_chapters on book_chapters.id = comments.chapter_id
join books on books.id = book_chapters.book_id
where comments.user_id = sqlc.arg('user_id')
order by comments.created_at desc
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: Moderation_GetUserHistory :many
select moderation_logs.*, users.name as actor_user_name
from moderation_logs
left join users on users.id = moderation_logs.actor_user_id
where target_type = 'user' and target_id = sqlc.arg('user_id')::uuid::text
order by "time" desc
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: Moderation_CountUserHistory :one
select count(*) from moderation_logs
where target_type = 'user' and target_id = sqlc.arg('user_id')::uuid::text;

-- name: Moderation_GetUserReports :many
select reports.*, users.name as reporter_user_name
from reports
left join users on users.id = reports.reporter_user_id
where
    (target_type = 'user' and target_id = sqlc.arg('user_id')::uuid::text)
    or (target_type = 'book' and exists (
        select 1 from books where books.id::text = reports.target_id and books.author_user_id = sqlc.arg('user_id')
    ))
    or (target_type = 'comment' and exists (
        select 1 from comments where comments.id::text = reports.target_id and comments.user_id = sqlc.arg('user_id')
    ))
order by "time" desc
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: Moderation_CountUserReports :one
select count(*) from reports
where
    (target_type = 'user' and target_id = sqlc.arg('user_id')::uuid::text)
    or (target_type = 'book' and exists (
        select 1 from books where books.id::text = reports.target_id and books.author_user_id = sqlc.arg('user_id')
    ))
    or (target_type = 'comment' and exists (
        select 1 from comments where comments.id::text = reports.target_id and comments.user_id = sqlc.arg('user_id')
    ));

-- name: Moderation_GetRandomUserID :one
select id from users order by random() limit 1;
