-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;


-- name: MarkFetched :exec
UPDATE feeds
SET last_fetched_at = $1
WHERE id = $2;

-- name: GetOldestFetched :one
SELECT * FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;

-- name: GetOldestFetchedByUserId :one
SELECT * FROM feeds
WHERE user_id = $1
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;