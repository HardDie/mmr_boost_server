-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reset_password (
    id         SERIAL    PRIMARY KEY,
    user_id    INT       NOT NULL UNIQUE REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    code       TEXT      NOT NULL UNIQUE,
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now())
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reset_password;
-- +goose StatementEnd
