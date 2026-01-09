create table comments (
    id int8 primary key,
    chapter_id int8 not null references book_chapters(id),
    user_id uuid not null references users(id),
    content text not null,
    created_at timestamptz not null,
    updated_at timestamptz,
    deleted_at timestamptz,
    parent_id int8 null references comments(id),
    subcomments int4 not null default 0,
    likes int4 not null default 0,
    likes_recalculated_at timestamptz not null default now()
);

create index ix_comments_special_root_comments on comments (chapter_id, created_at) where parent_id is null;
create index ix_comments_special_sub_comments on comments (chapter_id, parent_id, created_at);

create table comments_liked (
    comment_id int8 not null references comments(id),
    user_id uuid not null references users(id),
    liked_at timestamptz not null default now(),

    primary key (comment_id, user_id)
);