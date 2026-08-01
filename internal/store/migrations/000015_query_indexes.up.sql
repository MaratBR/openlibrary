-- Books and chapters.
create index ix_books_author_created_at on books (author_user_id, created_at desc);
create index ix_books_visible_author_pinned_created_at
    on books (author_user_id, is_pinned desc, created_at)
    where is_publicly_visible and not is_banned and not is_trashed and chapters > 0;
create index ix_books_visible_author_rating
    on books (author_user_id, rating desc)
    where is_publicly_visible;
create index ix_books_tag_ids on books using gin (tag_ids);
create index ix_book_chapters_book_order on book_chapters (book_id, "order");

-- Collections.
create index ix_collections_user_created_at on collections (user_id, created_at desc);
create index ix_collections_user_last_updated_at on collections (user_id, last_updated_at desc);
create index ix_collection_books_book_id on collection_books (book_id);
create index ix_collection_books_collection_order on collection_books (collection_id, "order");

-- Reading list filtering and ordering.
create index ix_reading_list_user_status_updated_at
    on reading_list (user_id, status, last_updated_at);

-- Draft lookup and publication.
create index ix_drafts_chapter_latest
    on drafts (chapter_id, (coalesce(updated_at, created_at)) desc);

-- Comments and reviews.
create index ix_comments_live_chapter on comments (chapter_id) where deleted_at is null;
create index ix_comments_user_created_at on comments (user_id, created_at desc);
create index ix_comments_liked_user_comment on comments_liked (user_id, comment_id);
create index ix_ratings_book_id on ratings (book_id);
create index ix_reviews_book_created_at on reviews (book_id, created_at desc);

-- Users, authentication, and moderation.
create index ix_user_follower_followed_id on user_follower (followed_id);
create index ix_user_2fa_user_id on user_2fa (user_id);
create index ix_user_2fa_uninitialized_created_at
    on user_2fa (created_at) where not initialized;
create index ix_sessions_sid on sessions (sid);
create index ix_sessions_user_created_at on sessions (user_id, created_at desc);
create index ix_user_bans_user_created_at on user_bans (user_id, created_at desc);

-- Tag search and cleanup.
create index ix_defined_tags_lowercased_name
    on defined_tags (lowercased_name text_pattern_ops);
create index ix_defined_tags_type_lowercased_name
    on defined_tags (tag_type, lowercased_name text_pattern_ops);
create index ix_defined_tags_synonym_of on defined_tags (synonym_of);
