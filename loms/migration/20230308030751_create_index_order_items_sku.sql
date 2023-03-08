-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS order_items_sku_idx ON order_items (sku);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS order_items_sku_idx;
-- +goose StatementEnd
