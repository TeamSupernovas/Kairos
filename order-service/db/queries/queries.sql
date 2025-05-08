-- name: GetUserOrders :many
SELECT * 
FROM orders
WHERE user_id = $1 
  AND deleted_at IS NULL
ORDER BY updated_at DESC, created_at DESC;


-- name: GetChefOrders :many
SELECT * 
FROM orders
WHERE chef_id = $1 
  AND deleted_at IS NULL
ORDER BY updated_at DESC, created_at DESC;


-- name: UpdateOrderItemStatus :exec
UPDATE order_items
SET dish_order_status = $2
WHERE order_item_id = $1 AND deleted_at IS NULL;

-- name: TouchOrderUpdatedAt :exec
UPDATE orders
SET updated_at = CURRENT_TIMESTAMP
WHERE order_id = $1
  AND deleted_at IS NULL;

-- name: SoftDeleteItem :exec
UPDATE order_items
SET deleted_at = CURRENT_TIMESTAMP
WHERE order_item_id = $1;

-- name: HardDeleteItem :exec
DELETE FROM order_items
WHERE order_item_id = $1;

-- name: SoftDeleteOrder :exec
UPDATE orders
SET deleted_at = CURRENT_TIMESTAMP
WHERE order_id = $1;

-- name: HardDeleteOrder :exec
DELETE FROM orders
WHERE order_id = $1;

-- name: CreateOrder :exec
INSERT INTO orders (
    order_id, user_id, chef_id, chef_name, user_name, total_price, pickup_time, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
);


-- name: AddOrderItem :exec
INSERT INTO order_items (
    order_item_id, order_id, dish_id, dish_name, dish_order_status, quantity, price_per_unit, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP
);

-- name: GetOrderItemsByOrderID :many
SELECT order_item_id, order_id, dish_id, dish_name, dish_order_status, quantity, price_per_unit, created_at, deleted_at
FROM order_items
WHERE order_id = $1
ORDER BY created_at DESC;

-- name: DeleteOrderItem :one
UPDATE order_items 
SET deleted_at = CURRENT_TIMESTAMP 
WHERE order_item_id = $1 
AND deleted_at IS NULL
RETURNING order_item_id;


-- name: GetOrderIDByOrderItem :one
SELECT order_id 
FROM order_items 
WHERE order_item_id = $1;

-- name: CountActiveOrderItems :one
SELECT COUNT(*) 
FROM order_items 
WHERE order_id = $1 
AND deleted_at IS NULL;

-- name: DeleteOrder :exec
UPDATE orders 
SET deleted_at = CURRENT_TIMESTAMP 
WHERE order_id = $1 
AND deleted_at IS NULL;

-- name: GetOrderItemStatus :one
SELECT dish_order_status 
FROM order_items 
WHERE order_item_id = $1 AND deleted_at IS NULL;

-- name: UpdateOrderItemStatusByOrderIDDishID :execrows
UPDATE order_items
SET dish_order_status = $3
WHERE order_id = $1
  AND dish_id = $2
  AND deleted_at IS NULL;

-- name: GetUserAndChefNameByOrderID :one
SELECT user_name, chef_name
FROM orders
WHERE order_id = $1
  AND deleted_at IS NULL;

-- name: GetDishNameByOrderIDAndDishID :one
SELECT dish_name
FROM order_items
WHERE order_id = $1
  AND dish_id = $2
  AND deleted_at IS NULL;