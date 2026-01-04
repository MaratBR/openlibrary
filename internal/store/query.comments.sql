-- name: GetChapterComments :many
select comments.*, users.name as user_name
from comments
join users on comments.user_id = users.id
where chapter_id = $1 and parent_id is null
order by created_at desc
limit $2;

-- name: GetChapterCommentsAfter :many
select comments.*, users.name as user_name
from comments
join users on comments.user_id = users.id
where chapter_id = $1 and parent_id is null and created_at < $3
order by created_at desc
limit $2;

-- name: GetChildComments :many
select comments.*, users.name as user_name
from comments
join users on comments.user_id = users.id
where parent_id = $1
order by created_at desc
limit $2;

-- name: GetChildCommentsAfter :many
select comments.*, users.name as user_name
from comments
join users on comments.user_id = users.id
where parent_id = $1 and created_at < $3
order by created_at desc
limit $2;


-- name: GetCommentByID :one
select * from comments where id = $1;

-- name: GetCommentWithUserByID :one
select c.*, u.name as user_name
from comments c
join users u on c.user_id = u.id
where c.id = $1;

-- name: InsertComment :exec 
insert into comments (id, chapter_id, parent_id, user_id, content, created_at, updated_at) 
values ($1, $2, $3, $4, $5, now(), now());

-- name: RecalculateCommentSubcomments :exec
update comments 
set subcomments = (select count(*) from comments c where c.parent_id = sqlc.arg('id'))
where id = sqlc.arg('id');



-- name: UpdateComment :execresult
update comments set content = $2, updated_at = now() where id = $1;