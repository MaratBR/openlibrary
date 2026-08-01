-- name: Report_Create :exec
insert into reports (id, "time", reporter_user_id, target_type, target_id, reason, description)
values ($1, $2, $3, $4, $5, $6, $7);

-- name: Report_UserExists :one
select exists(select 1 from users where id = $1);

-- name: Report_BookExists :one
select exists(select 1 from books where id = $1);

-- name: Report_CommentExists :one
select exists(select 1 from comments where id = $1);
