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
    words int4 not null default 0,
    chapters int4 not null default 0,
    tag_ids int8[] not null default '{}',
    cached_parent_tag_ids int8[] not null default '{}'
);

create index ix_books_author_user_id on books (author_user_id);
create index ix_books_age_rating on books (age_rating);
create index ix_books_tags on books using gin(cached_parent_tag_ids);
create index ix_books_name on books using gin(name);


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

