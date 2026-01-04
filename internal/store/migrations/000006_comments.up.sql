create table comments (
    id int8 primary key,
    chapter_id int8 not null references book_chapters(id),
    user_id uuid not null references users(id),
    content text not null,
    created_at timestamptz not null,
    updated_at timestamptz,
    deleted_at timestamptz,
    parent_id int8 null references comments(id),
    subcomments int4 not null default 0
);