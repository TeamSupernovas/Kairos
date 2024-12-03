package model

import (
	"time"
)

type DishImage struct {
	DishImageID		string		`json:"dishImageId"`
	DishID			string		`json:"dishId"`
	ImageURL		string		`json:"imageURL"`
	CreatedAt		time.Time	`json:"createdAt"`
	DeletedAt		*time.Time	`json:"deletedAt"`
}

