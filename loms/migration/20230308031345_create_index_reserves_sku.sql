-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS reserves_sku_idx ON reserves (sku);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS reserves_sku_idx;
-- +goose StatementEnd
