-- name: CLI_Util_MinifyContent_ListChapters :many
select id, content, words
from book_chapters
where id > $1
order by id
limit $2;

-- name: CLI_Util_MinifyContent_UpdateChapter :exec
update book_chapters
set content = $2, words = $3
where id = $1;

-- name: CLI_Util_MinifyContent_ListBookSummaries :many
select id, summary
from books
where id > $1
order by id
limit $2;

-- name: CLI_Util_MinifyContent_UpdateBookSummary :exec
update books
set summary = $2
where id = $1;
