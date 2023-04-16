-- +goose Up
-- +goose StatementBegin
INSERT INTO application_statuses (id, title)
VALUES (1, 'created'),
       (2, 'waiting for payment'),
       (3, 'payed'),
       (4, 'waiting for processing'),
       (5, 'in progress'),
       (6, 'paused'),
       (7, 'done')
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM application_statuses WHERE id IN (1, 2, 3, 4, 5, 6, 7);
-- +goose StatementEnd
