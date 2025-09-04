-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
-- User to get all of the feeds
-- Not an authenticated endpoint
SELECT * FROM feeds;

-- name: GetNextFeedsToFetch :many
-- Purpose is to get feed that next needs to be fetched
-- Find any feeds that have never been fetched before, priority
-- If every feed fetched, find fetched longest ago
-- Many goroutines to fetch different feeds
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
-- Pass in how many feeds we want
LIMIT $1;

-- name: MarkFeedAsFetched :one
UPDATE feeds
-- For auditing purposes
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING *;