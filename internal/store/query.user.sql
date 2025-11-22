-- name: User_Get :one
select *
from users
where id = $1
limit 1;

-- name: User_GetDetails :one
select 
    users.*, 
    (select count(*) from books where author_user_id = users.id and is_publicly_visible and not is_banned and chapters > 0) as books_total,
    (select count(*) from user_follower where followed_id = users.id) as followers,
    (select count(*) from user_follower where follower_id = users.id) as "following",
    (user_follower.created_at is not null)::bool as is_following
from users
left join user_follower on user_follower.followed_id = users.id and user_follower.follower_id = sqlc.narg(actor_user_id)
where users.id = $1
limit 1;

-- name: User_FindByLogin :one
select *
from users
where name = $1 or email = $1
limit 1;

-- name: User_ExistsByUsername :one
select exists(select 1
from users
where name = $1);

-- name: User_ExistsByEmail :one
select exists(select 1
from users
where email = $1);

-- name: User_Insert :exec
insert into users
(id, name, password_hash, joined_at, email, email_verified)
values ($1, $2, $3, $4, $5, $6);

-- name: User_SetEmailVerified :exec
update users
set email_verified = $1
where id = $2;

-- name: Session_Insert :exec
insert into sessions
(id, sid, user_id, created_at, user_agent, ip_address, expires_at)
values ($1, $2, $3, $4, $5, $6, $7);

-- name: Session_Terminate :exec
update sessions
set is_terminated = true
where sid = $1;


-- name: Session_TerminateAllByUserID :exec
update sessions
set is_terminated = true
where user_id = $1;

-- name: Session_GetUserSessions :many
select s.*, u.id as user_id, u.name as user_name, u.joined_at as user_joined_at
from sessions s
join users u on s.user_id = u.id
where s.user_id = $1;

-- name: Session_GetInfo :one
select s.*, u.name as user_name, u.joined_at as user_joined_at, u."role" as user_role
from sessions s
join users u on s.user_id = u.id
where s.sid = $1;


-- name: GetUserPrivacySettings :one
select
    privacy_hide_stats,
    privacy_hide_comments,
    privacy_hide_email,
    privacy_allow_searching
from users
where id = $1;

-- name: User_UpdatePrivacySettings :exec
update users
set privacy_hide_stats = $2,
    privacy_hide_comments = $3,
    privacy_hide_email = $4,
    privacy_allow_searching = $5
where id = $1;

-- name: User_GetModerationSettings :one
select
    show_adult_content,
    censored_tags,
    censored_tags_mode
from users
where id = $1;

-- name: User_UpdateModerationSettings :exec
update users
set show_adult_content = $2,
    censored_tags = $3,
    censored_tags_mode = $4
where id = $1;

-- name: User_GetAboutSettings :one
select 
    about,
    gender
from users
where id = $1;

-- name: User_UpdateAboutSettings :exec
update users
set about = $2, gender = $3
where id = $1;


-- name: User_GetCustomizationSettings :one
select 
    profile_css,
    enable_profile_css,
    default_theme
from users
where id = $1;

-- name: User_UpdateCustomizationSettings :exec
update users
set profile_css = $2, enable_profile_css = $3, default_theme = $4
where id = $1;

-- name: User_Get2FADevices :many
select *
from user_2fa
where user_id = $1;

-- name: DeleteInactive2FADevices :exec
delete
from user_2fa
where not initialized and created_at < $1;

-- name: DeleteUserFollow :exec
delete
from user_follower
where follower_id = $1 and followed_id = $2;

-- name: User_InsertFollow :exec
insert into user_follower
(follower_id, followed_id)
values ($1, $2);

-- name: User_IsFollowing :one
select exists(select 1
from user_follower
where follower_id = $1 and followed_id = $2);

-- name: User_UpdatePassword :exec
update users
set password_hash = $2
where id = $1;

-- name: User_UpdateRole :exec
update users
set "role" = $2
where id = $1;

-- name: User_GetNames :many
select name, id
from users
where id = any(sqlc.arg(ids)::uuid[]);

-- name: EmailVerification_Delete :exec
delete from email_verification where email = $1;

-- name: EmailVerification_Insert :exec
insert into email_verification (email, user_id, valid_through, verification_code_hash)
values ($1, $2, $3, $4);

-- name: EmailVerification_Get :one
select * from email_verification where email = $1;