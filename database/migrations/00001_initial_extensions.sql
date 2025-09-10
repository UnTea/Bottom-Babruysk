-- +goose Up
-- +goose StatementBegin
create extension if not exists "pgcrypto";
create extension if not exists "pg_trgm";
create extension if not exists "citext";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop extension if exists "pg_trgm";
drop extension if exists "citext";
drop extension if exists "pgcrypto";
-- +goose StatementEnd
