-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS status_history (
    id             SERIAL    PRIMARY KEY,
    user_id        INT       NOT NULL REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    application_id INT       NOT NULL REFERENCES applications(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    new_status_id  INT       NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT (now())
);
CREATE INDEX status_history_application_id_idx ON status_history (application_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE status_history;
-- +goose StatementEnd
