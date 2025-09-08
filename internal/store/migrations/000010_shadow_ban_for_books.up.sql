alter table books
    add column is_shadow_banned bool not null default false;
