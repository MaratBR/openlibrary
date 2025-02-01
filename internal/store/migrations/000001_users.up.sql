create type user_role as enum (
    'user',
    'admin',
    'moderator',
    'system'
);

create type censor_mode as enum (
    'hide',
    'censor',
    'none'
);

create table users (
    id uuid primary key,
    name varchar(255) not null,
    joined_at timestamptz not null default now(),
    password_hash text not null,
    "role" user_role not null default 'user',
    is_banned boolean not null default false,
    avatar_file text null,

    -- about columns
    about text not null default '',
    gender text not null default '',

    -- customization columns
    profile_css text not null default '',
    enable_profile_css boolean not null default false,
    default_theme text not null default '',

    -- privacy columns
    privacy_hide_stats boolean not null default false,
    privacy_hide_comments boolean not null default false,
    privacy_hide_email boolean not null default true,
    privacy_allow_searching boolean not null default false,

    -- moderation settings
    show_adult_content boolean not null default false,
    censored_tags text[] not null default '{}',
    censored_tags_mode censor_mode not null default 'none'
);

create type type_of_2fa as enum (
    'totp',
    'webauthn'
);

create table user_2fa (
    id uuid primary key,
    user_id uuid not null references users(id),
    "type" type_of_2fa not null,
    "key" text not null,
    created_at timestamptz not null default now(),
    initialized boolean not null default false,
    active boolean not null default true
);

create table user_bans (
    id bigint primary key,
    user_id uuid not null references users(id),
    created_at timestamptz not null default now(),
    banned_by_user_id uuid null references users(id),
    note text not null default '',
    expires_at timestamptz not null
);

create table user_follower (
    follower_id uuid not null references users(id),
    followed_id uuid not null references users(id),
    created_at timestamptz not null default now(),
    primary key (follower_id, followed_id)
);

create table "sessions" (
    id int8 primary key,
    "sid" text not null,
    user_id uuid not null references users(id),
    created_at timestamptz not null default now(),
    user_agent text not null,
    ip_address text not null,
    expires_at timestamptz not null,
    is_terminated boolean not null default false
);