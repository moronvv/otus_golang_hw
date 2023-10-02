-- +goose Up
CREATE TABLE events (
    id BIGINT PRIMARY KEY generated always AS IDENTITY,
    title TEXT NOT NULL,
    description TEXT,
    datetime TIMESTAMP NOT NULL,
    duration BIGINT NOT NULL,
    user_id UUID NOT NULL,
    notify_before BIGINT NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE events;
