-- extensions
create extension if not exists "pgcrypto";
create extension if not exists "pg_trgm";
create extension if not exists "citext";

-- visibility
do $$
begin
  if not exists (select 1 from pg_type where typname = 'visibility_t') then
create type visibility_t as enum ('private', 'unlisted', 'public');
end if;
end$$;

-- users
create table if not exists users (
                                     id            uuid primary key default gen_random_uuid(),
    email         citext unique not null,
    password_hash text not null,
    display_name  text,
    created_at    timestamptz default now(),
    updated_at    timestamptz default now()
    );

-- tracks
create table if not exists tracks (
                                      id                      uuid primary key default gen_random_uuid(),
    owner_id                uuid references users(id) on delete cascade,
    title                   text not null,
    artist                  text,
    album                   text,
    filename                text not null,
    s3_key                  text,
    mime                    text,
    size                    bigint,
    duration_seconds        numeric,
    is_processed            boolean default false,
    visibility              visibility_t default 'private',
    created_at              timestamptz default now(),
    updated_at              timestamptz default now()
    );

-- playlists
create table if not exists playlists (
                                         id          uuid primary key default gen_random_uuid(),
    owner_id    uuid references users(id) on delete cascade,
    title       text not null,
    description text,
    visibility  visibility_t default 'private',
    created_at  timestamptz default now(),
    updated_at  timestamptz default now()
    );

-- playlist elements
create table if not exists playlist_items (
                                              playlist_id uuid references playlists(id) on delete cascade,
    track_id    uuid references tracks(id) on delete cascade,
    pos         integer not null,
    added_at    timestamptz default now(),
    PRIMARY KEY (playlist_id, track_id)
    );

-- uploads
create table if not exists uploads (
                                       id                uuid primary key default gen_random_uuid(),
    owner_id          uuid references users(id),
    original_filename text,
    s3_key            text,
    mime              text,
    size              bigint,
    status            text default 'pending', -- pending, processing, done, failed
    created_at        timestamptz default now(),
    updated_at        timestamptz default now()
    );

-- indexes
create index if not exists idx_tracks_owner on tracks(owner_id);
create index if not exists idx_playlists_owner on playlists(owner_id);
create index if not exists idx_users_email on users(email);
