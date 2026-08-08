-- Temporary development data: whole-book, chapter, and selected-text reports.
-- Stable report numbers make this script safe to run repeatedly.

insert into reports (number, "time", reporter_user_id, target_type, target_id, reason, description)
select 'TMP-BOOK-' || books.id, now(), books.author_user_id, 'book', books.id::text,
       'Book requires review', 'Temporary whole-book report.'
from books
on conflict (number) do nothing;

insert into reports (number, "time", reporter_user_id, target_type, target_id, reason, description, book_chapter_id)
select 'TMP-CHAPTER-' || books.id, now(), books.author_user_id, 'book', books.id::text,
       'Chapter requires review', 'Temporary chapter report.', chapter.id
from books
join lateral (
    select book_chapters.id
    from book_chapters
    where book_chapters.book_id = books.id
    order by book_chapters."order", book_chapters.id
    limit 1
) chapter on true
on conflict (number) do nothing;

insert into reports (number, "time", reporter_user_id, target_type, target_id, reason, description, book_chapter_id, book_excerpt)
select 'TMP-TEXT-' || books.id, now(), books.author_user_id, 'book', books.id::text,
       'Selected text requires review', 'Temporary selected-text report.', chapter.id,
       coalesce(nullif(left(btrim(regexp_replace(chapter.content, '<[^>]*>', '', 'g')), 500), ''), '[selected text unavailable]')
from books
join lateral (
    select book_chapters.id, book_chapters.content
    from book_chapters
    where book_chapters.book_id = books.id and book_chapters.content <> ''
    order by book_chapters."order", book_chapters.id
    limit 1
) chapter on true
on conflict (number) do nothing;
