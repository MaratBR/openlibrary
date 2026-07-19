-- name: ReaderPreferences_Get :one
select *
from user_reader_preferences
where user_id = $1;

-- name: ReaderPreferences_Upsert :exec
insert into user_reader_preferences (user_id, font_size, font_family, page_color, theme)
values ($1, $2, $3, $4, $5)
on conflict (user_id) do update set
    font_size = excluded.font_size,
    font_family = excluded.font_family,
    page_color = excluded.page_color,
    theme = excluded.theme,
    updated_at = now();
