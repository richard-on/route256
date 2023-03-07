-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "order"
(
    order_id   BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id    BIGINT   NOT NULL,
    items      item[]   NOT NULL,
    status     SMALLINT NOT NULL DEFAULT 0,
    created_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "order";
-- +goose StatementEnd
