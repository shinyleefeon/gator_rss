-- +goose Up
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMP,
    feed_id TEXT NOT NULL REFERENCES feeds(name) ON DELETE CASCADE,
    UNIQUE (url)
);

-- +goose Down

DROP TABLE posts;