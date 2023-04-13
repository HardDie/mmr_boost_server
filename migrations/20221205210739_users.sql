-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id           SERIAL    PRIMARY KEY,
    email        TEXT      NOT NULL UNIQUE,
    username     TEXT      NOT NULL UNIQUE,
    role_id      INT       NOT NULL DEFAULT 4,
    is_activated BOOLEAN   NOT NULL DEFAULT false,
    created_at   TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at   TIMESTAMP NOT NULL DEFAULT (now()),
    deleted_at   TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
