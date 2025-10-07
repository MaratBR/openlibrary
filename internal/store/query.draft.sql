-- name: GetDraftById :one
select drafts.*, books.id as book_id, books.name as book_name, bc.is_publicly_visible as is_chapter_publicly_visible
from drafts
join book_chapters bc on bc.id = drafts.chapter_id
join books on books.id = bc.book_id
where drafts.id = $1;

-- name: InsertDraft :exec
insert into drafts (
    id, created_by, chapter_id, chapter_name, content, updated_at, created_at)
values ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateDraft :exec
update drafts
set chapter_name = $2, is_adult_override = $3, words = $4, content = $5, summary = $6, updated_at = now()
where id = $1;

-- name: UpdateDraftContent :exec
update drafts
set content = $2, words = $3, updated_at = now()
where id = $1;

-- name: UpdateDraftChapterName :exec
update drafts
set chapter_name = $2, updated_at = now()
where id = $1;

-- name: DeleteDraft :exec
delete from drafts where id = $1;

-- name: MarkDraftAsPublished :exec
update drafts
set published_at = now()
where id = $1;

-- name: GetLatestDraftID :one
select id
from drafts
where chapter_id = $1
order by coalesce(updated_at, created_at) desc
limit 1;