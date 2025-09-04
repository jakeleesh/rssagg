-- +goose Up
-- Relationship between a User and all the Feeds they follow
-- Many to many 
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    -- User is deleted, delete all data about what feeds they're following
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    -- Feed gets deleted, delete all data related to that Feed
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    -- UNIQUE make sure never have 2 instances of a follow for same User Feed relationship
    -- As a User, can only follow feed once
    UNIQUE(user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
