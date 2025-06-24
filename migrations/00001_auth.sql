-- +goose Up
-- +goose StatementBegin
create table if not exists users (
  id TEXT unique,
  name TEXT,
  email TEXT unique,
  password_hash TEXT unique,
  is_active boolean,
  created_at timestamp,
  updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
