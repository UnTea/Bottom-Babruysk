-- +goose Up
create type visibility as enum ('private', 'unlisted', 'public');
comment on type visibility is 'Тип видимости контента: private, unlisted, public';

create type upload_status as enum ('pending', 'processing', 'done', 'failed');
comment on type upload_status is 'Тип статуса загрузки: pending, processing, done, failed';

create type user_role as enum ('user', 'admin');
comment on type user_role is 'Роль пользователя: обычный пользователь или администратор.';

-- +goose StatementBegin

-- Users
create table if not exists users
(
    id            uuid primary key default gen_random_uuid(),
    email         citext unique not null,
    password_hash text          not null,
    display_name  text,
    role          user_role        default 'user',
    created_at    timestamptz      default now(),
    updated_at    timestamptz      default now()
);

comment on table users is 'Пользователи сервиса.';

comment on column users.id is 'Уникальный идентификатор.';
comment on column users.email is 'Email пользователя.';
comment on column users.password_hash is 'Хеш пароля пользователя.';
comment on column users.display_name is 'Отображаемое имя пользователя.';
comment on column users.role is 'Роль пользователя (user/admin).';
comment on column users.created_at is 'Время создания записи пользователя.';
comment on column users.updated_at is 'Время последнего обновления записи пользователя.';

create index if not exists users_email_index on users (email);
comment on index users_email_index is 'Индекс по email таблицы users для ускоренного поиска пользователя по email.';

-- Tracks
create table if not exists tracks
(
    id          uuid primary key default gen_random_uuid(),
    uploader_id uuid references users (id) on delete cascade,
    title       text        not null,
    subtitle    text,
    description text,
    duration    interval,
    visibility  visibility       default 'private',
    created_at  timestamptz      default now(),
    updated_at  timestamptz      default now(),
    uploaded_at timestamptz not null
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

create index if not exists tracks_owner_index on tracks (uploader_id);
comment on index tracks_owner_index is 'Индекс по owner_id таблицы tracks для ускорения поиска треков по владельцу.';

-- Artists
create table if not exists artists
(
    id         uuid primary key default gen_random_uuid(),
    name       text not null,
    bio        text,
    created_at timestamptz      default now(),
    updated_at timestamptz      default now()
);

comment on table artists is 'Исполнители/авторы.';

comment on column artists.id is 'Уникальный идентификатор.';
comment on column artists.name is 'Отображаемое имя артиста.';
comment on column artists.bio is 'Короткая биография/описание артиста.';
comment on column artists.created_at is 'Время создания записи артиста.';
comment on column artists.updated_at is 'Время последнего обновления записи артиста.';

-- TrackAuthors - связь трек <-> артист
create table if not exists track_authors
(
    id         uuid primary key default gen_random_uuid(),
    track_id   uuid not null references tracks (id) on delete cascade,
    artist_id  uuid not null references artists (id) on delete cascade,
    role       text             default 'performer',
    ord        integer          default 0,
    created_at timestamptz      default now(),
    unique (track_id, artist_id, ord)
);

comment on table track_authors is 'Связка трек — автор: role (performer/composer/etc.) и его порядок (ord).';

comment on column track_authors.id is 'Уникальный идентификатор.';
comment on column track_authors.track_id is 'Ссылка на трек.';
comment on column track_authors.artist_id is 'Ссылка на исполнителя.';
comment on column track_authors.role is 'Роль исполнителя для данного трека (performer, composer и т.д.).';
comment on column track_authors.ord is 'Порядок/позиция автора в списке авторов трека.';
comment on column track_authors.created_at is 'Время создания связи трек <-> артист.';

-- TrackFiles
create table if not exists track_files
(
    id          uuid primary key default gen_random_uuid(),
    track_id    uuid        not null references tracks (id) on delete cascade,
    filename    text        not null,
    s3_key      text,
    mime        text,
    format      text,
    codec       text,
    bitrate     integer,
    sample_rate integer,
    channels    integer,
    size        bigint,
    duration    interval,
    checksum    text,
    created_at  timestamptz      default now(),
    updated_at  timestamptz      default now(),
    uploaded_at timestamptz not null
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

create index if not exists track_files_track_id_index on track_files (track_id);
comment on index track_files_track_id_index is 'Индекс по track_id таблицы track_files для ускорения поиска файлов по треку.';

-- Albums
create table if not exists albums
(
    id           uuid primary key default gen_random_uuid(),
    owner_id     uuid references users (id) on delete set null,
    title        text not null,
    description  text,
    release_date date,
    created_at   timestamptz      default now(),
    updated_at   timestamptz      default now()
);

comment on table albums is 'Альбомы / релизы.';

comment on column albums.id is 'Уникальный идентификатор.';
comment on column albums.owner_id is 'Кто создал запись альбома (может быть NULL, если удалён пользователь).';
comment on column albums.title is 'Название альбома.';
comment on column albums.description is 'Описание альбома.';
comment on column albums.release_date is 'Дата релиза (если известна).';
comment on column albums.created_at is 'Время создания записи альбома.';
comment on column albums.updated_at is 'Время последнего обновления записи альбома.';

create index if not exists albums_owner_index on albums (owner_id);
comment on index albums_owner_index is 'Индекс по owner_id таблицы albums для ускорения поиска альбомов по владельцу.';

-- AlbumDiscs
create table if not exists album_discs
(
    album_id    uuid    not null references albums (id) on delete cascade,
    disc_number integer not null default 1,
    primary key (album_id, disc_number)
);

comment on table album_discs is 'Диски в альбоме (номер диска внутри релиза).';

comment on column album_discs.album_id is 'Ссылка на альбом.';
comment on column album_discs.disc_number is 'Номер диска внутри альбома.';

-- AlbumTracks — позиции треков внутри альбома/диска
create table if not exists album_tracks
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
create table if not exists playlists
(
    id          uuid primary key default gen_random_uuid(),
    owner_id    uuid references users (id) on delete cascade,
    title       text not null,
    description text,
    visibility  visibility       default 'private',
    created_at  timestamptz      default now(),
    updated_at  timestamptz      default now()
);

comment on table playlists is 'Плейлисты пользователей.';

comment on column playlists.id is 'Уникальный идентификатор.';
comment on column playlists.owner_id is 'Владелец плейлиста.';
comment on column playlists.title is 'Название плейлиста.';
comment on column playlists.description is 'Описание плейлиста.';
comment on column playlists.visibility is 'Видимость плейлиста.';
comment on column playlists.created_at is 'Время создания плейлиста.';
comment on column playlists.updated_at is 'Время последнего обновления плейлиста.';

create index if not exists playlists_owner_index on playlists (owner_id);
comment on index playlists_owner_index is 'Индекс по owner_id таблицы playlists для ускорения поиска плейлистов по владельцу.';

-- PlaylistItems
create table if not existS playlist_items
(
    playlist_id uuid references playlists (id) on delete cascade,
    track_id    uuid references tracks (id) on delete cascade,
    pos         integer not null,
    added_at    timestamptz default now(),
    primary key (playlist_id, track_id)
);

comment on table playlist_items is 'Связь плейлист -> трек с позицией.';

comment on column playlist_items.playlist_id is 'Ссылка на плейлист.';
comment on column playlist_items.track_id is 'Ссылка на трек.';
comment on column playlist_items.pos is 'Позиция трека в плейлисте.';
comment on column playlist_items.added_at is 'Время добавления трека в плейлист.';

-- Uploads
create table if not exists uploads
(
    id                uuid primary key default gen_random_uuid(),
    owner_id          uuid references users (id),
    original_filename text,
    s3_key            text,
    mime              text,
    size              bigint,
    status            upload_status    default 'pending',
    created_at        timestamptz      default now(),
    updated_at        timestamptz      default now()
);

comment on table uploads is 'Промежуточные записи загрузок.';

comment on column uploads.id is 'Уникальный идентификатор.';
comment on column uploads.owner_id is 'Пользователь, загрузивший файл.';
comment on column uploads.original_filename is 'Исходное имя загружаемого файла.';
comment on column uploads.s3_key is 'Ключ/путь в object storage (S3/MinIO).';
comment on column uploads.mime is 'MIME-тип файла.';
comment on column uploads.size is 'Размер файла в байтах.';
comment on column uploads.status is 'Статус загрузки: pending, processing, done, failed.';
comment on column uploads.created_at is 'Время создания записи загрузки.';
comment on column uploads.updated_at is 'Время последнего обновления записи загрузки.';

-- UserFollowers
create table if not exists user_followers
(
    follower_id uuid not null references users (id) on delete cascade,
    followee_id uuid not null references users (id) on delete cascade,
    created_at  timestamptz default now(),
    primary key (follower_id, followee_id)
);

comment on table user_followers is 'Подписки пользователей друг на друга.';

comment on column user_followers.follower_id is 'Пользователь, который подписывается (инициатор подписки).';
comment on column user_followers.followee_id is 'Пользователь, на которого оформлена подписка.';
comment on column user_followers.created_at is 'Дата и время создания подписки.';

-- TrackLikes
create table if not exists track_likes
(
    user_id    uuid not null references users (id) on delete cascade,
    track_id   uuid not null references tracks (id) on delete cascade,
    created_at timestamptz default now(),
    primary key (user_id, track_id)
);

comment on table track_likes is 'Лайки треков пользователями.';

comment on column track_likes.user_id is 'Пользователь, который поставил лайк.';
comment on column track_likes.track_id is 'Трек, которому поставлен лайк.';
comment on column track_likes.created_at is 'Дата и время, когда был поставлен лайк.';

-- pg_trgm индекс для поиска по названию
create index if not exists tracks_title_trgm_index on tracks using gin (title gin_trgm_ops);
comment on index tracks_title_trgm_index is 'GIN trigram индекс на столбце title таблицы tracks для ускоренного поиска по названию.';

-- Уникальный индекс: для одного трека — одно представление в конкретном формате
create unique index if not exists track_files_track_format_unique_index on track_files (track_id, format);
comment on index track_files_track_format_unique_index is 'Уникальный индекс для таблицы track_files по полям (track_id, format), гарантирующий, что один трек не имеет двух файлов одинакового формата.';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists track_files_track_format_unique_index;
drop index if exists tracks_title_trgm_index;

drop table if exists track_likes;
drop table if exists user_followers;
drop table if exists uploads;
drop table if exists playlist_items;

drop index if exists playlists_owner_index;
drop table if exists playlists;

drop table if exists album_tracks;
drop table if exists album_discs;

drop index if exists albums_owner_index;
drop table if exists albums;

drop index if exists track_files_track_id_index;
drop table if exists track_files;

drop table if exists track_authors;
drop table if exists artists;

drop index if exists tracks_owner_index;
drop table if exists tracks;

drop index if exists users_email_index;
drop table if exists users;

-- +goose StatementEnd

drop type upload_status;
drop type visibility;
drop type user_role;
