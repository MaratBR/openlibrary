create schema ol_analytics;

create table ol_analytics.view_bucket (
    "period" int4 not null,
    book_id int8 not null,
    count int8 not null default 0,

    primary key ("period", book_id)
);


create index ix_view_bucket_book_id on ol_analytics.view_bucket(book_id);
