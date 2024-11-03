create type tag_type as enum (
    'freeform',
    'warning',
    'fandom',
    'reltype',
    'rel'
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