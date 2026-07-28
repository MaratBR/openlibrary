create schema ol_analytics;

create table ol_analytics.interaction_event (
    id bigint generated always as identity primary key,
    user_key text not null,
    book_id int8 not null,
    event_type text not null,
    value double precision not null default 1,
    created_at timestamptz not null default now()
);

create type ol_analytics.bucket_period_type as enum (
    'all',
    'year',
    'month',
    'week',
    'day'
);

create table ol_analytics.bucket (
    book_id         int8 not null,
    metric          text not null,
    bucket_type     ol_analytics.bucket_period_type not null,
    bucket_start    timestamptz not null,
    samples_count   int8 not null,
    value_sum       double precision not null,
    updated_at      timestamptz not null default now(),

    primary key (
        book_id,
        metric,
        bucket_type,
        bucket_start
    )
);


create table ol_analytics.book_popularity_bucket (
    book_id      int8 not null,
    bucket_type  ol_analytics.bucket_period_type not null,
    bucket_start timestamptz not null,
    value        double precision not null,
    updated_at   timestamptz not null default now(),

    primary key (
        book_id,
        bucket_type,
        bucket_start
    )
);

create table ol_analytics.worker_state (
    worker_name text primary key,
    last_launch timestamptz not null default '-infinity',
    last_cursor int8 not null default 0
);