-- Temporary moderation data: one synthetic report for every comment and book.
-- The stable report numbers make this script safe to run more than once.

insert into reports (
    number,
    "time",
    reporter_user_id,
    target_type,
    target_id,
    reason,
    description
)
select
    'TMP-COMMENT-' || comments.id,
    now(),
    null,
    'comment',
    comments.id::text,
    'Inappropriate comment',
    'Temporary report generated for moderation portal development.'
from comments
on conflict (number) do nothing;

insert into reports (
    number,
    "time",
    reporter_user_id,
    target_type,
    target_id,
    reason,
    description
)
select
    'TMP-BOOK-' || books.id,
    now(),
    null,
    'book',
    books.id::text,
    'Book requires review',
    'Temporary report generated for moderation portal development.'
from books
on conflict (number) do nothing;
