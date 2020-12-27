-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS event (
    event_id serial PRIMARY KEY,
    title TEXT NOT NULL,
    start timestamptz NOT NULL,
    stop timestamptz NOT NULL,
    description TEXT,
    user_id int NOT NULL,
    notification bigint
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE event;
