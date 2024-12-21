create type reading_list_status as enum ('dnf', 'reading', 'paused', 'read');

create table reading_list (
    user_id uuid not null references users(id),
    book_id int8 not null references books(id),
    status reading_list_status not null,
    primary key (user_id, book_id)
);