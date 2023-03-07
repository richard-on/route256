-- +goose Up
-- +goose StatementBegin
CREATE TYPE "item" AS
(
    sku   BIGINT,
    count INT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS "item";
-- +goose StatementEnd
