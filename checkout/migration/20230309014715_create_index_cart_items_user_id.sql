-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS cart_items_user_id_idx ON cart_items (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS cart_items_user_id_idx;
-- +goose StatementEnd
