package dto

import "time"

type Address struct {
	Street		string		`json:"street" binding:"required"`
	City		string		`json:"city" binding:"required"`
	State		string		`json:"state" binding:"required"`
	PostalCode	string		`json:"postalCode" binding:"required"`
	Country		string		`json:"country" binding:"required"`
}

type GeoPoint struct {
	Latitude	float64		`json:"latitude"`
	Longitude	float64		`json:"longitude"`
}

type CreateDishRequest struct {
	DishName			string		`json:"dishName" binding:"required"`
	ChefID				string		`json:"chefId" binding:"required"`
	Description			string		`json:"description" binding:"required"`
	Price				float64		`json:"price" binding:"required,gt=0"`
	AvailablePortions	int			`json:"availablePortions" binding:"required,gt=0"`
	MealCourse			string		`json:"mealCourse"`
	DietaryCategory		string		`json:"dietaryCategory"`
	AvailableUntil		time.Time	`json:"availableUntil"`
	Address				Address		`json:"address" binding:"required"`
}


type UpdateDishRequest struct {
	DishID			  string	`json:"dishId" binding:"required"`
	DishName          string    `json:"dishName" binding:"required"`
	ChefID            string    `json:"chefId" binding:"required"`
	Description       string    `json:"description"`
	Price             float64   `json:"price" binding:"gt=0"`
	AvailablePortions int       `json:"availablePortions" binding:"gt=0"`
	MealCourse        string    `json:"mealCourse"`
	DietaryCategory   string    `json:"dietaryCategory"`
	AvailableUntil    time.Time `json:"availableUntil"`
	Address           Address   `json:"address"`
}

type DishResponse struct {
	DishID				string		`json:"dishId"`
	DishName			string		`json:"dishName"`
	ChefID				string		`json:"chefId" binding:"required"`
	Description			string		`json:"description"`
	Price				float64		`json:"price" binding:"gt=0"`
	AvailablePortions	int			`json:"availablePortions" binding:"gt=0"`
	MealCourse			string		`json:"mealCourse"`
	DietaryCategory		string		`json:"dietaryCategory"`
	AvailableUntil		time.Time	`json:"availableUntil"`
	Address				Address		`json:"address"`
	GeoPoint			GeoPoint	`json:"geoPoint"`
	CreatedAt			time.Time	`json:"createdAt"`
	UpdatedAt			time.Time	`json:"updatedAt"`
}

