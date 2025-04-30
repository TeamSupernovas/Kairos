package controllers

import (
	"net/http"
	"strconv"

	"kairos/rating-service/db"
	"kairos/rating-service/kafka"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type RatingRequest struct {
	DishID     string `json:"dishId" binding:"required"`
	ChefID     string `json:"chefId" binding:"required"`
	UserID     string `json:"userId" binding:"required"`
	Rating     int32  `json:"rating" binding:"required,min=1,max=5"`
	ReviewText string `json:"reviewText"`
}

// CreateRating handles the creation of a new rating
func CreateRating(c *gin.Context, q *db.Queries) {
	var req RatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := db.CreateRatingParams{
		DishID: req.DishID,
		ChefID: req.ChefID,
		UserID: req.UserID,
		Rating: pgtype.Int4{
			Int32: req.Rating,
			Valid: true,
		},
		ReviewText: pgtype.Text{
			String: req.ReviewText,
			Valid:  true,
		},
	}

	rating, err := q.CreateRating(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msg := "New rating submitted"
	notifType := "RatingService"
	_ = kafka.SendNotification(c.Request.Context(), gin.H{"user_id": req.UserID, "message": msg, "type": notifType})
	_ = kafka.SendNotification(c.Request.Context(), gin.H{"user_id": req.ChefID, "message": msg, "type": notifType})

	c.JSON(http.StatusCreated, rating)
}

// GetRating retrieves a rating by ID
func GetRating(c *gin.Context, q *db.Queries) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	rating, err := q.GetRating(c.Request.Context(), int32(id))
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "rating not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rating)
}

// ListRatings retrieves all ratings
func ListRatings(c *gin.Context, q *db.Queries) {
	ratings, err := q.ListRatings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ratings)
}

// UpdateRating updates an existing rating
func UpdateRating(c *gin.Context, q *db.Queries) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var req RatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := db.UpdateRatingParams{
		ID: int32(id),
		Rating: pgtype.Int4{
			Int32: req.Rating,
			Valid: true,
		},
		ReviewText: pgtype.Text{
			String: req.ReviewText,
			Valid:  true,
		},
	}

	rating, err := q.UpdateRating(c.Request.Context(), params)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "rating not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msg := "Rating updated"
	notifType := "RatingService"
	_ = kafka.SendNotification(c.Request.Context(), gin.H{"user_id": rating.UserID, "message": msg, "type": notifType})
	_ = kafka.SendNotification(c.Request.Context(), gin.H{"user_id": rating.ChefID, "message": msg, "type": notifType})

	c.JSON(http.StatusOK, rating)
}

// DeleteRating deletes a rating by ID
func DeleteRating(c *gin.Context, q *db.Queries) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	rating, err := q.GetRating(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rating not found"})
		return
	}

	err = q.DeleteRating(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msg := "Rating deleted"
	notifType := "RatingService"
	_ = kafka.SendNotification(c.Request.Context(), gin.H{"user_id": rating.UserID, "message": msg, "type": notifType})
	_ = kafka.SendNotification(c.Request.Context(), gin.H{"user_id": rating.ChefID, "message": msg, "type": notifType})

	c.JSON(http.StatusOK, gin.H{"message": "rating deleted successfully"})
}

// ListRatingsForDish returns all reviews for a given dish
func ListRatingsForDish(c *gin.Context, q *db.Queries) {
	dishID := c.Param("dishId")

	ratings, err := q.ListRatingsByDish(c.Request.Context(), dishID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ratings)
}

func ListRatingsByUser(c *gin.Context, q *db.Queries) {
	userID := c.Param("userId")

	ratings, err := q.ListRatingsByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ratings)
}

func ListRatingsByChef(c *gin.Context, q *db.Queries) {
	chefID := c.Param("chefId")

	ratings, err := q.ListRatingsByChef(c.Request.Context(), chefID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ratings)
}

func GetChefAverageRating(c *gin.Context, q *db.Queries) {
	chefID := c.Param("chefId")

	avg, err := q.GetChefAverageRating(c.Request.Context(), chefID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "no ratings found for chef"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, avg)
}

func GetDishAverageRating(c *gin.Context, q *db.Queries) {
	dishID := c.Param("dishId")

	avg, err := q.GetOrderAverageRating(c.Request.Context(), dishID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "no ratings found for dish"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, avg)
}
