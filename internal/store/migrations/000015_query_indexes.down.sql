drop index if exists ix_defined_tags_synonym_of;
drop index if exists ix_defined_tags_type_lowercased_name;
drop index if exists ix_defined_tags_lowercased_name;

drop index if exists ix_book_logs_book_action_time;
drop index if exists ix_book_logs_book_time;
drop index if exists ix_user_bans_user_created_at;
drop index if exists ix_sessions_user_created_at;
drop index if exists ix_sessions_sid;
drop index if exists ix_user_2fa_uninitialized_created_at;
drop index if exists ix_user_2fa_user_id;
drop index if exists ix_user_follower_followed_id;

drop index if exists ix_reviews_book_created_at;
drop index if exists ix_ratings_book_id;
drop index if exists ix_comments_liked_user_comment;
drop index if exists ix_comments_live_chapter;

drop index if exists ix_drafts_chapter_latest;

drop index if exists ix_reading_list_user_status_updated_at;

drop index if exists ix_collection_books_collection_order;
drop index if exists ix_collection_books_book_id;
drop index if exists ix_collections_user_last_updated_at;
drop index if exists ix_collections_user_created_at;

drop index if exists ix_book_chapters_book_order;
drop index if exists ix_books_tag_ids;
drop index if exists ix_books_visible_author_rating;
drop index if exists ix_books_visible_author_pinned_created_at;
drop index if exists ix_books_author_created_at;
