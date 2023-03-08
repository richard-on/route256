-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart_items
(
    user_id BIGINT PRIMARY KEY REFERENCES cart (user_id),
    sku     BIGINT NOT NULL,
    count   INT    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_items;
-- +goose StatementEnd
