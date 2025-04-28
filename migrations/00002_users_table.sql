-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    user_username  TEXT NOT NULL,
    user_full_name TEXT NOT NULL,
    user_sex BOOLEAN NOT NULL,
    user_height INTEGER NOT NULL,
    user_weight INTEGER NOT NULL,
    user_age INTEGER NOT NULL,
    PRIMARY KEY (id)  -- Explicitly define primary key
) INHERITS (base_entity);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
