-- +goose Up
-- +goose StatementBegin

alter table albums
    add column picture_id uuid;

comment on column albums.picture_id is 'Идентификатор записи в uploads, представляющей обложку альбома (оригинал). При удалении загрузки обложка сбрасывается в NULL.';


alter table albums
    add constraint albums_picture_id_fk
        foreign key (picture_id) references uploads (id) on delete set null;

-- Genres
create table genres
(
    id   uuid primary key default gen_random_uuid(),
    name text not null unique
);

comment on table genres is 'Справочник музыкальных жанров.';

comment on column genres.id is 'Уникальный идентификатор жанра.';
comment on column genres.name is 'Уникальное имя жанра.';

-- TrackGenres
create table track_genres
(
    track_id uuid not null references tracks (id) on delete cascade,
    genre_id uuid not null references genres (id) on delete cascade,
    primary key (track_id, genre_id)
);

comment on table track_genres is 'Связь треков с жанрами (многие-ко-многим).';

comment on column track_genres.track_id is 'Ссылка на трек.';
comment on column track_genres.genre_id is 'Ссылка на жанр.';

create index track_genres_track_id_idx on track_genres (track_id);
comment on index track_genres_track_id_idx is 'Индекс по track_id для быстрого получения жанров трека.';

create index track_genres_genre_id_idx on track_genres (genre_id);
comment on index track_genres_genre_id_idx is 'Индекс по genre_id для выборки треков по жанру.';

-- AlbumArtists: связь альбом <-> артист
create table album_artists
(
    album_id   uuid                      not null references albums (id) on delete cascade,
    artist_id  uuid                      not null references artists (id) on delete cascade,
    ord        integer     default 0     not null,
    created_at timestamptz default now() not null,
    primary key (album_id, artist_id)
);

comment on table album_artists is 'Связь альбом—артист с порядком отображения.';

comment on column album_artists.album_id is 'Ссылка на альбом.';
comment on column album_artists.artist_id is 'Ссылка на артиста.';
comment on column album_artists.ord is 'Порядок артиста в списке артистов альбома.';
comment on column album_artists.created_at is 'Время создания связи альбом—артист.';

create index album_artists_album_id_idx on album_artists (album_id);
comment on index album_artists_album_id_idx is 'Индекс для быстрого получения артистов конкретного альбома.';

create index album_artists_artist_id_idx on album_artists (artist_id);
comment on index album_artists_artist_id_idx is 'Индекс для быстрого получения альбомов конкретного артиста.';

create index artists_name_trgm_index on artists using gin (name gin_trgm_ops);
comment on index artists_name_trgm_index is 'GIN trigram индекс для быстрого поиска/подсказок по имени артиста.';

create index albums_title_trgm_index on albums using gin (title gin_trgm_ops);
comment on index albums_title_trgm_index is 'GIN trigram индекс для быстрого поиска/подсказок по названию альбома.';

create index album_tracks_track_id_idx on album_tracks (track_id);
comment on index album_tracks_track_id_idx is 'Индекс для быстрого поиска альбомов по треку (обратные связи).';

create index uploads_owner_id_idx on uploads (owner_id);
comment on index uploads_owner_id_idx is 'Индекс для выборок загрузок по владельцу.';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop index uploads_owner_id_idx;
drop index album_tracks_track_id_idx;
drop index albums_title_trgm_index;
drop index artists_name_trgm_index;

drop index album_artists_artist_id_idx;
drop index album_artists_album_id_idx;
drop table album_artists;

drop index track_genres_genre_id_idx;
drop index track_genres_track_id_idx;
drop table track_genres;
drop table genres;

alter table albums
    drop column picture_id;

alter table albums
    drop constraint albums_picture_id_fk;

-- +goose StatementEnd