create table collections (
    id int8 primary key,
    name varchar(255) not null,
    user_id uuid not null references users(id),
    created_at timestamptz not null default now(),
    books_count int4 not null default 0
);

create index ix_collections_user_id on collections (user_id);

create table collection_books (
    collection_id int8 not null references collections(id),
    book_id int8 not null references books(id),
    "order" int4 not null default 0,
    primary key (collection_id, book_id)
);