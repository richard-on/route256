-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart
(
    user_id    BIGINT PRIMARY KEY,
    updated_at timestamptz NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart;
-- +goose StatementEnd
