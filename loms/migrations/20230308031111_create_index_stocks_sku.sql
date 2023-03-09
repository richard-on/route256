-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS stocks_sku_idx ON stocks (sku);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS stocks_sku_idx;
-- +goose StatementEnd
