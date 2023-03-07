-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "stocks"
(
    warehouse_id BIGINT NOT NULL,
    sku          INT    NOT NULL,
    count        INT    NOT NULL,
    reserved     INT    NOT NULL DEFAULT 0,
    PRIMARY KEY (warehouse_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "stocks";
-- +goose StatementEnd
