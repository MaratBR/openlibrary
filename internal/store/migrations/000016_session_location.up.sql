alter table sessions
    add column location_country text not null default '',
    add column location_region text not null default '',
    add column location_city text not null default '';
