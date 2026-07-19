-- name: Draft_GetById :one
select drafts.*,
    books.id as book_id, books.name as book_name, 
    bc.is_publicly_visible as is_chapter_publicly_visible,
    bc.content_updated_at as chapter_content_updated_at
from drafts
join book_chapters bc on bc.id = drafts.chapter_id
join books on books.id = bc.book_id
where drafts.id = $1;

-- name: Draft_Insert :exec
insert into drafts (
    id, created_by, chapter_id, chapter_name, content, updated_at, created_at)
values ($1, $2, $3, $4, $5, $6, $7);

-- name: Draft_Update :exec
update drafts
set chapter_name = $2, words = $3, content = $4, summary = $5, updated_at = now()
where id = $1;

-- name: Draft_UpdateContent :exec
update drafts
set content = $2, words = $3, updated_at = now()
where id = $1;

-- name: Draft_UpdateChapterName :exec
update drafts
set chapter_name = $2, updated_at = now()
where id = $1;

-- name: Draft_Delete :exec
delete from drafts where id = $1;

-- name: Draft_MarkAsPublished :exec
update drafts
set published_at = now(), scheduled_at = null
where id = $1;

-- name: Draft_ClearChapterSchedules :exec
update drafts set scheduled_at = null where chapter_id = $1;

-- name: Draft_Schedule :exec
update drafts set scheduled_at = $2, updated_at = now() where id = $1;

-- name: Draft_PublishDue :many
with due as (
    select drafts.*
    from drafts
    where scheduled_at <= now()
    for update skip locked
), updated_chapters as (
    update book_chapters
    set name = due.chapter_name,
        content = due.content,
        words = due.words,
        summary = due.summary,
        is_publicly_visible = true,
        content_updated_at = due.updated_at,
        updated_at = now()
    from due
    where book_chapters.id = due.chapter_id
    returning book_chapters.book_id, due.id as draft_id
), published_drafts as (
    update drafts
    set published_at = now(), scheduled_at = null
    from updated_chapters
    where drafts.id = updated_chapters.draft_id
)
select distinct book_id from updated_chapters;

-- name: Draft_GetLatestID :one
select id
from drafts
where chapter_id = $1
order by coalesce(updated_at, created_at) desc
limit 1;

-- name: Draft_UserCanAccess :one
select exists (
    select 1
    from drafts
    join book_chapters on book_chapters.id = drafts.chapter_id
    join books on books.id = book_chapters.book_id
    where drafts.id = sqlc.arg('draft_id')
      and books.author_user_id = sqlc.arg('user_id')
      and (sqlc.arg('chapter_id')::int8 = 0 or book_chapters.id = sqlc.arg('chapter_id'))
      and (sqlc.arg('book_id')::int8 = 0 or books.id = sqlc.arg('book_id'))
);

-- name: Chapter_UserCanEdit :one
select exists (
    select 1
    from book_chapters
    join books on books.id = book_chapters.book_id
    where book_chapters.id = sqlc.arg('chapter_id')
	  and books.id = sqlc.arg('book_id')
      and books.author_user_id = sqlc.arg('user_id')
);
