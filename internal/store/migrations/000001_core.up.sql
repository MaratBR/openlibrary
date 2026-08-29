CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;

create type user_role as enum ('user', 'admin', 'moderator', 'system');
create type censor_mode as enum ('hide', 'censor', 'none');
create type type_of_2fa as enum ('totp', 'webauthn');
create type age_rating as enum ('?', 'G', 'PG', 'PG-13', 'R', 'NC-17');
create type tag_type as enum ('freeform', 'warning', 'fandom', 'reltype', 'rel', 'genre');
create type reading_list_status as enum ('dnf', 'reading', 'paused', 'read', 'want_to_read');

create table users (
    id uuid primary key,
    name varchar(255) not null,
    email varchar(255) not null default '',
    email_verified boolean not null default false,
    joined_at timestamptz not null default now(),
    password_hash text not null,
    "role" user_role not null default 'user',
    is_banned boolean not null default false,
    avatar_file text null,
    about text not null default '',
    gender text not null default '',
    profile_css text not null default '',
    enable_profile_css boolean not null default false,
    default_theme text not null default '',
    privacy_hide_stats boolean not null default false,
    privacy_hide_comments boolean not null default false,
    privacy_hide_email boolean not null default true,
    privacy_allow_searching boolean not null default false,
    show_adult_content boolean not null default false,
    censored_tags text[] not null default '{}',
    censored_tags_mode censor_mode not null default 'none'
);

create unique index ix_users_name on users (name);
create unique index ix_users_email on users (email) where email != '';

create table user_2fa (
    id uuid primary key,
    user_id uuid not null references users(id),
    "type" type_of_2fa not null,
    "key" text not null,
    created_at timestamptz not null default now(),
    initialized boolean not null default false,
    active boolean not null default true
);

create table user_bans (
    id bigint primary key,
    user_id uuid not null references users(id),
    created_at timestamptz not null default now(),
    banned_by_user_id uuid null references users(id),
    note text not null default '',
    expires_at timestamptz not null
);

create table email_verification (
    email text primary key,
    user_id uuid not null references users(id),
    verification_code_hash text not null,
    created_at timestamptz not null default now(),
    valid_through timestamptz not null
);

create table user_follower (
    follower_id uuid not null references users(id),
    followed_id uuid not null references users(id),
    created_at timestamptz not null default now(),
    primary key (follower_id, followed_id)
);

create table sessions (
    id int8 primary key,
    sid text not null,
    user_id uuid not null references users(id),
    created_at timestamptz not null default now(),
    user_agent text not null,
    ip_address text not null,
    expires_at timestamptz not null,
    is_terminated boolean not null default false,
    location_country text not null default '',
    location_region text not null default '',
    location_city text not null default ''
);

create table books (
    id int8 primary key,
    name text not null,
    slug varchar(80) not null,
    summary text not null,
    author_user_id uuid not null references users(id),
    created_at timestamptz not null default now(),
    age_rating age_rating not null default '?',
    is_publicly_visible boolean not null default true,
    is_banned boolean not null default false,
    is_trashed boolean not null default false,
    words int4 not null default 0,
    chapters int4 not null default 0,
    tag_ids int8[] not null default '{}',
    cached_parent_tag_ids int8[] not null default '{}',
    cover text not null default '',
    view int4 not null default 0,
    rating float8 null,
    total_reviews int4 not null default 0,
    total_ratings int4 not null default 0,
    is_pinned boolean not null default false
);

create index ix_books_author_user_id on books (author_user_id);
create index ix_books_age_rating on books (age_rating);
create index ix_books_tags on books using gin (cached_parent_tag_ids);
create index ix_books_name on books using gin (name);

create table book_chapters (
    id int8 primary key,
    name varchar(255) not null,
    book_id int8 not null references books(id),
    content text not null,
    content_updated_at timestamptz not null default now(),
    "order" int4 not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    words int4 not null default 0,
    summary text not null default '',
    fonts text[] not null default '{}',
    is_publicly_visible bool not null default false
);

create index ix_bok_chapters_book_id on book_chapters (book_id);
create index ix_bok_chapters_order on book_chapters ("order");

create table collections (
    id int8 primary key,
    name varchar(255) not null,
    slug varchar(80) not null,
    summary text not null default '',
    user_id uuid not null references users(id),
    created_at timestamptz not null default now(),
    books_count int4 not null default 0,
    last_updated_at timestamptz null,
    is_public boolean not null default true
);

create index ix_collections_user_id on collections (user_id);

create table collection_books (
    collection_id int8 not null references collections(id),
    book_id int8 not null references books(id),
    added_at timestamptz not null default now(),
    "order" int4 not null default 0,
    primary key (collection_id, book_id)
);

create table defined_tags (
    id int8 primary key,
    name text not null,
    description text not null default '',
    is_spoiler boolean not null default false,
    is_adult boolean not null default false,
    created_at timestamptz not null default now(),
    tag_type tag_type not null,
    synonym_of int8 null references defined_tags(id),
    is_default boolean not null default false,
    lowercased_name text not null generated always as (lower(name)) stored
);

create unique index ix_defined_tags_name on defined_tags (name text_pattern_ops);

create table reading_list (
    user_id uuid not null references users(id),
    book_id int8 not null references books(id),
    status reading_list_status not null,
    last_accessed_chapter_id int8 null references book_chapters(id),
    last_updated_at timestamptz not null default now(),
    primary key (user_id, book_id)
);

create index ix_reading_list_chapter_id on reading_list (last_accessed_chapter_id);

create table reading_list_history (
    user_id uuid not null references users(id),
    book_id int8 not null references books(id),
    chapter_id int8 not null references book_chapters(id),
    finished_reading boolean not null default false,
    progress int4 not null default 0,
    primary key (user_id, chapter_id)
);

create index ix_reading_list_history_chapter_id on reading_list_history (chapter_id);

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
create index ix_comments_special_sub_comments on comments (parent_id, created_at);

create table comments_liked (
    comment_id int8 not null references comments(id),
    user_id uuid not null references users(id),
    liked_at timestamptz not null default now(),
    primary key (comment_id, user_id)
);

create table ratings (
    user_id uuid not null references users(id),
    book_id int8 not null references books(id),
    rating int2 not null,
    updated_at timestamptz not null default now(),
    primary key (user_id, book_id)
);

create table reviews (
    user_id uuid not null references users(id),
    book_id int8 not null references books(id),
    content text not null default '',
    created_at timestamptz not null default now(),
    last_updated_at timestamptz null,
    likes int4 not null default 0,
    primary key (user_id, book_id)
);

create table drafts (
    id int8 primary key,
    created_by uuid not null references users(id),
    chapter_id int8 not null references book_chapters(id),
    chapter_name text not null,
    content text not null default '',
    words int4 not null default 0,
    summary text not null default '',
    updated_at timestamptz null,
    created_at timestamptz not null default now(),
    published_at timestamptz null,
    scheduled_at timestamptz null
);

create index drafts_scheduled_at_idx on drafts (scheduled_at) where scheduled_at is not null;

create table draft_log (
    id int8 primary key,
    draft_id int8 not null references drafts(id),
    created_at timestamptz not null default now(),
    user_id uuid null references users(id),
    payload jsonb not null
);

create table site_config (
    value jsonb,
    key text primary key
);

create table user_reader_preferences (
    user_id uuid primary key references users(id) on delete cascade,
    font_size smallint not null default 18 check (font_size in (12, 14, 16, 18, 20, 22, 26, 30, 36, 42, 48)),
    font_family text not null default 'serif' check (font_family in ('serif', 'sans', 'dyslexic')),
    page_color text not null default 'background' check (page_color in ('background', 'surface')),
    theme text not null default 'system' check (theme in ('system', 'light', 'dark')),
    updated_at timestamptz not null default now()
);

create table user_data (
    key text not null,
    user_id uuid not null references users(id) on delete cascade,
    data jsonb not null,

    primary key (key, user_id)
);

create index ix_books_author_created_at on books (author_user_id, created_at desc);
create index ix_books_visible_author_pinned_created_at on books (author_user_id, is_pinned desc, created_at)
    where is_publicly_visible and not is_banned and not is_trashed and chapters > 0;
create index ix_books_visible_author_rating on books (author_user_id, rating desc) where is_publicly_visible;
create index ix_books_tag_ids on books using gin (tag_ids);
create index ix_book_chapters_book_order on book_chapters (book_id, "order");
create index ix_collections_user_created_at on collections (user_id, created_at desc);
create index ix_collections_user_last_updated_at on collections (user_id, last_updated_at desc);
create index ix_collection_books_book_id on collection_books (book_id);
create index ix_collection_books_collection_order on collection_books (collection_id, "order");
create index ix_reading_list_user_status_updated_at on reading_list (user_id, status, last_updated_at);
create index ix_drafts_chapter_latest on drafts (chapter_id, (coalesce(updated_at, created_at)) desc);
create index ix_comments_live_chapter on comments (chapter_id) where deleted_at is null;
create index ix_comments_user_created_at on comments (user_id, created_at desc);
create index ix_comments_liked_user_comment on comments_liked (user_id, comment_id);
create index ix_ratings_book_id on ratings (book_id);
create index ix_reviews_book_created_at on reviews (book_id, created_at desc);
create index ix_user_follower_followed_id on user_follower (followed_id);
create index ix_user_2fa_user_id on user_2fa (user_id);
create index ix_user_2fa_uninitialized_created_at on user_2fa (created_at) where not initialized;
create index ix_sessions_sid on sessions (sid);
create index ix_sessions_user_created_at on sessions (user_id, created_at desc);
create index ix_user_bans_user_created_at on user_bans (user_id, created_at desc);
create index ix_defined_tags_lowercased_name on defined_tags (lowercased_name text_pattern_ops);
create index ix_defined_tags_type_lowercased_name on defined_tags (tag_type, lowercased_name text_pattern_ops);
create index ix_defined_tags_synonym_of on defined_tags (synonym_of);
