create schema ol_analytics;

create type ol_analytics.interaction_event_type as enum (
    'book_view',
    'chapter_view',
    'started_reading',
    'completed',
    'dropped',
    'finished_chapter',
    'search_click'
);

create table ol_analytics.interaction_event (
    id bigint generated always as identity primary key,
    user_key text not null,
    book_id int8 not null,
    event_type ol_analytics.interaction_event_type not null,
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

create type ol_analytics.counter_type as enum (
    'views',
    'search_clicks'
);

create table ol_analytics.bucket_counter (
    book_id      int8 not null,
    metric       ol_analytics.counter_type not null,
    bucket_type  ol_analytics.bucket_period_type not null,
    bucket_start timestamptz not null,

    value        int8 not null,
    updated_at   timestamptz not null default now(),

    primary key (
        book_id,
        metric,
        bucket_type,
        bucket_start
    )
);


create table ol_analytics.bucket_popularity (
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
