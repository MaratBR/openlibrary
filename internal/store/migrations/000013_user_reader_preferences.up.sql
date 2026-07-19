create table user_reader_preferences (
    user_id uuid primary key references users(id) on delete cascade,
    font_size smallint not null default 18 check (font_size in (12, 14, 16, 18, 20, 22, 26, 30, 36, 42, 48)),
    font_family text not null default 'serif' check (font_family in ('serif', 'sans', 'dyslexic')),
    page_color text not null default 'background' check (page_color in ('background', 'surface')),
    theme text not null default 'system' check (theme in ('system', 'light', 'dark')),
    updated_at timestamptz not null default now()
);
