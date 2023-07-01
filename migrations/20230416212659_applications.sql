-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS applications (
    id             SERIAL   PRIMARY KEY,
    user_id        INT      NOT NULL           REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    status_id      INT      NOT NULL,
    type_id        INT      NOT NULL,
    current_mmr    INT      NOT NULL,
    target_mmr     INT      NOT NULL,
    tg_contact     TEXT     NOT NULL,
    price          INT      NOT NULL,
    comment        TEXT,
    steam_login    TEXT,
    steam_password TEXT,
    created_at     TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at     TIMESTAMP NOT NULL DEFAULT (now()),
    deleted_at     TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE applications;
-- +goose StatementEnd
