-- +goose Up
-- +goose StatementBegin
INSERT INTO application_statuses (id, title)
VALUES (1, 'created'),
       (2, 'awaits payment'),
       (3, 'paid'),
       (4, 'in progress'),
       (5, 'done'),
       (6, 'deleted'),
       (7, 'canceled'),
       (8, 'suspended')
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM application_statuses WHERE id IN (1, 2, 3, 4, 5, 6, 7);
-- +goose StatementEnd
