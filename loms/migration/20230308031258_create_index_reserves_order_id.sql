-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS reserves_order_id_idx ON reserves (order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS reserves_order_id_idx;
-- +goose StatementEnd
