package handler

import (
	"dishmanagementservice/internal/dto"
	"dishmanagementservice/internal/mapper"
	"dishmanagementservice/internal/model"
	"dishmanagementservice/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DishHandler struct {
	dishService *service.DishService
}

func NewDishHandler(dishService *service.DishService) *DishHandler {
	return &DishHandler{
		dishService: dishService,
	}
}

func (dishHandler *DishHandler) AddDish(context *gin.Context) {
	var req dto.CreateDishRequest

	err := context.BindJSON(&req)

	if (err != nil) {
		log.Printf("Invalid input error: %+v", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	// map the request to the model
	newDish := mapper.ToDishModel(req)

	err = dishHandler.dishService.AddDish(newDish)
	if (err != nil) {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Dish added successfully"})
}

func (dishHandler *DishHandler) UpdateDish(context *gin.Context) {
	// Extract the dishID from the URL parameters
	dishIDParam := context.Param("dishId")
	if dishIDParam == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Dish ID is required"})
		return
	}

	// Validate if the dishIDParam is a valid UUID
	_, err := uuid.Parse(dishIDParam)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Dish ID format"})
		return
	}

	// Bind the request JSON to a dish model
	var req dto.UpdateDishRequest

	if err := context.BindJSON(&req); err != nil {
		log.Printf("Invalid input error: %+v", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var updatedDish model.Dish
	updatedDish.DishID = req.DishID
	updatedDish.DishName = req.DishName
	updatedDish.ChefID = req.ChefID
	updatedDish.Description = req.Description
	updatedDish.Price = req.Price
	updatedDish.AvailablePortions = req.AvailablePortions
	updatedDish.MealCourse = req.MealCourse
	updatedDish.DietaryCategory = req.DietaryCategory
	updatedDish.AvailableUntil = req.AvailableUntil
	updatedDish.Address = model.Address(req.Address)
	
	// Call the service layer to update the dish
	err = dishHandler.dishService.UpdateDish(dishIDParam, &updatedDish)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	context.JSON(http.StatusOK, gin.H{"message": "Dish updated successfully"})
}

func (dishHandler *DishHandler) GetDishByID(context *gin.Context) {
	dishIDParam := context.Param("dishId")

	if (dishIDParam == "") {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Dish ID is required"})
		return
	}

	// Validate if the dishIDParam is a valid UUID
	_, err := uuid.Parse(dishIDParam)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Dish ID format"})
		return
	}

	dish, err := dishHandler.dishService.GetDishByID(dishIDParam)
	if err != nil {
		if err.Error() == "no dish found with the provided ID" {
			context.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	// Return the dish as JSON response
	context.JSON(http.StatusOK, dish)
}

func (dishHandler *DishHandler) DeleteDishById(context *gin.Context) {
	dishIDParam := context.Param("dishId")

	if (dishIDParam == "") {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Dish ID is required"})
		return
	}

	// Validate if the dishIDParam is a valid UUID
	_, err := uuid.Parse(dishIDParam)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Dish ID format"})
		return
	}

	// Call the DeleteDishByID method from the repository
	err = dishHandler.dishService.DeleteDishByID(dishIDParam)
	if err != nil {
		log.Printf("Error deleting dish: %+v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete dish"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Dish deleted successfully"})

}
