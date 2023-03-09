-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS order_items_order_id_idx ON order_items (order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS order_items_order_id_uidx;
-- +goose StatementEnd
