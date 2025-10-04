-- +goose Up
create type visibility as enum (
    'private',
    'unlisted',
    'public'
    );

comment on type visibility is 'Тип видимости контента: private, unlisted, public';

create type upload_status as enum (
    'pending',
    'processing',
    'done',
    'failed'
    );

comment on type upload_status is 'Тип статуса загрузки: pending, processing, done, failed';

create type user_role as enum (
    'user',
    'admin'
    );

comment on type user_role is 'Роль пользователя: обычный пользователь или администратор.';

create type codec as enum (
    -- Lossless
    'wav',
    'wv',
    'wvc',
    'flac',
    'alac',
    'lpac',
    'ltac',
    'off',
    'ofr',
    'ofs',
    'thd',
    'ape',
    'shn',
    'opus',
    'vorbis',
    'pcm',

    -- Lossy
    'mp3',
    'aac',
    'wma',
    'ogg'
    );

comment on type codec is 'Программные энкодеры используемые используемых для сжатия и распаковки видеоряда или звуковой дорожки';

create type format as enum (
    'mp3',
    'mp4',
    'flac',
    'wav',
    'ogg',
    'ts',
    'm4a',
    'webm',
    'aac'
    );

comment on type format is 'Идентификатор формата файла (.mp3, .mp4)';

-- +goose StatementBegin

-- Users
create table users
(
    id            uuid primary key default gen_random_uuid(),
    email         citext unique                   not null,
    password_hash text                            not null,
    display_name  text                            not null,
    role          user_role        default 'user' not null,
    created_at    timestamptz      default now()  not null,
    updated_at    timestamptz      default null
);

comment on table users is 'Пользователи сервиса.';

comment on column users.id is 'Уникальный идентификатор.';
comment on column users.email is 'Email пользователя.';
comment on column users.password_hash is 'Хеш пароля пользователя.';
comment on column users.display_name is 'Отображаемое имя пользователя.';
comment on column users.role is 'Роль пользователя (user/admin).';
comment on column users.created_at is 'Время создания записи пользователя.';
comment on column users.updated_at is 'Время последнего обновления записи пользователя.';

create index users_email_index on users (email);
comment on index users_email_index is 'Индекс по email таблицы users для ускоренного поиска пользователя по email.';

-- Tracks
create table tracks
(
    id          uuid primary key default gen_random_uuid(),
    uploader_id uuid references users (id),
    title       text                               not null,
    subtitle    text                               not null,
    description text                               not null,
    duration    interval                           not null,
    visibility  visibility       default 'private' not null,
    created_at  timestamptz      default now()     not null,
    updated_at  timestamptz      default null,
    uploaded_at timestamptz                        not null
);

comment on table tracks is 'Логические треки. Одна запись — единый музыкальный/медийный объект.';

comment on column tracks.id is 'Уникальный идентификатор.';
comment on column tracks.uploader_id is 'Пользователь/владелец, загрузивший трек.';
comment on column tracks.title is 'Название трека.';
comment on column tracks.subtitle is 'Подзаголовок или дополнительный заголовок.';
comment on column tracks.description is 'Дополнительное описание трека.';
comment on column tracks.duration is 'Длительность трека.';
comment on column tracks.visibility is 'Видимость трека.';
comment on column tracks.created_at is 'Время создания трека.';
comment on column tracks.updated_at is 'Время последнего обновления трека.';
comment on column tracks.uploaded_at is 'Дата загрузки трека.';

create index tracks_owner_index on tracks (uploader_id);
comment on index tracks_owner_index is 'Индекс по owner_id таблицы tracks для ускорения поиска треков по владельцу.';

-- Person // TODO добавить

-- Person <-> Artists: связь многие ко многим

-- Artists
create table artists
(
    id                      uuid primary key default gen_random_uuid(),
    display_name            text                           not null,
    normalized_display_name text unique                    not null,
    bio                     text                           not null,
    created_at              timestamptz      default now() not null,
    updated_at              timestamptz      default null
);

comment on table artists is 'Исполнители/авторы.';

comment on column artists.id is 'Уникальный идентификатор.';
comment on column artists.display_name is 'Отображаемое имя артиста.';
comment on column artists.normalized_display_name is 'Нормализованное имя артиста';
comment on column artists.bio is 'Короткая биография/описание артиста.';
comment on column artists.created_at is 'Время создания записи артиста.';
comment on column artists.updated_at is 'Время последнего обновления записи артиста.';

-- TrackArtists- связь трек <-> артист
create table track_artists
(
    id         uuid primary key default gen_random_uuid(),
    track_id   uuid                           not null references tracks (id) on delete cascade,
    artist_id  uuid                           not null references artists (id) on delete cascade,
    ord        integer          default 0     not null,
    created_at timestamptz      default now() not null,
    unique (track_id, artist_id)
);

comment on table track_artists is 'Связка трек — автор: role (performer/composer/etc.) и его порядок (ord).';

comment on column track_artists.id is 'Уникальный идентификатор.';
comment on column track_artists.track_id is 'Ссылка на трек.';
comment on column track_artists.artist_id is 'Ссылка на исполнителя.';
comment on column track_artists.ord is 'Порядок/позиция автора в списке авторов трека.';
comment on column track_artists.created_at is 'Время создания связи трек <-> артист.';

-- TrackFiles
create table track_files
(
    id          uuid primary key default gen_random_uuid(),
    track_id    uuid                           not null references tracks (id) on delete cascade,
    filename    text                           not null,
    s3_key      text,
    mime        text                           not null,
    format      format                         not null,
    codec       codec                          not null,
    bitrate     integer                        not null,
    sample_rate integer                        not null,
    channels    integer                        not null,
    size        bigint                         not null,
    duration    interval                       not null,
    checksum    bytea unique                   not null,
    created_at  timestamptz      default now() not null,
    updated_at  timestamptz      default null,
    uploaded_at timestamptz                    not null
);

comment on table track_files is 'Файлы/копии трека в разных форматах (mp3, flac, mp4 и т.п.).';

comment on column track_files.id is 'Уникальный идентификатор.';
comment on column track_files.track_id is 'Ссылка на трек.';
comment on column track_files.filename is 'Исходное имя файла или имя в хранилище.';
comment on column track_files.s3_key is 'Ключ/путь в object storage (S3/MinIO).';
comment on column track_files.mime is 'MIME-тип файла.';
comment on column track_files.format is 'Формат файла (mp3, flac, mp4 и т.д.).';
comment on column track_files.codec is 'Кодек (например "aac", "opus").';
comment on column track_files.bitrate is 'Битрейт в килобитах/с.';
comment on column track_files.sample_rate is 'Частота дискретизации в герцах.';
comment on column track_files.channels is 'Число каналов (1,2 и т.д.).';
comment on column track_files.size is 'Размер файла в байтах.';
comment on column track_files.duration is 'Длительность файла.';
comment on column track_files.checksum is 'Контрольная сумма файла (md5/sha256 и т.п.).';
comment on column track_files.created_at is 'Время создания записи файла трека.';
comment on column track_files.updated_at is 'Время последнего обновления файла трека.';
comment on column track_files.uploaded_at is 'Дата загрузки файла трека.';

create index track_files_track_id_index on track_files (track_id);
comment on index track_files_track_id_index is 'Индекс по track_id таблицы track_files для ускорения поиска файлов по треку.';

-- Albums
create table albums
(
    id           uuid primary key default gen_random_uuid(),
    owner_id     uuid                               references users (id) on delete set null,
    title        text                               not null,
    description  text                               not null,
    visibility   visibility       default 'private' not null,
    release_date date                               not null,
    created_at   timestamptz      default now()     not null,
    updated_at   timestamptz      default null
);

comment on table albums is 'Альбомы / релизы.';

comment on column albums.id is 'Уникальный идентификатор.';
comment on column albums.owner_id is 'Кто создал запись альбома (может быть NULL, если удалён пользователь).';
comment on column albums.title is 'Название альбома.';
comment on column albums.description is 'Описание альбома.';
comment on column albums.visibility is 'Видимость альбома.';
comment on column albums.release_date is 'Дата релиза (если известна).';
comment on column albums.created_at is 'Время создания записи альбома.';
comment on column albums.updated_at is 'Время последнего обновления записи альбома.';

create index albums_owner_index on albums (owner_id);
comment on index albums_owner_index is 'Индекс по owner_id таблицы albums для ускорения поиска альбомов по владельцу.';

-- AlbumTracks — позиции треков внутри альбома/диска
create table album_tracks
(
    album_id    uuid    not null references albums (id) on delete cascade,
    disc_number integer not null default 1,
    position    integer not null,
    track_id    uuid    not null references tracks (id) on delete cascade,
    created_at  timestamptz      default now(),
    primary key (album_id, disc_number, position),
    unique (album_id, disc_number, track_id)
);

comment on table album_tracks is 'Позиции треков в альбоме: диск + позиция.';

comment on column album_tracks.album_id is 'Ссылка на альбом.';
comment on column album_tracks.disc_number is 'Номер диска внутри альбома.';
comment on column album_tracks.position is 'Позиция трека на диске (1..N).';
comment on column album_tracks.track_id is 'Ссылка на трек.';
comment on column album_tracks.created_at is 'Время создания позиции трека в альбоме.';

-- Playlists
create table playlists
(
    id          uuid primary key default gen_random_uuid(),
    owner_id    uuid references users (id) on delete cascade,
    title       text                               not null,
    description text                               not null,
    visibility  visibility       default 'private' not null,
    created_at  timestamptz      default now()     not null,
    updated_at  timestamptz      default null
);

comment on table playlists is 'Плейлисты пользователей.';

comment on column playlists.id is 'Уникальный идентификатор.';
comment on column playlists.owner_id is 'Владелец плейлиста.';
comment on column playlists.title is 'Название плейлиста.';
comment on column playlists.description is 'Описание плейлиста.';
comment on column playlists.visibility is 'Видимость плейлиста.';
comment on column playlists.created_at is 'Время создания плейлиста.';
comment on column playlists.updated_at is 'Время последнего обновления плейлиста.';

create index playlists_owner_index on playlists (owner_id);
comment on index playlists_owner_index is 'Индекс по owner_id таблицы playlists для ускорения поиска плейлистов по владельцу.';

-- PlaylistItems
create table playlist_items
(
    playlist_id uuid references playlists (id) on delete cascade,
    track_id    uuid references tracks (id) on delete cascade,
    position    integer                   not null,
    added_at    timestamptz default now() not null,
    primary key (playlist_id, track_id)
);

comment on table playlist_items is 'Связь плейлист -> трек с позицией.';

comment on column playlist_items.playlist_id is 'Ссылка на плейлист.';
comment on column playlist_items.track_id is 'Ссылка на трек.';
comment on column playlist_items.position is 'Позиция трека в плейлисте.';
comment on column playlist_items.added_at is 'Время добавления трека в плейлист.';

-- Uploads
create table uploads
(
    id         uuid primary key default gen_random_uuid(),
    owner_id   uuid references users (id),
    filename   text                               not null,
    s3_key     text                               not null,
    mime       text                               not null,
    size       bigint                             not null,
    status     upload_status    default 'pending' not null,
    created_at timestamptz      default now()     not null,
    updated_at timestamptz      default null
);

comment on table uploads is 'Промежуточные записи загрузок.';

comment on column uploads.id is 'Уникальный идентификатор.';
comment on column uploads.owner_id is 'Пользователь, загрузивший файл.';
comment on column uploads.filename is 'Исходное имя загружаемого файла.';
comment on column uploads.s3_key is 'Ключ/путь в object storage (S3/MinIO).';
comment on column uploads.mime is 'MIME-тип файла.';
comment on column uploads.size is 'Размер файла в байтах.';
comment on column uploads.status is 'Статус загрузки: pending, processing, done, failed.';
comment on column uploads.created_at is 'Время создания записи загрузки.';
comment on column uploads.updated_at is 'Время последнего обновления записи загрузки.';

-- UserFollowers
create table user_followers
(
    follower_id uuid                      not null references users (id) on delete cascade,
    followee_id uuid                      not null references users (id) on delete cascade,
    created_at  timestamptz default now() not null,
    primary key (follower_id, followee_id)
);

comment on table user_followers is 'Подписки пользователей друг на друга.';

comment on column user_followers.follower_id is 'Пользователь, который подписывается (инициатор подписки).';
comment on column user_followers.followee_id is 'Пользователь, на которого оформлена подписка.';
comment on column user_followers.created_at is 'Дата и время создания подписки.';

-- TrackLikes
create table track_likes
(
    user_id    uuid                      not null references users (id) on delete cascade,
    track_id   uuid                      not null references tracks (id) on delete cascade,
    created_at timestamptz default now() not null,
    primary key (user_id, track_id)
);

comment on table track_likes is 'Лайки треков пользователями.';

comment on column track_likes.user_id is 'Пользователь, который поставил лайк.';
comment on column track_likes.track_id is 'Трек, которому поставлен лайк.';
comment on column track_likes.created_at is 'Дата и время, когда был поставлен лайк.';

-- pg_trgm индекс для поиска по названию
create index tracks_title_trgm_index on tracks using gin (title gin_trgm_ops);
comment on index tracks_title_trgm_index is 'GIN trigram индекс на столбце title таблицы tracks для ускоренного поиска по названию.';

-- Уникальный индекс: для одного трека — одно представление в конкретном формате
create unique index track_files_track_format_unique_index on track_files (track_id, format);
comment on index track_files_track_format_unique_index is 'Уникальный индекс для таблицы track_files по полям (track_id, format), гарантирующий, что один трек не имеет двух файлов одинакового формата.';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index track_files_track_format_unique_index;
drop index tracks_title_trgm_index;

drop table track_likes;
drop table user_followers;
drop table uploads;
drop table playlist_items;

drop index playlists_owner_index;
drop table playlists;

drop table album_tracks;

drop index albums_owner_index;
drop table albums;

drop index track_files_track_id_index;
drop table track_files;

drop table track_artists;
drop table artists;

drop index tracks_owner_index;
drop table tracks;

drop index users_email_index;
drop table users;

-- +goose StatementEnd

drop type upload_status;
drop type visibility;
drop type user_role;
