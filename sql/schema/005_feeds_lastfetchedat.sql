-- +goose Up
-- Keep track of when we last fetched posts for a given feed
-- Nullable
ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE feeds DROP COLUMN last_fetched_at;