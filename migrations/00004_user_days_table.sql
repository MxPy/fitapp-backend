-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_days
(
    user_id UUID REFERENCES users(id),
    user_date  DATE NOT NULL,
    daily_kcal INTEGER,
    daily_proteins INTEGER,
    daily_carbs INTEGER,
    daily_fats INTEGER,
    PRIMARY KEY (id)  -- Explicitly define primary key
) INHERITS (base_entity);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_days;
-- +goose StatementEnd
