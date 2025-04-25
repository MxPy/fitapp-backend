-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS base_entity
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS users
(
    username  TEXT NOT NULL,
    full_name TEXT NOT NULL
) INHERITS (base_entity);

CREATE TABLE IF NOT EXISTS products
(
    name  TEXT NOT NULL,
    kcal INTEGER NOT NULL
    proteins INTEGER NOT NULL
    carbs INTEGER NOT NULL
    fats INTEGER NOT NULL
) INHERITS (base_entity);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS base_entity;
-- +goose StatementEnd
