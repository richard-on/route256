-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS stocks_warehouse_id_idx ON stocks (warehouse_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS stocks_warehouse_id_idx;
-- +goose StatementEnd
