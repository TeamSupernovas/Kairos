package controllers

import (
	"net/http"
	"strconv"

	"kairos/rating-service/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type RatingRequest struct {
	DishID     int32  `json:"dishId" binding:"required"`
	UserID     int32  `json:"userId" binding:"required"`
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

	c.JSON(http.StatusCreated, rating)
}

// GetRating retrieves a rating by ID
func GetRating(c *gin.Context, q *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
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

	c.JSON(http.StatusOK, rating)
}

// DeleteRating deletes a rating by ID
func DeleteRating(c *gin.Context, q *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	err = q.DeleteRating(c.Request.Context(), int32(id))
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "rating not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rating deleted successfully"})
}