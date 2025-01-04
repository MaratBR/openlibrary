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