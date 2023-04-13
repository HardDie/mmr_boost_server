-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS access_tokens (
    id         SERIAL    PRIMARY KEY,
    user_id    INT       NOT NULL UNIQUE REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    token_hash TEXT      NOT NULL UNIQUE,
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    deleted_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE access_tokens;
-- +goose StatementEnd
