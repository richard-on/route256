-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reserves
(
    order_id     BIGINT NOT NULL REFERENCES orders (order_id),
    sku          BIGINT NOT NULL,
    warehouse_id BIGINT NOT NULL,
    count        INT    NOT NULL,
    PRIMARY KEY (order_id, sku, warehouse_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reserves;
-- +goose StatementEnd
