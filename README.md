### Bottom Babruysk

#### Project requirements

- Go 1.24+
- Docker
- Docker Compose
- Taskfile (go install github.com/go-task/task/v3/cmd/task@latest)

#### Quick start

1. Copy .env (see example below) to the root of the repository.
2. Start Postgres and run migrations:

```bash
    task compose:up:detached
```

3. Build both binaries:

```bash
    task build:api
    task build:imposter
```

Binaries will appear in `./bin/` folder:

- ./bin/bottom_babruysk.exe — HTTP API
- ./bin/imposter.exe — FLAC importer

#### Docker Compose (DB + migrations)

The repository contains `docker-compose.yml`, which launches:

- `postgres` (17.6-alpine) on `5432` port
- `migrator` (Dockerfile: `migrator.dockerfile`), which reads `.env` and runs goose up from the `/migrations` folder (
  externally mounted `./migrations:/migrations:ro`).

Useful commands:

```bash
    # start and leave running in the background
    task compose:up:detached
    
    # stop and remove containers (volume with data will remain)
    task compose:down
```

You can see other useful commands by simply calling without parameters:

```bash
  task
```

#### Environment variables (.env)

Example `.env` (place next to `docker-compose.yml` and `Taskfile.yml`):

```text
# Database
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=bottom_babruysk
DB_SSLMODE=disable

# HTTP (API)
HTTP_ADDRESS=:42069
API_LOGGER_LEVEL=debug

# Goose (для контейнера migrator)
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://admin:admin@postgres:5432/bottom_babruysk?sslmode=disable
GOOSE_MIGRATION_DIR=/migrations
GOOSE_TABLE=goose_migrations
```

**Note**: The API usually collects DSN from DB_\*, while the migrator inside the container uses GOOSE_\*.

#### API assembly and launch

Build:

```bash
  task build:api
```

(Binary: `./bin/bottom_babruysk.exe`, source code: `cmd/api`)

Run:

```bash
    task run:api
    # or simply ./bin/bottom_babruysk.exe
```

- Listens to the address from HTTP_ADDRESS (default: 42069).
- Connects to the database via DB_* from .env.

#### FLAC Importer (CLI)

Build:

```bash
  task build:imposter
```

(Binary: `./bin/imposter.exe`, source code: `cmd/flac_importer`)

Run:

```bash
    ./bin/imposter.exe \
      --path "C:\Music" \
      --dsn "postgres://user:password@db_ip:db_port/bottom_babruysk?sslmode=disable" \
      --uploader-email importer@local
```

Parameters:

- --path, -p - root folder for scanning `.flac` files
- --dsn, -d - Postgres `DSN`
- --uploader-email - email address of the user who owns the imported tracks (default is `imposter@sosal.cock`). If such
  a user does not exist, they will be created.

**What the importer does:**

- Reads FLAC metadata (`STREAMINFO`, `VorbisComment`, etc.).
- Determines the fields: `title/album/artist/date` (fallback: from the file name, ‘Unknown Album/Artist’).
- Calculates the bitrate using the formula `sampleRate * bitsPerSample * channels / 1000`.
- Deduplication:
    - track_files - by `checksum` (FLAC MD5 from `STREAMINFO`). In the database, track_files.checksum must be `unique`.
    - tracks - by pair (artist + title) (normalised artist — see below — and case-insensitive title).
- Artists:
  - Stored as `display_name` and `normalised_name` (`unique`). 
  - Normalisation: remove all special characters and spaces, leave only letters and numbers, convert to lowercase. 
- Creates a link in `track_artists (unique track_id + artist_id)`. 
- Writes the track file to `track_files` (mimetype/format/codec/bitrate/frequency/channels/size/duration/checksum/download
time).

#### DSN Examples

Locally (host machine):
```text
    postgres://admin:admin@127.0.0.1:5432/bottom_babruysk?sslmode=disable
```

Inside the docker-compose network (if running inside a container):
```text
    postgres://admin:admin@postgres:5432/bottom_babruysk?sslmode=disable
```
