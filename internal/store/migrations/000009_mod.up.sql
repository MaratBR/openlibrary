alter table books
    add column is_perm_removed bool not null default false;

create table moderation_logs (
    id int8 primary key,
    "time" timestamptz not null,
    "type" text not null,
    target_type text not null,
    target_id text not null,
    payload jsonb null,
    actor_user_id uuid null references users (id),
    reason text not null default ''
);

create index ix_moderation_logs_target_time
    on moderation_logs (target_type, target_id, "time" desc);
create index ix_moderation_logs_target_type_time
    on moderation_logs (target_type, target_id, "type", "time" desc);

create table reports (
    id bigint generated always as identity primary key,
    number text not null unique,
    "time" timestamptz not null,
    reporter_user_id uuid null references users (id),
    target_type text not null check (target_type in ('user', 'book', 'comment')),
    target_id text not null,
    reason text not null,
    description text not null
);

create table report_number_counters (
    "day" date primary key,
    counter bigint not null check (counter > 0)
);

create index ix_reports_target_time
    on reports (target_type, target_id, "time" desc);
