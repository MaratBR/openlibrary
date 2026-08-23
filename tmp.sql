-- Update an existing development database after the consolidated core schema
-- gained chapter font metadata. This is intentionally idempotent.
alter table book_chapters
    add column if not exists fonts text[] not null default '{}';
