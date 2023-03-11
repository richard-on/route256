-- +goose Up
-- +goose StatementBegin
INSERT INTO stocks VALUES (1625903, 1, 3);
INSERT INTO stocks VALUES (1625903, 2, 2);
INSERT INTO stocks VALUES (1625903, 3, 1);
INSERT INTO stocks VALUES (147298765, 4, 5);
INSERT INTO stocks VALUES (147298765, 1, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM stocks;
-- +goose StatementEnd
