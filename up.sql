alter table reports
    add column if not exists book_chapter_id int8 null references book_chapters (id),
    add column if not exists book_excerpt text not null default '',
    add column if not exists status text not null default 'unreviewed',
    add column if not exists priority text not null default 'medium';

alter table reports
    drop column if exists assigned_to,
    drop column if exists assigned_team,
    drop column if exists channel,
    drop column if exists sla_deadline,
    drop column if exists tags;

update reports
set reporter_user_id = coalesce(
    (
        select books.author_user_id
        from books
        where reports.target_type = 'book' and books.id::text = reports.target_id
    ),
    (select users.id from users order by users.id limit 1)
)
where reporter_user_id is null;

alter table reports alter column reporter_user_id set not null;
