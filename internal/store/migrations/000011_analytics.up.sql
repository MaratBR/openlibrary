create schema ol_analytics;

create table ol_analytics.view_counter (
    period        int4    not null,
    entity_type   smallint not null,
    entity_id     int8    not null,
    book_id       int8    not null,   -- always populated
    view_count    int8    not null default 0,

    primary key (period, entity_type, entity_id)
);

create index ix_view_counter_book_period
    on ol_analytics.view_counter (period, book_id);
