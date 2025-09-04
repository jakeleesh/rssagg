-- name: CreatePost :one
INSERT INTO posts (
    id, created_at, updated_at, title, description, published_at, url, feed_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* FROM posts
-- Know what feed every post in the database belongs to
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
-- Tells us which feeds User is following
WHERE feed_follows.user_id = $1
-- Newest stuff first
ORDER BY posts.published_at DESC
LIMIT $2;