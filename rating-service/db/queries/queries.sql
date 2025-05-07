-- name: CreateRating :one
INSERT INTO ratings (
    dish_id, dish_name,
    chef_id, chef_name,
    user_id, user_name,
    rating, review_text
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, dish_id, dish_name, chef_id, chef_name, user_id, user_name, rating, review_text, created_at;

-- name: GetRating :one
SELECT id, dish_id, dish_name, chef_id, chef_name, user_id, user_name, rating, review_text, created_at
FROM ratings
WHERE id = $1;

-- name: ListRatings :many
SELECT id, dish_id, dish_name, chef_id, chef_name, user_id, user_name, rating, review_text, created_at
FROM ratings
ORDER BY created_at DESC;

-- name: UpdateRating :one
UPDATE ratings
SET rating = $2, review_text = $3
WHERE id = $1
RETURNING id, dish_id, dish_name, chef_id, chef_name, user_id, user_name, rating, review_text, created_at;

-- name: DeleteRating :exec
DELETE FROM ratings
WHERE id = $1;

-- name: ListRatingsByDish :many
SELECT id, dish_id, dish_name, chef_id, chef_name, user_id, user_name, rating, review_text, created_at
FROM ratings
WHERE dish_id = $1
ORDER BY created_at DESC;

-- name: ListRatingsByUser :many
SELECT id, dish_id, dish_name, chef_id, chef_name, user_id, user_name, rating, review_text, created_at
FROM ratings
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListRatingsByChef :many
SELECT id, dish_id, dish_name, chef_id, chef_name, user_id, user_name, rating, review_text, created_at
FROM ratings
WHERE chef_id = $1
ORDER BY created_at DESC;

-- name: GetChefAverageRating :one
SELECT chef_id, AVG(rating)::FLOAT AS average_rating
FROM ratings
WHERE chef_id = $1
GROUP BY chef_id;

-- name: GetOrderAverageRating :one
SELECT dish_id, AVG(rating)::FLOAT AS average_rating
FROM ratings
WHERE dish_id = $1
GROUP BY dish_id;
