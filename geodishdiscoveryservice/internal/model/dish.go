package model

import "time"

// Dish represents a dish document in the MongoDB collection
type Dish struct {
	DishID           string     `bson:"DishID" json:"DishID"`
	DishName         string     `bson:"DishName" json:"DishName"`
	ChefID           string     `bson:"ChefID" json:"ChefID"`
	Description      string     `bson:"Description" json:"Description"`
	Price            float64    `bson:"Price" json:"Price"`
	AvailablePortions int       `bson:"AvailablePortions" json:"AvailablePortions"`
	MealCourse       string     `bson:"MealCourse" json:"MealCourse"`
	DietaryCategory  string     `bson:"DietaryCategory" json:"DietaryCategory"`
	AvailableUntil   *time.Time  `bson:"AvailableUntil" json:"AvailableUntil"`
	Location         Location   `bson:"location" json:"location"` // GeoJSON format for geospatial queries
	Address          Address    `bson:"Address" json:"Address"`
	CreatedAt        time.Time  `bson:"CreatedAt" json:"CreatedAt"`
	UpdatedAt        time.Time  `bson:"UpdatedAt" json:"UpdatedAt"`
	DeletedAt        *time.Time `bson:"DeletedAt" json:"DeletedAt"`
}

// Location represents the GeoJSON format for location
type Location struct {
	Type        string    `bson:"type" json:"type"`               // Should always be "Point"
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // [longitude, latitude]
}

// Address represents the address information of the dish
type Address struct {
	Street     string `bson:"Street" json:"Street"`
	City       string `bson:"City" json:"City"`
	State      string `bson:"State" json:"State"`
	PostalCode string `bson:"PostalCode" json:"PostalCode"`
	Country    string `bson:"Country" json:"Country"`
}

// Image represents an image URL for the dish (if included in your system later)
type Image struct {
	ImageURL string `bson:"ImageURL" json:"ImageURL"`
}
