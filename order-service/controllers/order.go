package controllers

import (
	"fmt"
	"kairos/order-service/db"
	"kairos/order-service/kafka"
	"github.com/IBM/sarama"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderResponse struct {
	Order      db.Order       `json:"order"`
	OrderItems []db.OrderItem `json:"order_items"`
}

// CreateOrder handles the creation of an order with multiple order items
func CreateOrder(queries *db.Queries, kafkaProducer sarama.SyncProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newOrder struct {
			UserID     string      `json:"user_id" binding:"required"`
			ChefID     string      `json:"chef_id" binding:"required"`
			TotalPrice float64     `json:"total_price"`
			PickupTime **time.Time `json:"pickup_time"`
			OrderItems []struct {
				DishID          string  `json:"dish_id" binding:"required"`
				DishOrderStatus string  `json:"dish_order_status"`
				Quantity        int32   `json:"quantity" binding:"required"`
				PricePerUnit    float64 `json:"price_per_unit" binding:"required"`
			} `json:"order_items" binding:"required"`
		}

		// Bind incoming JSON data to struct
		if err := c.ShouldBindJSON(&newOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate a unique order ID
		orderID := uuid.New().String()

		// Call the CreateOrder query
		err := queries.CreateOrder(c, db.CreateOrderParams{
			OrderID:    orderID,
			UserID:     newOrder.UserID,
			ChefID:     newOrder.ChefID,
			TotalPrice: newOrder.TotalPrice,
			PickupTime: newOrder.PickupTime,
		})

		if err != nil {
			log.Println("Error creating order:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}

		// Insert the order items into the order_items table
		for _, item := range newOrder.OrderItems {
			orderItemID := uuid.New().String()

			err := queries.AddOrderItem(c, db.AddOrderItemParams{
				OrderItemID:     orderItemID,
				OrderID:         orderID,
				DishID:          item.DishID,
				DishOrderStatus: item.DishOrderStatus,
				Quantity:        item.Quantity,
				PricePerUnit:    item.PricePerUnit,
			})

			if err != nil {
				log.Println("Error adding order item:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add order items"})
				return
			}
	// Publish order placed event to Kafka
	orderPlacedEvent := kafka.OrderPlacedEvent{
	OrderID:    orderID,
	DishID:     item.DishID,
	Portions:   item.Quantity,
	}

	err = kafka.PublishOrderPlaced(kafkaProducer, orderPlacedEvent)
	if err != nil {
		log.Println("Error publishing to Kafka:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish order placed event"})
		return
	}


		}

		// Respond with success message
		c.JSON(http.StatusOK, gin.H{
			"message":  "Order created successfully",
			"order_id": orderID,
		})
	}
}

func GetOrdersByConsumer(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user_id from query parameters
		userID := c.DefaultQuery("user_id", "")

		// Validate that the user_id is not empty
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		// Fetch orders for the consumer using db.Queries
		orders, err := queries.GetUserOrders(c, userID)
		if err != nil {
			// Log the error and return a response with a failure message
			fmt.Println("Error fetching orders:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
			return
		}

		// Create a slice to hold our response objects
		orderResponses := make([]OrderResponse, len(orders))

		// For each order, fetch its corresponding order items
		for i, order := range orders {
			orderItems, err := queries.GetOrderItemsByOrderID(c, order.OrderID)
			if err != nil {
				// Handle error while fetching order items
				fmt.Println("Error fetching order items:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
				return
			}

			// Create the response object combining order and its items
			orderResponses[i] = OrderResponse{
				Order:      order,
				OrderItems: orderItems,
			}
		}

		// Return the combined response as JSON
		c.JSON(http.StatusOK, gin.H{"orders": orderResponses})
	}
}

func GetOrdersByProvider(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract chef_id from query parameters
		chefID := c.DefaultQuery("chef_id", "")
		fmt.Println(chefID)

		// Validate that the chef_id is not empty
		if chefID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "chef_id is required"})
			return
		}

		// Fetch all orders for the chef using db.Queries
		orders, err := queries.GetChefOrders(c, chefID)
		if err != nil {
			fmt.Println("Error fetching chef orders:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
			return
		}

		// Create a slice to hold our response objects
		orderResponses := make([]OrderResponse, len(orders))

		// For each order, fetch its corresponding order items
		for i, order := range orders {
			orderItems, err := queries.GetOrderItemsByOrderID(c, order.OrderID)
			if err != nil {
				fmt.Println("Error fetching order items:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
				return
			}

			// Create the response object combining order and its items
			orderResponses[i] = OrderResponse{
				Order:      order,
				OrderItems: orderItems,
			}
		}

		// Return the combined response as JSON
		c.JSON(http.StatusOK, gin.H{
			"orders": orderResponses,
		})
	}
}

// DeleteOrderItem handles the deletion of an order item and potentially the parent order
func DeleteOrderItem(pool *pgxpool.Pool, queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		orderItemID := c.Param("orderItemId")

		// Start a transaction
		tx, err := pool.Begin(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}
		txQueries := db.New(tx)

		defer func() {
			if p := recover(); p != nil {
				tx.Rollback(ctx)
				panic(p) // Re-throw panic after rollback
			} else if err != nil {
				tx.Rollback(ctx)
			} else {
				err = tx.Commit(ctx)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
				}
			}
		}()

		// Get the order ID associated with this item
		orderID, err := txQueries.GetOrderIDByOrderItem(ctx, orderItemID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order ID"})
			return
		}

		// Delete the order item
		_, err = txQueries.DeleteOrderItem(ctx, orderItemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order item"})
			return
		}

		// Check if this was the last active item in the order
		count, err := txQueries.CountActiveOrderItems(ctx, orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count active items"})
			return
		}

		// If no active items remain, delete the parent order
		if count == 0 {
			err = txQueries.DeleteOrder(ctx, orderID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete parent order"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "Order item deleted successfully",
			"orderDeleted": count == 0,
		})
	}
}

// UpdateOrderItemStatus handles updating the status of an order item.
func UpdateOrderItemStatus(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get the order item ID from the route parameters.
		orderItemID := c.Param("orderId")
		if orderItemID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order item ID is required"})
			return
		}

		// Get the new status from the request body.
		var requestBody struct {
			Status string `json:"status" binding:"required"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Ensure the status is valid.
		validStatuses := map[string]bool{
			"pending":   true,
			"confirmed": true,
			"ready":     true,
			"canceled":  true,
			"completed": true,
		}
		if !validStatuses[requestBody.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
			return
		}

		// Update the order item status in the database.
		err := queries.UpdateOrderItemStatus(ctx, db.UpdateOrderItemStatusParams{
			OrderItemID:     orderItemID,
			DishOrderStatus: requestBody.Status,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order item status"})
			return
		}

		// Respond with a success message.
		c.JSON(http.StatusOK, gin.H{
			"message": "Order item status updated successfully",
		})
	}
}

func GetOrderItemStatus(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		orderItemID := c.Param("orderId")

		// Fetch the order item status from the database
		status, err := queries.GetOrderItemStatus(ctx, orderItemID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order item status"})
			return
		}

		// Return the status as a JSON response
		c.JSON(http.StatusOK, gin.H{
			"orderItemId":     orderItemID,
			"dishOrderStatus": status,
		})
	}
}
