-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart_items
(
    user_id BIGINT NOT NULL,
    sku     BIGINT NOT NULL,
    count   INT    NOT NULL,
    PRIMARY KEY (user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_items;
-- +goose StatementEnd
