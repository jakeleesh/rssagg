-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollows :many
SELECT * FROM feed_follows WHERE user_id = $1;

-- name: DeleteFeedFollow :exec
-- Not returning record, just run a SQL query
-- Don't actually need user_id, id already unique
-- Tacking on user_id is prevent someone who doesn't own FeedFollow to unfollow
-- Ensures only user who owns follow record can unfollow
DELETE FROM feed_follows WHERE id = $1 AND user_id = $2;
