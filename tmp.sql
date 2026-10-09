-- Update an existing development database after the consolidated core schema
-- gained chapter font metadata. This is intentionally idempotent.
alter table book_chapters
    add column if not exists fonts text[] not null default '{}';

-- Reader content width for existing databases.
alter table user_reader_preferences
    add column if not exists content_width smallint not null default 72
    check (content_width between 48 and 100 and content_width % 2 = 0);
