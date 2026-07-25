alter table user_logs
    alter column action_type type text using action_type::text,
    add column "time" timestamptz not null default now(),
    add column reason text not null default '';

create index ix_user_logs_user_time on user_logs (user_id, "time" desc);

create table comment_logs (
    id int8 primary key,
    "time" timestamptz not null default now(),
    comment_id int8 not null references comments(id),
    action_type text not null,
    payload jsonb null,
    actor_user_id uuid null references users(id),
    reason text not null default ''
);

create index ix_comment_logs_comment_time on comment_logs (comment_id, "time" desc);
