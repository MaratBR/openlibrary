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
    created_at timestamptz not null default now(),
    published_at timestamptz null
);

create table draft_log (
    id int8 primary key,
    draft_id int8 not null references drafts(id),
    created_at timestamptz not null default now(),
    user_id uuid null references users(id),
    payload jsonb not null
);