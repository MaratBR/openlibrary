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
