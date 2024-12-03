package model

import (
	"time"
)

type Dish struct {
	DishID             string
	DishName           string
	ChefID             string
	Description        string
	Price              float64
	AvailablePortions  int
	MealCourse		   string
	DietaryCategory    string
	AvailableUntil 	   time.Time
	Address			   Address
	GeoPoint 		   GeoPoint
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt         *time.Time
}

type Address struct {
	Street           string
	City             string
	State            string
	PostalCode       string
	Country          string
}

type GeoPoint struct {
	Latitude         float64
	Longitude        float64
}

func NewDish(dishName string, chefID string, price float64, availablePortions int) *Dish {
	dish := &Dish{
		DishName: dishName,
		ChefID: chefID,
		Price: price,
		AvailablePortions: availablePortions,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return dish;
}

func (dish *Dish)UpdateAvailablePortions(availablePortions int) {
	dish.AvailablePortions = availablePortions
	dish.UpdatedAt = time.Now()
}

func (dish *Dish) UpdatePrice(newPrice float64) {
	dish.Price = newPrice
	dish.UpdatedAt = time.Now()
}

func (dish *Dish) MarkAsDeleted() {
	currentTime := time.Now()
	dish.DeletedAt = &currentTime
	dish.UpdatedAt = currentTime
}

func (dish *Dish) SetGeoPoint(latitude, longitude float64) {
	dish.GeoPoint = GeoPoint{
		Latitude: latitude,
		Longitude: longitude,
	}
	dish.UpdatedAt = time.Now()
}

func (dish *Dish) SetAddress(street, city, state, postalCode, country string) {
	dish.Address = Address{
		Street: street,
		City: city,
		State: state,
		PostalCode: postalCode,
		Country: country,
	}
	dish.UpdatedAt = time.Now()
}

func (dish *Dish) SetAvailableUntil(availableUntil time.Time) {
	dish.AvailableUntil = availableUntil
	dish.UpdatedAt = time.Now()
}

func (dish *Dish) SetDescription(description string) {
	dish.Description = description
	dish.UpdatedAt = time.Now()
}

func (dish *Dish) SetMealCourse(mealCourse string) {
	dish.MealCourse = mealCourse
	dish.UpdatedAt = time.Now()
}

func (dish *Dish) SetDietaryCategory(dietaryCategory string) {
	dish.DietaryCategory = dietaryCategory
	dish.UpdatedAt = time.Now()
}