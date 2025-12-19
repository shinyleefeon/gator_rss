-- +goose Up
CREATE TABLE feed_follows (
    ID SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id TEXT NOT NULL REFERENCES feeds(name) ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);



-- +goose Down
DROP TABLE feed_follows;