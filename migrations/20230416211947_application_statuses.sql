-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS application_statuses (
    id    SERIAL PRIMARY KEY,
    title TEXT   NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE application_statuses;
-- +goose StatementEnd
