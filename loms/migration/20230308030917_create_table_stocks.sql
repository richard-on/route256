-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stocks
(
    sku          BIGINT NOT NULL,
    warehouse_id BIGINT NOT NULL,
    count        INT    NOT NULL,
    PRIMARY KEY (sku, warehouse_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stocks;
-- +goose StatementEnd
