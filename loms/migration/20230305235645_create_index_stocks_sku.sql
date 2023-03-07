-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS "sku_idx" ON stocks (sku);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "sku_idx";
-- +goose StatementEnd
