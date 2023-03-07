-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS items_gin_idx ON "order" USING gin (items);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS items_gin_idx;
-- +goose StatementEnd
