-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS reserves_warehouse_id_idx ON reserves (warehouse_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS reserves_warehouse_id_idx;
-- +goose StatementEnd
