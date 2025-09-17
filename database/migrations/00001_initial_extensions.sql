-- +goose Up
-- +goose StatementBegin
create extension "pgcrypto";
create extension "pg_trgm";
create extension "citext";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop extension "pg_trgm";
drop extension "citext";
drop extension "pgcrypto";
-- +goose StatementEnd
