-- name: GetDraftById :one
select *
from drafts 
where id = $1;

-- name: InsertDraft :exec
insert into drafts (
    id, created_by, chapter_id, chapter_name, content, updated_at, created_at)
values ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateDraft :exec
update drafts
set chapter_name = $2, is_adult_override = $3, words = $4, content = $5, summary = $6, updated_at = now()
where id = $1;

-- name: DeleteDraft :exec
delete from drafts where id = $1;

-- name: MarkDraftAsPublished :exec
update drafts
set published_at = now()
where id = $1;