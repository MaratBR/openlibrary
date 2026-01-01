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