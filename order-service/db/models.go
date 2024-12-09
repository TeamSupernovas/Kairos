// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"
)

type Order struct {
	OrderID     string      `json:"orderId"`
	UserID      string      `json:"userId"`
	ChefID      string      `json:"chefId"`
	TotalPrice  float64     `json:"totalPrice"`
	PickupTime  **time.Time `json:"pickupTime"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	CanceledAt  **time.Time `json:"canceledAt"`
	CompletedAt **time.Time `json:"completedAt"`
	OrderStatus string      `json:"orderStatus"`
	DeletedAt   **time.Time `json:"deletedAt"`
}

type OrderItem struct {
	OrderItemID     string      `json:"orderItemId"`
	OrderID         string      `json:"orderId"`
	DishID          string      `json:"dishId"`
	DishOrderStatus string      `json:"dishOrderStatus"`
	Quantity        int32       `json:"quantity"`
	PricePerUnit    float64     `json:"pricePerUnit"`
	CreatedAt       time.Time   `json:"createdAt"`
	DeletedAt       **time.Time `json:"deletedAt"`
}
