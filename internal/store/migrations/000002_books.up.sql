create table books (
    id int8 primary key,
    name varchar(255) not null,
    author_user_id uuid not null references users(id),
    created_at timestamptz not null default now()
);

create table book_chapters (
    id int8 primary key,
    name varchar(255),
    book_id int8 not null references books(id),
    content text not null,
    "order" int4 not null,
    created_at timestamptz not null default now()  
);
