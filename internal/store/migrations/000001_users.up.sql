create table users (
    id uuid primary key,
    name varchar(255) not null,
    joined_at timestamptz not null default now(),
    password_hash text not null
);

create table "sessions" (
    id text primary key,
    user_id uuid not null references users(id),
    created_at timestamptz not null default now(),
    user_agent text not null,
    ip_address text not null,
    expires_at timestamptz not null
);