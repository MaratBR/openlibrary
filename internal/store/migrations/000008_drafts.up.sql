create table drafts (
    id int8 primary key,
    created_by uuid not null references users(id),
    chapter_id int8 not null references book_chapters(id),
    chapter_name text not null,
    content text not null default '',
    words int4 not null default 0,
    summary text not null default '',
    is_adult_override boolean not null default false,
    updated_at timestamptz null,
    created_at timestamptz not null default now()
);