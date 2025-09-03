-- +goose Up

CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    -- Means if try to create a feed for a user_id that does not exist in users table, get an error
    -- Don't want feeds to exist without a user who created them
    -- ON DELETE CASCADE: When a user is deleted, all feeds associated with that user to be automatially deleted
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;