-- +goose Up
-- Difference between TEXT, VARCHAR exactly 64 characters. Want API keys to be 64 characters long
ALTER TABLE users ADD COLUMN api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (
    -- Set DEFAULT because already have users in database
    -- If add column UNIQUE NOT NULL, need to provide unique default for every existing record
    -- Use random number generation
    -- Basically generating random bytes and casting into byte array, using sha256 to get fixed size output
    -- Saying take a big random slice of bytes, hash so we get fixed size and encode it in hexadecimal
    -- Get 64 unique hexadecimal characters
    encode(sha256(random()::text::bytea), 'hex')
);

-- +goose Down
ALTER TABLE users DROP COLUMN api_key;