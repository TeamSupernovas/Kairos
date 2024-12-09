-- name: GetUserOrders :many
SELECT * 
FROM orders
WHERE user_id = $1 
  AND deleted_at IS NULL;

-- name: GetChefOrders :many
SELECT * 
FROM orders
WHERE chef_id = $1 
  AND deleted_at IS NULL;

-- name: UpdateOrderItemStatus :exec
UPDATE order_items
SET dish_order_status = $2
WHERE order_item_id = $1 AND deleted_at IS NULL;

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
    order_id, user_id, chef_id, total_price, pickup_time, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
);

-- name: AddOrderItem :exec
INSERT INTO order_items (
    order_item_id, order_id, dish_id, dish_order_status, quantity, price_per_unit, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP
);

-- name: GetOrderItemsByOrderID :many
SELECT order_item_id, order_id, dish_id, dish_order_status, quantity, price_per_unit, created_at, deleted_at
FROM order_items
WHERE order_id = $1;

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
