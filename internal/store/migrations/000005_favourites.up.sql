create table favorites (
    user_id uuid not null,
    book_id int8 not null,
    is_favorite bool not null default true,
    created_at timestamptz not null default now(),
    primary key (user_id, book_id)
)