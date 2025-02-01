create type reading_list_status as enum ('dnf', 'reading', 'paused', 'read', 'want_to_read');

create table reading_list (
    user_id uuid not null references users(id),
    book_id int8 not null references books(id),
    status reading_list_status not null,
    last_accessed_chapter_id int8 null references book_chapters(id),
    last_updated_at timestamptz not null default now(),
    primary key (user_id, book_id)
);

-- create index ix_reading_list_user_id on reading_list (user_id);
-- create index ix_reading_list_book_id on reading_list (book_id);
create index ix_reading_list_chapter_id on reading_list (last_accessed_chapter_id);

create table reading_list_history (
    user_id uuid not null references users(id),
    book_id int8 not null references books(id),
    chapter_id int8 not null references book_chapters(id),
    finished_reading boolean not null default false,
    progress int4 not null default 0,
    primary key (user_id, chapter_id)
);

-- create index ix_reading_list_history_user_id on reading_list_history (user_id);
-- create index ix_reading_list_history_book_id on reading_list_history (book_id);
create index ix_reading_list_history_chapter_id on reading_list_history (chapter_id);
