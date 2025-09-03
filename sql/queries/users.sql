-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, api_key)
-- Update schema, need to be creating new API keys for new users
-- SQL will handle creating of new API keys every time a new user is created
-- Not using $5, function signature won't change
VALUES ($1, $2, $3, $4, encode(sha256(random()::text::bytea), 'hex'))
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;