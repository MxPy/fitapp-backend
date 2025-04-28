-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS base_entity
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS base_entity;
-- +goose StatementEnd
