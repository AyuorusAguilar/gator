-- name: CreateFeedFollow :one
WITH inserted_row AS (
    INSERT INTO feed_follows
        (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    ) RETURNING *)
SELECT inserted_row.*,
    users.name AS user_name,
    feeds.name AS feed_name,
    feeds.url AS feed_url
FROM inserted_row
INNER JOIN users ON users.id = inserted_row.user_id
INNER JOIN feeds on feeds.id = inserted_row.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT feeds.name, feeds.url FROM feed_follows
 INNER JOIN feeds ON feeds.id = feed_follows.feed_id
 WHERE feed_follows.user_id = $1;

-- name: DeleteFollow :exec
DELETE FROM feed_follows
WHERE feed_id = $1 AND user_id = $2;
