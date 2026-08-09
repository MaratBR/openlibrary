-- Development data for comments and reviews.
--
-- Run with:
--   psql "$DATABASE_URL" -f populate.sql
--
-- The script uses existing users, books, and chapters. A book needs at least
-- reviews_per_book + additional_reviews_per_book distinct users to receive the
-- requested maximum number of reviews.

begin;

do $populate$
declare
    chapter_limit integer := 100;
    comments_per_chapter integer := 20;
    reviews_per_book integer := 10;
    additional_reviews_per_book integer := 15;
    books_with_additional_reviews integer := 10;
begin
    if chapter_limit < 0
        or comments_per_chapter < 0
        or reviews_per_book < 0
        or additional_reviews_per_book < 0
        or books_with_additional_reviews < 0
    then
        raise exception 'Population limits must not be negative';
    end if;

    -- Add comments to at most chapter_limit chapters. IDs are allocated above
    -- the current maximum because comments.id does not have a sequence.
    with selected_chapters as (
        select book_chapters.id, book_chapters.name
        from book_chapters
        order by book_chapters.id
        limit chapter_limit
    ),
    selected_users as (
        select users.id,
               row_number() over (order by users.id) as user_number
        from users
    ),
    generated_comments as (
        select selected_chapters.id as chapter_id,
               selected_chapters.name as chapter_name,
               comment_number,
               selected_users.id as user_id,
               row_number() over (order by selected_chapters.id, comment_number) as id_offset
        from selected_chapters
        cross join generate_series(1, comments_per_chapter) as comment_number
        join selected_users
          on selected_users.user_number =
             ((comment_number - 1) % greatest((select count(*) from selected_users), 1)) + 1
    ),
    current_maximum as (
        select coalesce(max(comments.id), 0) as maximum_comment_id
        from comments
    )
    insert into comments (id, chapter_id, user_id, content, created_at)
    select current_maximum.maximum_comment_id + generated_comments.id_offset,
           generated_comments.chapter_id,
           generated_comments.user_id,
           format(
               'Generated comment %s for chapter "%s".',
               generated_comments.comment_number,
               generated_comments.chapter_name
           ),
           now() - generated_comments.id_offset * interval '1 minute'
    from generated_comments
    cross join current_maximum;

    -- Add reviews_per_book reviews to every book. Existing reviewers are
    -- excluded, so running the script again adds another batch when enough
    -- unused users remain.
    insert into reviews (user_id, book_id, content, created_at)
    select reviewer.id,
           books.id,
           format('Generated review %s for "%s".', reviewer.review_number, books.name),
           now() - reviewer.review_number * interval '1 hour'
    from books
    cross join lateral (
        select users.id,
               row_number() over (order by users.id) as review_number
        from users
        where not exists (
            select 1
            from reviews
            where reviews.book_id = books.id
              and reviews.user_id = users.id
        )
        order by users.id
        limit reviews_per_book
    ) as reviewer;

    -- Give a limited number of books an additional batch of reviews.
    with selected_books as (
        select books.id, books.name
        from books
        order by books.id
        limit books_with_additional_reviews
    )
    insert into reviews (user_id, book_id, content, created_at)
    select reviewer.id,
           selected_books.id,
           format(
               'Generated additional review %s for "%s".',
               reviewer.review_number,
               selected_books.name
           ),
           now() - (reviews_per_book + reviewer.review_number) * interval '1 hour'
    from selected_books
    cross join lateral (
        select users.id,
               row_number() over (order by users.id) as review_number
        from users
        where not exists (
            select 1
            from reviews
            where reviews.book_id = selected_books.id
              and reviews.user_id = users.id
        )
        order by users.id
        limit additional_reviews_per_book
    ) as reviewer;

    -- Review pages join reviews to ratings, so every generated review needs a
    -- rating. Preserve ratings that already exist.
    insert into ratings (user_id, book_id, rating)
    select reviews.user_id,
           reviews.book_id,
           (1 + mod(
               hashtextextended(reviews.user_id::text || ':' || reviews.book_id, 0)
                   & 9223372036854775807,
               5
           ))::int2
    from reviews
    on conflict (user_id, book_id) do nothing;

    update books
    set rating = (select avg(ratings.rating::float8) from ratings where ratings.book_id = books.id),
        total_ratings = (select count(*) from ratings where ratings.book_id = books.id),
        total_reviews = (select count(*) from reviews where reviews.book_id = books.id);
end
$populate$;

commit;
