-- name: CreateRating :one
INSERT INTO ratings (dish_id, user_id, rating, review_text)
VALUES ($1, $2, $3, $4)
RETURNING id, dish_id, user_id, rating, review_text, created_at;

-- name: GetRating :one
SELECT id, dish_id, user_id, rating, review_text, created_at
FROM ratings
WHERE id = $1;

-- name: ListRatings :many
SELECT id, dish_id, user_id, rating, review_text, created_at
FROM ratings
ORDER BY created_at DESC;

-- name: UpdateRating :one
UPDATE ratings
SET rating = $2, review_text = $3
WHERE id = $1
RETURNING id, dish_id, user_id, rating, review_text, created_at;

-- name: DeleteRating :exec
DELETE FROM ratings
WHERE id = $1;
