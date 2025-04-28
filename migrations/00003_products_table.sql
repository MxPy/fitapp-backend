-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products
(
    product_name  TEXT NOT NULL,
    kcal INTEGER NOT NULL,
    proteins INTEGER NOT NULL,
    carbs INTEGER NOT NULL,
    fats INTEGER NOT NULL,
    PRIMARY KEY (id)  -- Explicitly define primary key
) INHERITS (base_entity);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
