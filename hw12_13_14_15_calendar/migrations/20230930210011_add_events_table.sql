-- +goose Up
CREATE TABLE events (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    datetime TIMESTAMP NOT NULL,
    duration INT NOT NULL,
    user_id UUID NOT NULL,
    notify_before INT NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE events;
