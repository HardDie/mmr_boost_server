-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS history (
    id             SERIAL    PRIMARY KEY,
    user_id        INT       NOT NULL,
    affect_user_id INT,
    message        TEXT      NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT (now())
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE history;
-- +goose StatementEnd
