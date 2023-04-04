-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS outbox
(
    id      BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    key     TEXT,
    payload BYTEA    NOT NULL,
    status  SMALLINT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox;
-- +goose StatementEnd
