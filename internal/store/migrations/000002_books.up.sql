CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;

create type age_rating as enum (
    '?',
    'G',
    'PG',
    'PG-13',
    'R',
    'NC-17'
);

create table books (
    id int8 primary key,
    name text not null,
    summary text not null,
    author_user_id uuid not null references users(id),
    created_at timestamptz not null default now(),
    age_rating age_rating not null default '?',
    is_publicly_visible boolean not null default true,
    is_banned boolean not null default false,
    words int4 not null default 0,
    chapters int4 not null default 0,
    tag_ids int8[] not null default '{}',
    cached_parent_tag_ids int8[] not null default '{}',
    favorites int4 not null default 0,
    has_cover boolean not null default false,
    view int4 not null default 0,
    rating float8 null,
    total_reviews int4 not null default 0,
    total_ratings int4 not null default 0
);

create index ix_books_author_user_id on books (author_user_id);
create index ix_books_age_rating on books (age_rating);
create index ix_books_tags on books using gin(cached_parent_tag_ids);
create index ix_books_name on books using gin(name);


create table book_ban_history (
    book_id int8 not null references books(id) primary key,
    user_id uuid not null references users(id),
    created_at timestamptz not null default now(),
    reason text not null default '',
    action text not null
);

create table book_chapters (
    id int8 primary key,
    name varchar(255) not null,
    book_id int8 not null references books(id),
    content text not null,
    "order" int4 not null,
    created_at timestamptz not null default now(),
    words int4 not null default 0,
    is_adult_override bool not null default false,
    summary text not null default ''
);

create index ix_bok_chapters_book_id on book_chapters (book_id);
create index ix_bok_chapters_order on book_chapters ("order");

create table book_view (
    ip_address inet not null,
    book_id int8 not null references books(id),
    recorded_at timestamptz not null 
);

