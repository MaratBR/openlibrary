alter table books 
    add column is_perm_removed bool not null default false;

create type book_action_type as enum (
    -- normal updates
    'significant_update',
    'author_transfer',
    'coauthor_added',
    'coauthor_removed',

    -- moderation actions
    'ban',
    'shadow_ban',
    'perm_removal',
    'un_ban',
    'un_shadow_ban',

    -- special actions
    'reindex'
);

create table book_logs (
    id int8 primary key,
    "time" timestamptz not null,
    book_id int8 not null references books (id),
    action_type book_action_type not null,
    payload jsonb null,
    actor_user_id uuid null references users (id),
    reason text not null default ''
);

create type user_action_type as enum (
    -- security actions
    'sec_password_reset',
    'sec_2fa_umbrella', -- umbrella event for all 2fa events
    
    -- moderation actions
    'ban',
    'unban',
    'mute',
    'unmute'
);

create table user_logs (
    id int8 primary key,
    user_id uuid null references users (id),
    actor_user_id uuid null references users (id),
    action_type user_action_type not null,
    payload jsonb null
);