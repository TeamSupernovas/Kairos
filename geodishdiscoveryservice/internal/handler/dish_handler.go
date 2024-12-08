package handler

import (
	"geodishdiscoveryservice/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DishHandler struct {
	dishService *service.DishService
}

func NewDishHandler(dishService *service.DishService) *DishHandler {
	return &DishHandler{
		dishService: dishService,
	}
}

func (dishHandler *DishHandler) GetNearbyDishes(context *gin.Context) {
	// Extract and validate required query parameters
	latitudeStr := context.Query("latitude")
	longitudeStr := context.Query("longitude")

	if latitudeStr == "" || longitudeStr == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "latitude and longitude are required"})
		return
	}

	// Parse latitude and longitude to float64
	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude format"})
		return
	}

	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude format"})
		return
	}

	// Parse optional parameters
	radius := 5000.0 // default radius 5km
	if radiusStr := context.Query("radius"); radiusStr != "" {
		radius, err = strconv.ParseFloat(radiusStr, 64)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radius format"})
			return
		}
	}

	mealCourse := context.Query("mealCourse")
	dietaryCategory := context.Query("dietaryCategory")

	maxPrice := -1.0
	if maxPriceStr := context.Query("maxPrice"); maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maxPrice format"})
			return
		}
	}

	page := 1
	if pageStr := context.Query("page"); pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
			return
		}
	}

	pageSize := 10
	if pageSizeStr := context.Query("pageSize"); pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
			return
		}
	}

	// Call the service layer to fetch nearby dishes
	dishes, total, err := dishHandler.dishService.GetNearbyDishes(
		context,
		latitude, longitude, radius, mealCourse, dietaryCategory, maxPrice, page, pageSize,
	)
	if err != nil {
		log.Printf("Error fetching nearby dishes: %+v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Return paginated response
	context.JSON(http.StatusOK, gin.H{
		"totalCount": total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": (total + pageSize - 1) / pageSize,
		"data":       dishes,
	})
}
