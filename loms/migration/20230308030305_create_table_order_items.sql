-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_items
(
    order_id BIGINT NOT NULL REFERENCES orders (order_id),
    sku      BIGINT NOT NULL,
    count    INT    NOT NULL,
    PRIMARY KEY (order_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_items;
-- +goose StatementEnd
