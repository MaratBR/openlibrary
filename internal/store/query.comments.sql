-- name: Comment_GetByChapter :many
select comments.*, users.name as user_name
from comments
join users on comments.user_id = users.id
where chapter_id = $1 and parent_id is null
order by
    case when sqlc.arg('sort')::text = 'oldest' then comments.created_at end asc,
    case when sqlc.arg('sort')::text = 'popular' then (select count(*) from comments_liked where comment_id = comments.id) end desc,
    case when sqlc.arg('sort')::text in ('newest', 'popular') then comments.created_at end desc,
    comments.id desc
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: Comment_CountByChapter :one
select count(*) from comments where chapter_id = $1 and deleted_at is null;

-- name: Comment_GetByChapterAfter :many
select comments.*, users.name as user_name
from comments
join users on comments.user_id = users.id
where chapter_id = $1 and parent_id is null and created_at < $3
order by created_at desc
limit $2;

-- name: Comment_GetChildComments :many
select comments.*, users.name as user_name
from comments
join users on comments.user_id = users.id
where parent_id = sqlc.arg('parent_id')::int8
order by created_at asc
limit $1;

-- name: Comment_GetChildCommentsAfter :many
select comments.*, users.name as user_name
from comments
join users on comments.user_id = users.id
where parent_id = sqlc.arg('parent_id')::int8 and created_at > $2
order by created_at asc
limit $1;


-- name: Comment_GetByID :one
select * from comments where id = $1;

-- name: Comment_GetWithUserByID :one
select c.*, u.name as user_name
from comments c
join users u on c.user_id = u.id
where c.id = $1;

-- name: Comment_Insert :exec
insert into comments (id, chapter_id, parent_id, user_id, content, created_at)
values ($1, $2, $3, $4, $5, now());

-- name: Comment_RecalculateSubcomments :exec
update comments 
set subcomments = (select count(*) from comments c where c.parent_id = sqlc.arg('id') and c.deleted_at is null)
where id = sqlc.arg('id');


-- name: Comment_Update :execresult
update comments set content = $2, updated_at = now() where id = $1;

-- name: Comment_GetLikedComments :many
select comment_id, liked_at
from comments_liked
where user_id = $1 and comment_id = ANY(sqlc.arg('ids')::int8[]);

-- name: Comment_Like :exec
insert into comments_liked (comment_id, user_id)
values ($1, $2) on conflict (comment_id, user_id) do nothing;

-- name: Comment_UnLike :exec
delete from comments_liked where comment_id = $1 and user_id = $2;
