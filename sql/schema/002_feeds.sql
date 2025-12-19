-- +goose Up
CREATE TABLE feeds (
    name TEXT PRIMARY KEY,
    url TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose DOWN
DROP TABLE feeds;