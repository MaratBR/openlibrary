-- name: Moderation_SearchChapters :many
select c.id, c.name as title, c.created_at, c.updated_at, (c.is_publicly_visible)::boolean as visible, b.id as book_id, b.name as book_name, u.id as author_id, u.name as author_name, c.words
from book_chapters c join books b on b.id = c.book_id join users u on u.id = b.author_user_id
where (sqlc.arg('search')::text = '' or c.id::text = sqlc.arg('search') or c.name ilike '%' || sqlc.arg('search') || '%' or b.name ilike '%' || sqlc.arg('search') || '%' or u.name ilike '%' || sqlc.arg('search') || '%')
 and (sqlc.arg('book_id')::text = '' or b.id::text = sqlc.arg('book_id'))
 and (sqlc.arg('author_id')::text = '' or u.id::text = sqlc.arg('author_id'))
 and (sqlc.arg('state')::text = '' or (sqlc.arg('state') = 'visible' and (c.is_publicly_visible)) or (sqlc.arg('state') = 'hidden' and not (c.is_publicly_visible)))
order by c.created_at desc, c.id desc
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: Moderation_CountChapters :one
select count(*)
from book_chapters c join books b on b.id = c.book_id join users u on u.id = b.author_user_id
where (sqlc.arg('search')::text = '' or c.id::text = sqlc.arg('search') or c.name ilike '%' || sqlc.arg('search') || '%' or b.name ilike '%' || sqlc.arg('search') || '%' or u.name ilike '%' || sqlc.arg('search') || '%')
 and (sqlc.arg('book_id')::text = '' or b.id::text = sqlc.arg('book_id'))
 and (sqlc.arg('author_id')::text = '' or u.id::text = sqlc.arg('author_id'))
 and (sqlc.arg('state')::text = '' or (sqlc.arg('state') = 'visible' and (c.is_publicly_visible)) or (sqlc.arg('state') = 'hidden' and not (c.is_publicly_visible)));

-- name: Moderation_ChapterDetail :one
select c.id, c.name as title, c.created_at, c.updated_at, (c.is_publicly_visible)::boolean as visible, b.id as book_id, b.name as book_name, u.id as author_id, u.name as author_name, c.words, c.content
from book_chapters c join books b on b.id = c.book_id join users u on u.id = b.author_user_id
where c.id = $1;

-- name: Moderation_SearchComments :many
select c.id, c.content as title, c.created_at, c.updated_at, (c.deleted_at is null)::boolean as visible, b.id as book_id, b.name as book_name, u.id as author_id, u.name as author_name, c.chapter_id, ch.name as chapter_name, c.parent_id
from comments c join book_chapters ch on ch.id = c.chapter_id join books b on b.id = ch.book_id join users u on u.id = c.user_id
where (sqlc.arg('search')::text = '' or c.id::text = sqlc.arg('search') or c.content ilike '%' || sqlc.arg('search') || '%' or b.name ilike '%' || sqlc.arg('search') || '%' or u.name ilike '%' || sqlc.arg('search') || '%')
 and (sqlc.arg('book_id')::text = '' or b.id::text = sqlc.arg('book_id'))
 and (sqlc.arg('author_id')::text = '' or u.id::text = sqlc.arg('author_id'))
 and (sqlc.arg('state')::text = '' or (sqlc.arg('state') = 'visible' and (c.deleted_at is null)) or (sqlc.arg('state') = 'hidden' and not (c.deleted_at is null)))
 and (sqlc.arg('chapter_id')::text = '' or c.chapter_id::text = sqlc.arg('chapter_id'))
order by c.created_at desc, c.id desc
limit sqlc.arg('page_limit') offset sqlc.arg('page_offset');

-- name: Moderation_CountComments :one
select count(*)
from comments c join book_chapters ch on ch.id = c.chapter_id join books b on b.id = ch.book_id join users u on u.id = c.user_id
where (sqlc.arg('search')::text = '' or c.id::text = sqlc.arg('search') or c.content ilike '%' || sqlc.arg('search') || '%' or b.name ilike '%' || sqlc.arg('search') || '%' or u.name ilike '%' || sqlc.arg('search') || '%')
 and (sqlc.arg('book_id')::text = '' or b.id::text = sqlc.arg('book_id'))
 and (sqlc.arg('author_id')::text = '' or u.id::text = sqlc.arg('author_id'))
 and (sqlc.arg('state')::text = '' or (sqlc.arg('state') = 'visible' and (c.deleted_at is null)) or (sqlc.arg('state') = 'hidden' and not (c.deleted_at is null)))
 and (sqlc.arg('chapter_id')::text = '' or c.chapter_id::text = sqlc.arg('chapter_id'));

-- name: Moderation_CommentDetail :one
select c.id, c.content as title, c.created_at, c.updated_at, (c.deleted_at is null)::boolean as visible, b.id as book_id, b.name as book_name, u.id as author_id, u.name as author_name, c.chapter_id, ch.name as chapter_name, c.parent_id, c.content
from comments c join book_chapters ch on ch.id = c.chapter_id join books b on b.id = ch.book_id join users u on u.id = c.user_id
where c.id = $1;

-- name: Moderation_CommentReplies :many
select id from comments where parent_id = $1 order by created_at, id limit 20;
